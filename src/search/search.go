package search

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	cfp "github.com/magneticstain/ip-2-cloudresource/src/plugin/cloudfront"
	ec2p "github.com/magneticstain/ip-2-cloudresource/src/plugin/ec2"
	elbp "github.com/magneticstain/ip-2-cloudresource/src/plugin/elb"
	iamp "github.com/magneticstain/ip-2-cloudresource/src/plugin/iam"
	orgp "github.com/magneticstain/ip-2-cloudresource/src/plugin/organizations"
	generalResource "github.com/magneticstain/ip-2-cloudresource/src/resource"
	ipfuzzing "github.com/magneticstain/ip-2-cloudresource/src/svc/ip_fuzzing"
)

type Search struct {
	ac     awsconnector.AWSConnector
	ipAddr string
}

func NewSearch(ac awsconnector.AWSConnector, ipAddr string) Search {
	search := Search{ac: ac, ipAddr: ipAddr}

	return search
}

func reconcileCloudSvcParam(cloudSvc string) []string {
	var cloudSvcs []string

	if cloudSvc == "all" {
		cloudSvcs = []string{"cloudfront", "ec2", "elbv1", "elbv2"}
	} else if strings.Contains(cloudSvc, ",") {
		// csv provided, split the values into a slice
		cloudSvcs = strings.Split(cloudSvc, ",")
	} else {
		cloudSvcs = []string{cloudSvc}
	}

	return cloudSvcs
}

func (search Search) RunIPFuzzing(doAdvIPFuzzing bool) (string, error) {
	var fuzzedSvc string

	fuzzedSvc, err := ipfuzzing.FuzzIP(search.ipAddr, doAdvIPFuzzing)
	if err != nil {
		return fuzzedSvc, err
	} else if fuzzedSvc == "" || fuzzedSvc == "UNKNOWN" {
		log.Info("could not determine service via IP fuzzing")
	} else {
		log.Info("IP fuzzing determined the associated cloud service is: ", fuzzedSvc)
	}

	return fuzzedSvc, err
}

func (search Search) fetchOrgAcctIds(orgSearchOrgUnitID string, orgSearchXaccountRoleARN string) ([]string, error) {
	var acctIds []string

	// assume xaccount role first if ARN is provided
	if orgSearchXaccountRoleARN != "" {
		arac, err := awsconnector.NewAWSConnectorAssumeRole(orgSearchXaccountRoleARN, search.ac.AwsConfig)
		if err != nil {
			return acctIds, err
		}
		search.ac = arac
	}

	orgp := orgp.NewOrganizationsPlugin(search.ac, orgSearchOrgUnitID)
	orgAccts, err := orgp.GetResources()
	if err != nil {
		return acctIds, err
	}

	for _, acct := range orgAccts {
		log.Debug("org account found: ", *acct.Id, " (", *acct.Name, ") [ ", acct.Status, " ]")
		acctIds = append(acctIds, *acct.Id)
	}

	return acctIds, nil
}

func (search Search) SearchAWS(cloudSvc string) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource
	var err error

	cloudSvc = strings.ToLower(cloudSvc)

	log.Debug("searching ", cloudSvc, " in AWS")

	switch cloudSvc {
	case "cloudfront":
		pluginConn := cfp.NewCloudfrontPlugin(search.ac)
		matchingResource, err = pluginConn.SearchResources(search.ipAddr)
		if err != nil {
			return matchingResource, err
		}
	case "ec2":
		pluginConn := ec2p.NewEC2Plugin(search.ac)
		matchingResource, err = pluginConn.SearchResources(search.ipAddr)
		if err != nil {
			return matchingResource, err
		}
	case "elbv1": // classic ELBs
		pluginConn := elbp.NewELBv1Plugin(search.ac)
		matchingResource, err = pluginConn.SearchResources(search.ipAddr)
		if err != nil {
			return matchingResource, err
		}
	case "elbv2":
		pluginConn := elbp.NewELBPlugin(search.ac)
		matchingResource, err = pluginConn.SearchResources(search.ipAddr)
		if err != nil {
			return matchingResource, err
		}
	default:
		return matchingResource, errors.New("invalid cloud service provided for AWS search")
	}

	return matchingResource, nil
}

func (search Search) doAccountSearch(cloudSvcs []string, acctID string) (generalResource.Resource, error) {
	var acctAliases []string
	var matchingResource generalResource.Resource
	var err error

	if acctID != "current" {
		// resolve account's aliases
		iamp := iamp.NewIAMPlugin(search.ac)
		acctAliases, err = iamp.GetResources()
		if err != nil {
			return matchingResource, err
		}

		log.Info("starting AWS resource search in account: ", acctID, " ", acctAliases)
	} else {
		log.Info("starting AWS resource search in principal account")
	}

	for _, svc := range cloudSvcs {
		matchingResource, err = search.SearchAWS(svc)

		if err != nil {
			return matchingResource, err
		} else if matchingResource.RID != "" {
			// resource was found
			matchingResource.AccountID = acctID
			matchingResource.AccountAliases = acctAliases
			return matchingResource, nil
		}
	}

	return matchingResource, nil
}

func (search Search) runSearchWorker(matchingResourceBuffer chan<- generalResource.Resource, acctID string, cloudSvcs []string, orgSearchRoleName string, wg *sync.WaitGroup) {
	defer wg.Done()

	if acctID != "current" {
		// replace connector with assumed role connector
		acctRoleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", acctID, orgSearchRoleName)
		ac, err := awsconnector.NewAWSConnectorAssumeRole(acctRoleArn, aws.Config{})
		if err != nil {
			log.Error("error when assuming role for account search worker: ", err)
			return
		}

		search.ac = ac
	}

	resultResource, err := search.doAccountSearch(cloudSvcs, acctID)
	if err != nil {
		log.Error("error when running search withing account search worker: ", err)
	} else if resultResource.RID != "" {
		matchingResourceBuffer <- resultResource
		return
	}
}

func (search Search) InitSearch(cloudSvc string, doIPFuzzing bool, doAdvIPFuzzing bool, doOrgSearch bool, orgSearchXaccountRoleARN string, orgSearchRoleName string, orgSearchOrgUnitID string) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource
	var err error

	cloudSvcs := reconcileCloudSvcParam(cloudSvc)

	if doIPFuzzing || doAdvIPFuzzing {
		fuzzedSvc, err := search.RunIPFuzzing(doAdvIPFuzzing)
		if err != nil {
			return matchingResource, err
		}

		normalizedSvcName := strings.ToLower(fuzzedSvc)

		if normalizedSvcName != "unknown" {
			cloudSvcs = []string{normalizedSvcName}

			// all ELBs act within EC2 infrastructure, so we will need to add the elb services as well if that's the case
			if normalizedSvcName == "ec2" {
				cloudSvcs = []string{normalizedSvcName, "elbv1", "elbv2"}
			}
		}
	}

	var acctsToSearch []string
	if doOrgSearch {
		log.Info("starting org account enumeration")

		acctsToSearch, err = search.fetchOrgAcctIds(orgSearchOrgUnitID, orgSearchXaccountRoleARN)
		if err != nil {
			return matchingResource, err
		}
	} else {
		acctsToSearch = append(acctsToSearch, "current")
	}

	log.Info("beginning resource gathering")
	matchingResourceBuffer := make(chan generalResource.Resource, 1)
	var wg sync.WaitGroup

	for _, acctID := range acctsToSearch {
		wg.Add(1)
		go search.runSearchWorker(matchingResourceBuffer, acctID, cloudSvcs, orgSearchRoleName, &wg)
	}

	go func() {
		wg.Wait()
		close(matchingResourceBuffer)
	}()

	resultResource, found := <-matchingResourceBuffer
	if found {
		matchingResource = resultResource
	}

	return matchingResource, nil
}

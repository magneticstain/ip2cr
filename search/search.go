package search

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/rollbar/rollbar-go"
	log "github.com/sirupsen/logrus"

	awscontroller "github.com/magneticstain/ip-2-cloudresource/aws"
	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws/aws_connector"
	cfp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/cloudfront"
	ec2p "github.com/magneticstain/ip-2-cloudresource/aws/plugin/ec2"
	elbp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/elb"
	iamp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/iam"
	orgp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/organizations"
	ipfuzzing "github.com/magneticstain/ip-2-cloudresource/aws/svc/ip_fuzzing"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type Search struct {
	AWSCtrlr        awscontroller.AWSController
	CloudSvcs       []string
	IpAddr          string
	MatchedResource generalResource.Resource
	Platform        string
}

func (search *Search) connectToPlatform() (bool, error) {
	// generate a connection to the specified platform via plugin
	switch search.Platform {
	case "aws":
		ac, err := awscontroller.New()
		if err != nil {
			return false, err
		}

		search.AWSCtrlr = ac
	}

	return true, nil
}

func ReconcileCloudSvcParam(cloudSvc string) []string {
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

func (search Search) RunIPFuzzing(doAdvIPFuzzing bool) ([]string, error) {
	var svcSet []string
	var fuzzedSvc string
	var err error

	fuzzedSvc, err = ipfuzzing.FuzzIP(search.IpAddr, doAdvIPFuzzing)
	if err != nil {
		return svcSet, err
	}

	if fuzzedSvc == "" || fuzzedSvc == "unknown" {
		log.Info("could not determine service via IP fuzzing")
		return svcSet, err
	}

	log.Info("IP fuzzing determined the associated cloud service is: ", fuzzedSvc)
	svcSet = append(svcSet, fuzzedSvc)

	// all ELBs act within EC2 infrastructure, so we will need to add the elb services as well if that's the case
	if fuzzedSvc == "ec2" {
		svcSet = append(search.CloudSvcs, "elbv1", "elbv2")
	}

	return svcSet, err
}

func (search Search) fetchOrgAcctIds(orgSearchOrgUnitID string, orgSearchXaccountRoleARN string) ([]string, error) {
	var acctIds []string
	var err error

	// assume xaccount role first if ARN is provided
	var arac awsconnector.AWSConnector
	if orgSearchXaccountRoleARN != "" {
		arac, err = awsconnector.NewAWSConnectorAssumeRole(orgSearchXaccountRoleARN, search.AWSCtrlr.PrincipalAWSConn.AwsConfig)
		if err != nil {
			return acctIds, err
		}
	} else {
		arac = search.AWSCtrlr.PrincipalAWSConn
	}

	var orgAccts []types.Account
	orgp := orgp.OrganizationsPlugin{AwsConn: arac, OrgUnitID: orgSearchOrgUnitID}
	orgAccts, err = orgp.GetResources()
	if err != nil {
		return acctIds, err
	}

	for _, acct := range orgAccts {
		if acct.Status == "ACTIVE" {
			log.Debug("org account found: ", *acct.Id, " (", *acct.Name, ")")
			acctIds = append(acctIds, *acct.Id)
		} else {
			log.Debug("org account found, but not active: ", *acct.Id, " (", *acct.Name, ")")
		}
	}

	return acctIds, nil
}

func (search Search) SearchAWSSvc(cloudSvc string, doNetMapping bool) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource
	var err error

	cloudSvc = strings.ToLower(cloudSvc)

	log.Debug("searching ", cloudSvc, " in AWS")

	switch cloudSvc {
	case "cloudfront":
		pluginConn := cfp.CloudfrontPlugin{AwsConn: search.AWSCtrlr.PrincipalAWSConn, NetworkMapping: doNetMapping}
		matchingResource, err = pluginConn.SearchResources(search.IpAddr)
		if err != nil {
			return matchingResource, err
		}
	case "ec2":
		pluginConn := ec2p.EC2Plugin{AwsConn: search.AWSCtrlr.PrincipalAWSConn, NetworkMapping: doNetMapping}
		matchingResource, err = pluginConn.SearchResources(search.IpAddr)
		if err != nil {
			return matchingResource, err
		}
	case "elbv1": // classic ELBs
		pluginConn := elbp.ELBv1Plugin{AwsConn: search.AWSCtrlr.PrincipalAWSConn, NetworkMapping: doNetMapping}
		matchingResource, err = pluginConn.SearchResources(search.IpAddr)
		if err != nil {
			return matchingResource, err
		}
	case "elbv2":
		pluginConn := elbp.ELBPlugin{AwsConn: search.AWSCtrlr.PrincipalAWSConn, NetworkMapping: doNetMapping}
		matchingResource, err = pluginConn.SearchResources(search.IpAddr)
		if err != nil {
			return matchingResource, err
		}
	default:
		return matchingResource, errors.New("invalid cloud service provided for AWS search")
	}

	return matchingResource, nil
}

func (search Search) doAccountLevelSearch(acctID string, doNetMapping bool) (generalResource.Resource, error) {
	var acctAliases []string
	var matchingResource generalResource.Resource
	var err error

	if acctID != "current" && search.Platform == "aws" {
		// resolve account's aliases
		iamp := iamp.IAMPlugin{AwsConn: search.AWSCtrlr.PrincipalAWSConn}
		acctAliases, err = iamp.GetResources()
		if err != nil {
			return matchingResource, err
		}

		log.Info("starting resource search in AWS account: ", acctID, " ", acctAliases)
	} else {
		log.Info("starting resource search in principal account")
	}

	for _, svc := range search.CloudSvcs {
		switch search.Platform {
		case "aws":
			matchingResource, err = search.SearchAWSSvc(svc, doNetMapping)
		default:
			errorMsg := fmt.Sprintf("%s is not a supported platform for searching", search.Platform)
			return matchingResource, errors.New(errorMsg)
		}

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

func (search Search) runSearchWorker(matchingResourceBuffer chan<- generalResource.Resource, acctID string, orgSearchRoleName string, doNetMapping bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if acctID != "current" && search.Platform == "aws" {
		// replace connector with assumed role connector
		acctRoleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", acctID, orgSearchRoleName)
		ac, err := awsconnector.NewAWSConnectorAssumeRole(acctRoleArn, aws.Config{})
		if err != nil {
			log.Error("error when assuming role for account search worker: ", err)
			return
		}

		search.AWSCtrlr.PrincipalAWSConn = ac
	}

	resultResource, err := search.doAccountLevelSearch(acctID, doNetMapping)
	if err != nil {
		log.Error("error when running search within account search worker: ", err)
	} else if resultResource.RID != "" {
		matchingResourceBuffer <- resultResource
		return
	}
}

func (search Search) StartSearchWorkers(acctsToSearch []string, orgSearchRoleName string, doNetMapping bool) bool {
	log.Info("beginning resource gathering")

	matchingResourceBuffer := make(chan generalResource.Resource, 1)
	var wg sync.WaitGroup

	for _, acctID := range acctsToSearch {
		wg.Add(1)
		go rollbar.WrapAndWait(search.runSearchWorker, matchingResourceBuffer, acctID, orgSearchRoleName, doNetMapping, &wg)
	}

	go func() {
		wg.Wait()
		close(matchingResourceBuffer)
	}()

	resultResource, found := <-matchingResourceBuffer
	if found {
		search.MatchedResource = resultResource
	}

	return found
}

func (search Search) StartSearch(cloudSvc string, doIPFuzzing bool, doAdvIPFuzzing bool, doOrgSearch bool, orgSearchXaccountRoleARN string, orgSearchRoleName string, orgSearchOrgUnitID string, doNetMapping bool) (bool, error) {
	var resourceFound bool
	var err error

	_, err = search.connectToPlatform()
	if err != nil {
		log.Fatal("error when connecting to ", search.Platform, ": ", err)
	}

	// TODO: move this to init function
	search.CloudSvcs = ReconcileCloudSvcParam(cloudSvc)

	if doIPFuzzing || doAdvIPFuzzing {
		search.CloudSvcs, err = search.RunIPFuzzing(doAdvIPFuzzing)
		if err != nil {
			return resourceFound, err
		}
	}

	var acctsToSearch []string
	if doOrgSearch {
		log.Info("starting org account enumeration")

		acctsToSearch, err = search.fetchOrgAcctIds(orgSearchOrgUnitID, orgSearchXaccountRoleARN)
		if err != nil {
			return resourceFound, err
		}
	} else {
		acctsToSearch = append(acctsToSearch, "current")
	}

	resourceFound = search.StartSearchWorkers(acctsToSearch, orgSearchRoleName, doNetMapping)

	return resourceFound, nil
}

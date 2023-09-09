package search

import (
	"errors"
	"fmt"
	"strings"

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
	ac     *awsconnector.AWSConnector
	ipAddr *string
}

func NewSearch(ac *awsconnector.AWSConnector, ipAddr *string) Search {
	search := Search{ac: ac, ipAddr: ipAddr}

	return search
}

func (search Search) RunIPFuzzing(doAdvIPFuzzing bool) (*string, error) {
	var fuzzedSvc *string

	fuzzedSvc, err := ipfuzzing.FuzzIP(search.ipAddr, doAdvIPFuzzing)
	if err != nil {
		return fuzzedSvc, err
	} else if *fuzzedSvc == "" || *fuzzedSvc == "UNKNOWN" {
		log.Info("could not determine service via IP fuzzing")
	} else {
		log.Info("IP fuzzing determined the associated cloud service is: ", *fuzzedSvc)
	}

	return fuzzedSvc, err
}

func (search Search) fetchOrgAcctIds() (*[]string, error) {
	var acctIds []string

	orgp := orgp.NewOrganizationsPlugin(search.ac)
	orgAccts, err := orgp.GetResources()
	if err != nil {
		return &acctIds, err
	}

	for _, acct := range *orgAccts {
		log.Debug("org account found: ", *acct.Id, " (", *acct.Name, ") [ ", acct.Status, " ]")
		acctIds = append(acctIds, *acct.Id)
	}

	return &acctIds, nil
}

func (search Search) SearchAWS(cloudSvc string) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource
	cloudSvc = strings.ToLower(cloudSvc)

	log.Debug("searching ", cloudSvc, " in AWS")

	switch cloudSvc {
	case "cloudfront":
		pluginConn := cfp.NewCloudfrontPlugin(search.ac)
		cfResource, err := pluginConn.SearchResources(search.ipAddr)
		if err != nil {
			return matchingResource, err
		}

		if cfResource.ARN != nil {
			matchingResource.RID = *cfResource.ARN
			log.Debug("IP found as CloudFront distribution -> ", matchingResource.RID)
		}
	case "ec2":
		pluginConn := ec2p.NewEC2Plugin(search.ac)
		ec2Resource, err := pluginConn.SearchResources(search.ipAddr)
		if err != nil {
			return matchingResource, err
		}

		if ec2Resource.InstanceId != nil {
			matchingResource.RID = *ec2Resource.InstanceId // for some reason, the EC2 Instance object doesn't contain the ARN of the instance :/
			log.Debug("IP found as EC2 instance -> ", matchingResource.RID)
		}
	case "elbv1": // classic ELBs
		pluginConn := elbp.NewELBv1Plugin(search.ac)
		elbResource, err := pluginConn.SearchResources(search.ipAddr)
		if err != nil {
			return matchingResource, err
		}

		if elbResource.LoadBalancerName != nil { // no ARN available here either
			matchingResource.RID = *elbResource.LoadBalancerName
			log.Debug("IP found as Classic Elastic Load Balancer -> ", matchingResource.RID)
		}
	case "elbv2":
		pluginConn := elbp.NewELBPlugin(search.ac)
		elbResource, err := pluginConn.SearchResources(search.ipAddr)
		if err != nil {
			return matchingResource, err
		}

		if elbResource.LoadBalancerArn != nil {
			matchingResource.RID = *elbResource.LoadBalancerArn
			log.Debug("IP found as Elastic Load Balancer -> ", matchingResource.RID)
		}
	default:
		return matchingResource, errors.New("invalid cloud service provided for AWS search")
	}

	return matchingResource, nil
}

func (search Search) runSearch(cloudSvcs *[]string, acctID *string) (*generalResource.Resource, error) {
	var matchingResource generalResource.Resource

	if *acctID != "current" {
		// resolve account's aliases
		iamp := iamp.NewIAMPlugin(search.ac)
		acctAliases, err := iamp.GetResources()
		if err != nil {
			return &matchingResource, err
		}

		matchingResource.AccountAliases = acctAliases

		log.Info("starting AWS resource search in account: ", *acctID, " ", acctAliases)
	} else {
		log.Info("starting AWS resource search in principal account")
	}

	for _, svc := range *cloudSvcs {
		matchingResource, err := search.SearchAWS(svc)

		if err != nil {
			return &matchingResource, err
		} else if matchingResource.RID != "" {
			// resource was found
			matchingResource.AccountID = *acctID
			return &matchingResource, nil
		}
	}

	return &matchingResource, nil
}

func (search Search) InitSearch(doIPFuzzing bool, doAdvIPFuzzing bool, doOrgSearch bool, orgSearchRoleName string) (*generalResource.Resource, error) {
	var matchingResource *generalResource.Resource
	cloudSvcs := []string{"cloudfront", "ec2", "elbv1", "elbv2"}
	var err error

	if doIPFuzzing || doAdvIPFuzzing {
		fuzzedSvc, err := search.RunIPFuzzing(doAdvIPFuzzing)
		if err != nil {
			return matchingResource, err
		}

		normalizedSvcName := strings.ToLower(*fuzzedSvc)

		if normalizedSvcName != "unknown" {
			cloudSvcs = []string{normalizedSvcName}

			// all ELBs act within EC2 infrastructure, so we will need to add the elb services as well if that's the case
			if normalizedSvcName == "ec2" {
				cloudSvcs = append(cloudSvcs, "elbv1", "elbv2")
			}
		}
	}

	var acctsToSearch *[]string
	if doOrgSearch {
		log.Info("starting org account enumeration")

		acctsToSearch, err = search.fetchOrgAcctIds()
		if err != nil {
			return matchingResource, err
		}
	} else {
		acctsToSearch = &[]string{"current"}
	}

	log.Info("beginning resource gathering")
	var acctRoleArn string
	for _, acctID := range *acctsToSearch {
		if acctID != "current" {
			// replace connector with assumed role connector
			acctRoleArn = fmt.Sprintf("arn:aws:iam::%s:role/%s", acctID, orgSearchRoleName)
			ac, err := awsconnector.NewAWSConnectorAssumeRole(&acctRoleArn)
			if err != nil {
				return matchingResource, err
			}

			search.ac = &ac
		}

		matchingResource, err = search.runSearch(&cloudSvcs, &acctID)
		if err != nil {
			return matchingResource, err
		}

		if matchingResource.RID != "" {
			return matchingResource, err
		}
	}

	return matchingResource, nil
}

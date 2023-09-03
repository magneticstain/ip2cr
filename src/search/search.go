package search

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	"github.com/magneticstain/ip-2-cloudresource/src/plugin/cloudfront"
	"github.com/magneticstain/ip-2-cloudresource/src/plugin/ec2"
	"github.com/magneticstain/ip-2-cloudresource/src/plugin/elb"
	"github.com/magneticstain/ip-2-cloudresource/src/plugin/iam"
	"github.com/magneticstain/ip-2-cloudresource/src/plugin/organizations"
	generalResource "github.com/magneticstain/ip-2-cloudresource/src/resource"
	ipfuzzing "github.com/magneticstain/ip-2-cloudresource/src/svc/ip_fuzzing"
)

type Search struct {
	ac *awsconnector.AWSConnector
}

func NewSearch(ac *awsconnector.AWSConnector) Search {
	search := Search{ac: ac}

	return search
}

func (search Search) RunIpFuzzing(ipAddr *string) (*string, error) {
	var fuzzedSvc *string

	fuzzedSvc, err := ipfuzzing.FuzzIP(ipAddr, true)
	if err != nil {
		return fuzzedSvc, err
	} else if *fuzzedSvc == "" || *fuzzedSvc == "UNKNOWN" {
		log.Info("could not determine service via IP fuzzing")
	} else {
		log.Info("IP fuzzing determined the associated cloud service is: ", *fuzzedSvc)
	}

	return fuzzedSvc, err
}

func (search Search) FetchOrgAcctIds() (*[]string, error) {
	var acctIds []string

	orgp := organizations.NewOrganizationsPlugin(search.ac)
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

func (search Search) RunSearch(matchingResource *generalResource.Resource, cloudSvcs *[]string, ipAddr *string, acctId *string) (*generalResource.Resource, error) {
	for _, svc := range *cloudSvcs {
		cloudResource, err := search.SearchAWS(svc, ipAddr, matchingResource)

		if err != nil {
			return matchingResource, err
		} else if cloudResource.RID != "" {
			// resource was found
			matchingResource.AccountId = *acctId

			if *acctId != "current" {
				// resolve account's aliases
				iamp := iam.NewIAMPlugin(search.ac)
				acctAliases, err := iamp.GetResources()
				if err != nil {
					return matchingResource, err
				}

				matchingResource.AccountAliases = acctAliases
			}

			break
		}
	}

	return matchingResource, nil
}

func (search Search) SearchAWS(cloudSvc string, ipAddr *string, matchingResource *generalResource.Resource) (*generalResource.Resource, error) {
	cloudSvc = strings.ToLower(cloudSvc)

	log.Info("searching ", cloudSvc, " in AWS")

	switch cloudSvc {
	case "cloudfront":
		pluginConn := cloudfront.NewCloudfrontPlugin(search.ac)
		cf_resource, err := pluginConn.SearchResources(ipAddr)
		if err != nil {
			return matchingResource, err
		}

		if cf_resource.ARN != nil {
			matchingResource.RID = *cf_resource.ARN
			log.Debug("IP found as CloudFront distribution -> ", matchingResource.RID)
		}
	case "ec2":
		pluginConn := ec2.NewEC2Plugin(search.ac)
		ec2Resource, err := pluginConn.SearchResources(ipAddr)
		if err != nil {
			return matchingResource, err
		}

		if ec2Resource.InstanceId != nil {
			matchingResource.RID = *ec2Resource.InstanceId // for some reason, the EC2 Instance object doesn't contain the ARN of the instance :/
			log.Debug("IP found as EC2 instance -> ", matchingResource.RID)
		}
	case "elbv1": // classic ELBs
		pluginConn := elb.NewELBv1Plugin(search.ac)
		elb_resource, err := pluginConn.SearchResources(ipAddr)
		if err != nil {
			return matchingResource, err
		}

		if elb_resource.LoadBalancerName != nil { // no ARN available here either
			matchingResource.RID = *elb_resource.LoadBalancerName
			log.Debug("IP found as Classic Elastic Load Balancer -> ", matchingResource.RID)
		}
	case "elbv2":
		pluginConn := elb.NewELBPlugin(search.ac)
		elb_resource, err := pluginConn.SearchResources(ipAddr)
		if err != nil {
			return matchingResource, err
		}

		if elb_resource.LoadBalancerArn != nil {
			matchingResource.RID = *elb_resource.LoadBalancerArn
			log.Debug("IP found as Elastic Load Balancer -> ", matchingResource.RID)
		}
	default:
		return matchingResource, errors.New("invalid cloud service provided for AWS search")
	}

	return matchingResource, nil
}

func (search Search) StartSearch(ipAddr *string, doIpFuzzing bool, doAdvIpFuzzing bool, doOrgSearch bool, orgSearchRoleName string) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource
	cloudSvcs := []string{"cloudfront", "ec2", "elbv1", "elbv2"}

	if doIpFuzzing {
		fuzzedSvc, err := search.RunIpFuzzing(ipAddr)
		if err != nil {
			return matchingResource, err
		}

		normalizedSvcName := strings.ToLower(*fuzzedSvc)

		if normalizedSvcName != "unknown" {
			cloudSvcs = []string{}

			// all ELBs act within EC2 infrastructure, so we will need to add the elb services as well if that's the case
			if *fuzzedSvc == "EC2" {
				cloudSvcs = append(cloudSvcs, *fuzzedSvc, "elbv1", "elbv2")
			}
		}
	}

	var acctsToSearch *[]string
	if doOrgSearch {
		var err error

		log.Info("starting org account enumeration")
		acctsToSearch, err = search.FetchOrgAcctIds()
		if err != nil {
			return matchingResource, err
		}
	} else {
		acctsToSearch = &[]string{"current"}
	}

	log.Debug("beginning resource gathering")
	var acctRoleArn string
	for _, acctId := range *acctsToSearch {
		if acctId != "current" {
			// replace connector with assumed role connector
			acctRoleArn = fmt.Sprintf("arn:aws:iam::%s:role/%s", acctId, orgSearchRoleName)
			ac, err := awsconnector.NewAWSConnectorAssumeRole(&acctRoleArn)
			if err != nil {
				return matchingResource, err
			} else {
				search.ac = &ac
			}
		}

		search.RunSearch(&matchingResource, &cloudSvcs, ipAddr, &acctId)
	}

	return matchingResource, nil
}

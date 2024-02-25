package aws

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws/aws_connector"
	cfp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/cloudfront"
	ec2p "github.com/magneticstain/ip-2-cloudresource/aws/plugin/ec2"
	elbp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/elb"
	orgp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/organizations"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type AWSController struct {
	PrincipalAWSConn awsconnector.AWSConnector
}

func New() (AWSController, error) {
	awsConn, err := awsconnector.New()

	awsCtrlr := AWSController{PrincipalAWSConn: awsConn}

	return awsCtrlr, err
}

func (awsCtrlr AWSController) FetchOrgAcctIds(orgSearchOrgUnitID string, orgSearchXaccountRoleARN string) ([]string, error) {
	var acctIds []string
	var err error

	// assume xaccount role first if ARN is provided
	var arac awsconnector.AWSConnector
	if orgSearchXaccountRoleARN != "" {
		arac, err = awsconnector.NewAWSConnectorAssumeRole(orgSearchXaccountRoleARN, awsCtrlr.PrincipalAWSConn.AwsConfig)
		if err != nil {
			return acctIds, err
		}
	} else {
		arac = awsCtrlr.PrincipalAWSConn
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

func (awsCtrlr *AWSController) SearchAWSSvc(ipAddr, cloudSvc string, doNetMapping bool) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource
	var err error

	log.Debug("searching ", cloudSvc, " in AWS controller")

	switch cloudSvc {
	case "cloudfront":
		pluginConn := cfp.CloudfrontPlugin{AwsConn: awsCtrlr.PrincipalAWSConn, NetworkMapping: doNetMapping}
		matchingResource, err = pluginConn.SearchResources(ipAddr)
		if err != nil {
			return matchingResource, err
		}
	case "ec2":
		pluginConn := ec2p.EC2Plugin{AwsConn: awsCtrlr.PrincipalAWSConn, NetworkMapping: doNetMapping}
		matchingResource, err = pluginConn.SearchResources(ipAddr)
		if err != nil {
			return matchingResource, err
		}
	case "elbv1": // classic ELBs
		pluginConn := elbp.ELBv1Plugin{AwsConn: awsCtrlr.PrincipalAWSConn, NetworkMapping: doNetMapping}
		matchingResource, err = pluginConn.SearchResources(ipAddr)
		if err != nil {
			return matchingResource, err
		}
	case "elbv2":
		pluginConn := elbp.ELBPlugin{AwsConn: awsCtrlr.PrincipalAWSConn, NetworkMapping: doNetMapping}
		matchingResource, err = pluginConn.SearchResources(ipAddr)
		if err != nil {
			return matchingResource, err
		}
	default:
		return matchingResource, errors.New("invalid cloud service provided for AWS search")
	}

	return matchingResource, nil
}

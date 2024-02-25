package aws

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws/aws_connector"
	cfp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/cloudfront"
	ec2p "github.com/magneticstain/ip-2-cloudresource/aws/plugin/ec2"
	elbp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/elb"
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

func (awsCtrlr AWSController) SearchAWSSvc(ipAddr, cloudSvc string, doNetMapping bool) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource
	var err error

	cloudSvc = strings.ToLower(cloudSvc)

	log.Debug("searching ", cloudSvc, " in AWS")

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

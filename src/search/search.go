package search

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	"github.com/magneticstain/ip2cr/src/plugin/cloudfront"
	"github.com/magneticstain/ip2cr/src/plugin/ec2"
	"github.com/magneticstain/ip2cr/src/plugin/elb"
	generalResource "github.com/magneticstain/ip2cr/src/resource"
)

type Search struct {
	ac *awsconnector.AWSConnector
}

func NewSearch(ac *awsconnector.AWSConnector) Search {
	search := Search{ac: ac}

	return search
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
			matchingResource.RID = *ec2Resource.InstanceId  // for some reason, the EC2 Instance object doesn't contain the ARN of the instance :/
			log.Debug("IP found as EC2 instance -> ", matchingResource.RID)
		}
	case "elb":
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

func (search Search) StartSearch(ipAddr *string) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource
	cloudSvcs := []string{"cloudfront", "ec2", "elb"}

	log.Debug("beginning resource gathering")

	for _, svc := range cloudSvcs {
		cloudResource, err := search.SearchAWS(svc, ipAddr, &matchingResource)

		if err != nil {
			return matchingResource, err
		} else if cloudResource.RID != "" {
			// resource was found
			break
		}
	}

	return matchingResource, nil
}

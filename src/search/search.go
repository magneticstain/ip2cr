package search

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	"github.com/magneticstain/ip2cr/src/plugin/cloudfront"
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
	case "elb":
		pluginConn := elb.NewELBPlugin(search.ac)
		elb_resource := pluginConn.SearchResources(ipAddr)
		if elb_resource.LoadBalancerArn != nil { // TODO: the compiler freaks out when I try to point to this (and below) - idk why
			matchingResource.RID = *elb_resource.LoadBalancerArn
			log.Debug("IP found as Elastic Load Balancer -> ", matchingResource.RID)
		}
	case "cloudfront":
		pluginConn := cloudfront.NewCloudfrontPlugin(search.ac)
		cf_resource := pluginConn.SearchResources(ipAddr)
		if cf_resource.ARN != nil {
			matchingResource.RID = *cf_resource.ARN
			log.Debug("IP found as CloudFront distribution -> ", matchingResource.RID)
		}
	default:
		return matchingResource, errors.New("invalid cloud service provided for AWS search")
	}

	return matchingResource, nil
}

func (search Search) StartSearch(ipAddr *string) generalResource.Resource {
	var matchingResource generalResource.Resource
	cloudSvcs := []string{"cloudfront", "elb"}

	log.Debug("beginning resource gathering")

	for _, svc := range cloudSvcs {
		cloudResource, err := search.SearchAWS(svc, ipAddr, &matchingResource)
		if err != nil {
			log.Error(err)
		} else if cloudResource.RID != "" {
			return matchingResource
		}
	}

	return matchingResource
}

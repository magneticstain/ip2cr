package search

import (
	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip2cr/src/models/aws_connector"
	"github.com/magneticstain/ip2cr/src/models/cloudfront"
	generalResource "github.com/magneticstain/ip2cr/src/models/resource"
)

type Search struct{}

func StartSearch_AWSCloudfront(ac *awsconnector.AWSConnector, ipAddr *string, matchingResource *generalResource.Resource) generalResource.Resource {
	log.Info("searching AWS CloudFront")

	cfp := cloudfront.NewCloudfrontPlugin(ac)
	cf_resource := cfp.SearchResources(ipAddr)
	if *cf_resource.ARN != "" {
		matchingResource.RID = *cf_resource.ARN
		log.Debug("IP found as CloudFront distribution ", matchingResource.RID)
	}

	return *matchingResource
}

func StartSearch(ac *awsconnector.AWSConnector, ipAddr *string) generalResource.Resource {
	var matchingResource generalResource.Resource

	log.Info("beginning resource gathering")

	cf := StartSearch_AWSCloudfront(ac, ipAddr, &matchingResource)
	if cf.RID != "" {
		return matchingResource
	}

	return matchingResource
}

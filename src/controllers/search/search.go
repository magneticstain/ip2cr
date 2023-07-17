package search

import (
	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip2cr/src/models/aws_connector"
	"github.com/magneticstain/ip2cr/src/models/cloudfront"
)

type Search struct{}

func StartSearch(ac *awsconnector.AWSConnector) {
	log.Info("beginning resource gathering")

	resource_set := make(map[string]any)

	log.Info("searching AWS CloudFront")
	cfp := cloudfront.NewCloudfrontPlugin(ac)
	resource_set["cloudfront"] = cfp.SearchResources()

	log.Info(resource_set)
}

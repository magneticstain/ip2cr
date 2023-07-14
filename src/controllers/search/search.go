package search

import (
	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip2cr/src/models/aws_connector"
)

func StartSearch(aws_conn *awsconnector.AWSConnector) {
	log.Info("Beginning search...")
}

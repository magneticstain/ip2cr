package awsconnector

import (
	"context"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSConnector struct {
	cfg aws.Config
}

func New() AWSConnector {
	ac := AWSConnector{cfg: ConnectToAWS()}

	return ac
}

func ConnectToAWS() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

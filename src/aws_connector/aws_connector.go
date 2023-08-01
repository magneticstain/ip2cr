package awsconnector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSConnector struct {
	AwsConfig aws.Config
}

func New() (AWSConnector, error) {
	cfg, err := ConnectToAWS()

	ac := AWSConnector{AwsConfig: cfg}

	return ac, err
}

func ConnectToAWS() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	return cfg, err
}

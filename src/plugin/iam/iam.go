package plugin

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
)

type IAMPlugin struct {
	AwsConn awsconnector.AWSConnector
}

func NewIAMPlugin(awsConn *awsconnector.AWSConnector) IAMPlugin {
	iamp := IAMPlugin{AwsConn: *awsConn}

	return iamp
}

func (iamp IAMPlugin) GetResources() ([]string, error) {
	var acctAliases []string

	iamClient := iam.NewFromConfig(iamp.AwsConn.AwsConfig)

	iamResources, err := iamClient.ListAccountAliases(context.TODO(), &iam.ListAccountAliasesInput{})
	if err != nil {
		return acctAliases, err
	}

	acctAliases = iamResources.AccountAliases

	return acctAliases, nil
}

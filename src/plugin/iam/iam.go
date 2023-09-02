package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	generalPlugin "github.com/magneticstain/ip-2-cloudresource/src/plugin"
)

type IAMPlugin struct {
	GenPlugin *generalPlugin.Plugin
	AwsConn   awsconnector.AWSConnector
}

func NewIAMPlugin(aws_conn *awsconnector.AWSConnector) IAMPlugin {
	iamp := IAMPlugin{GenPlugin: &generalPlugin.Plugin{}, AwsConn: *aws_conn}

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

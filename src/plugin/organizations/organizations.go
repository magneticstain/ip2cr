package organizations

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
)

type OrganizationsPlugin struct {
	AwsConn awsconnector.AWSConnector
}

func NewOrganizationsPlugin(aws_conn *awsconnector.AWSConnector) OrganizationsPlugin {
	orgp := OrganizationsPlugin{AwsConn: *aws_conn}

	return orgp
}

func (orgp OrganizationsPlugin) GetResources() (*[]types.Account, error) {
	var orgAccts []types.Account

	orgClient := organizations.NewFromConfig(orgp.AwsConn.AwsConfig)
	paginator := organizations.NewListAccountsPaginator(orgClient, &organizations.ListAccountsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return &orgAccts, err
		}

		orgAccts = append(orgAccts, output.Accounts...)
	}

	return &orgAccts, nil
}

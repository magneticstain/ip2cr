package plugin

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
)

type OrganizationsPlugin struct {
	AwsConn   awsconnector.AWSConnector
	OrgUnitID string
}

func NewOrganizationsPlugin(awsConn *awsconnector.AWSConnector) OrganizationsPlugin {
	orgp := OrganizationsPlugin{AwsConn: *awsConn}

	return orgp
}

func listAllAccountsInOrganization(orgClient *organizations.Client) (*[]types.Account, error) {
	var orgAccts []types.Account

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

func (orgp OrganizationsPlugin) listAllAccountsInOrganizationalUnit(orgClient *organizations.Client) (*[]types.Account, error) {
	var orgAccts []types.Account

	paginator := organizations.NewListAccountsForParentPaginator(orgClient, &organizations.ListAccountsForParentInput{
		ParentId: &orgp.OrgUnitID,
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return &orgAccts, err
		}

		orgAccts = append(orgAccts, output.Accounts...)
	}

	return &orgAccts, nil
}

func (orgp OrganizationsPlugin) GetResources() (*[]types.Account, error) {
	var orgAccts *[]types.Account
	var err error

	orgClient := organizations.NewFromConfig(orgp.AwsConn.AwsConfig)

	if orgp.OrgUnitID != "" {
		orgAccts, err = orgp.listAllAccountsInOrganizationalUnit(orgClient)
	} else {
		orgAccts, err = listAllAccountsInOrganization(orgClient)
	}

	if err != nil {
		return orgAccts, err
	}

	return orgAccts, nil
}

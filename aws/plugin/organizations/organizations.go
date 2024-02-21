package plugin

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws/aws_connector"
)

type OrganizationsPlugin struct {
	AwsConn   awsconnector.AWSConnector
	OrgUnitID string
}

func listAllAccountsInOrganization(orgClient organizations.ListAccountsAPIClient) ([]types.Account, error) {
	var orgAccts []types.Account

	paginator := organizations.NewListAccountsPaginator(orgClient, &organizations.ListAccountsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return orgAccts, err
		}

		orgAccts = append(orgAccts, output.Accounts...)
	}

	return orgAccts, nil
}

func (orgp OrganizationsPlugin) listAllAccountsInOrganizationalUnit(orgClient organizations.ListAccountsForParentAPIClient) ([]types.Account, error) {
	var orgAccts []types.Account

	paginator := organizations.NewListAccountsForParentPaginator(orgClient, &organizations.ListAccountsForParentInput{
		ParentId: &orgp.OrgUnitID,
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return orgAccts, err
		}

		orgAccts = append(orgAccts, output.Accounts...)
	}

	return orgAccts, nil
}

func (orgp OrganizationsPlugin) GetResources() ([]types.Account, error) {
	var orgAccts []types.Account
	var err error

	orgClient := organizations.NewFromConfig(orgp.AwsConn.AwsConfig)

	if orgp.OrgUnitID != "" {
		log.Debug("fetching accounts from specified OU (", orgp.OrgUnitID, ")")
		orgAccts, err = orgp.listAllAccountsInOrganizationalUnit(orgClient)
	} else {
		orgAccts, err = listAllAccountsInOrganization(orgClient)
	}

	if err != nil {
		return orgAccts, err
	}

	return orgAccts, nil
}

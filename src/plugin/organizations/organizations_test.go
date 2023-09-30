package plugin_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	plugin "github.com/magneticstain/ip-2-cloudresource/src/plugin/organizations"
)

func orgFactory() plugin.OrganizationsPlugin {
	ac, _ := awsconnector.New()
	OUID := ""

	orgp := plugin.NewOrganizationsPlugin(&ac, &OUID)

	return orgp
}

func TestGetResources(t *testing.T) {
	orgp := orgFactory()

	orgResources, _ := orgp.GetResources()

	expectedType := "Account"
	for _, acct := range *orgResources {
		orgType := reflect.TypeOf(acct)
		if orgType.Name() != expectedType {
			t.Errorf("Fetching resources via AWS Organizations Plugin failed; wanted %s type, received %s", expectedType, orgType.Name())
		}
	}
}

func TestGetResources_TgtOUID(t *testing.T) {
	// REF: https://docs.aws.amazon.com/organizations/latest/APIReference/API_Organization.html#organizations-Type-Organization-Id

	var tests = []struct {
		orgID, ipAddr string
	}{
		{"o-0000000000", "1.1.1.1"},
		{"o-9999999999", "1.1.1.1"},
		{"o-1234567890abcde", "1.1.1.1"},
		{"o-00000000000000000000000000000000", "1.1.1.1"},
		{"o-99999999999999999999999999999999", "1.1.1.1"},
	}

	for _, td := range tests {
		testName := td.orgID

		orgp := orgFactory()
		orgp.OrgUnitID = td.orgID

		t.Run(testName, func(t *testing.T) {
			orgResources, _ := orgp.GetResources()

			expectedType := "Account"
			for _, acct := range *orgResources {
				orgType := reflect.TypeOf(acct)
				if orgType.Name() != expectedType {
					t.Errorf("Fetching resources with specific Organizational Unit ID via AWS Organizations Plugin failed; wanted %s type, received %s", expectedType, orgType.Name())
				}
			}
		})
	}
}

package organizations_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	"github.com/magneticstain/ip-2-cloudresource/src/plugin/organizations"
)

func orgFactory() organizations.OrganizationsPlugin {
	ac, _ := awsconnector.New()

	orgp := organizations.NewOrganizationsPlugin(&ac)

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

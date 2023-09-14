package plugin_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	plugin "github.com/magneticstain/ip-2-cloudresource/src/plugin/iam"
)

func iampFactory() plugin.IAMPlugin {
	ac, _ := awsconnector.New()

	iamp := plugin.NewIAMPlugin(&ac)

	return iamp
}

func TestGetResources(t *testing.T) {
	iamp := iampFactory()

	iamResources, _ := iamp.GetResources()

	expectedType := "slice"
	for _, alias := range iamResources {
		iamType := reflect.TypeOf(alias)
		if iamType.Name() != expectedType {
			t.Errorf("Fetching resources via IAM Plugin failed; wanted %s type, received %s", expectedType, iamType.Name())
		}
	}
}

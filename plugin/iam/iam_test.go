package plugin_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws_connector"
	plugin "github.com/magneticstain/ip-2-cloudresource/plugin/iam"
)

func iampFactory() plugin.IAMPlugin {
	ac, _ := awsconnector.New()

	iamp := plugin.IAMPlugin{AwsConn: ac}

	return iamp
}

func TestGetResources(t *testing.T) {
	iamp := iampFactory()

	iamResources, _ := iamp.GetResources()

	expectedType := "string"
	for _, alias := range iamResources {
		iamAliasType := reflect.TypeOf(alias)
		if iamAliasType.Name() != expectedType {
			t.Errorf("Fetching account alias via IAM Plugin failed; wanted %s type, received %s", expectedType, iamAliasType.Name())
		}
	}
}

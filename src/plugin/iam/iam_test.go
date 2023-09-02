package iam_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	"github.com/magneticstain/ip-2-cloudresource/src/plugin/iam"
)

func iampFactory() iam.IAMPlugin {
	ac, _ := awsconnector.New()

	iamp := iam.NewIAMPlugin(&ac)

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

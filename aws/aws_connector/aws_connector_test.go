package awsconnector_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws/aws_connector"
)

func TestConnectToAWS(t *testing.T) {
	ac, _ := awsconnector.New()

	acType := reflect.TypeOf(ac.AwsConfig)

	if acType.Name() != "Config" {
		t.Errorf("AWS connector failed to connect; wanted aws.Config type, received %s", acType.Name())
	}
}
func TestNewAWSConnectorAssumeRole(t *testing.T) {

	var tests = []struct {
		roleArn string
	}{
		{"arn:aws:iam::123456789012:role/valid_role"},
		{"arn:aws:bad::123456:role/invalid_role"},
		{"arn:aws:iam::123456789012:user/invalid_user"}, // only roles are supported by this function
	}

	for _, td := range tests {
		testName := td.roleArn

		t.Run(testName, func(t *testing.T) {
			ac, _ := awsconnector.NewAWSConnectorAssumeRole(td.roleArn, aws.Config{})

			acType := reflect.TypeOf(ac.AwsConfig)

			if acType.Name() != "Config" {
				t.Errorf("AWS connector failed to connect when assuming test role; wanted aws.Config type, received %s", acType.Name())
			}
		})
	}
}

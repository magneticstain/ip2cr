package awsconnector_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
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
		valid   bool
	}{
		{"arn:aws:iam::123456789012:role/valid_role", true},
		// {"arn:aws:bad::123456:role/invalid_role", false},
		// {"arn:aws:iam::123456789012:user/invalid_user", false}, // only roles are supported by this function
	}

	for _, td := range tests {
		testName := td.roleArn

		t.Run(testName, func(t *testing.T) {
			ac, err := awsconnector.NewAWSConnectorAssumeRole(&td.roleArn)

			if td.valid {
				if err != nil {
					t.Errorf("Instantiation of AWS connector using role unexpectedly failed; expected no error, received %s", err)
				} else {
					acType := reflect.TypeOf(ac.AwsConfig)

					if acType.Name() != "Config" {
						t.Errorf("AWS connector failed to connect when assuming test role; wanted aws.Config type, received %s", acType.Name())
					}
				}
			} else {
				if err == nil {
					acType := reflect.TypeOf(ac.AwsConfig)

					t.Errorf("Instantiation of AWS connector using role was expected to fail with an error, but succeeded; obj type: %s", acType.Name())
				}
			}
		})
	}
}

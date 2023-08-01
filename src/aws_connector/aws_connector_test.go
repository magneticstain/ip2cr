package awsconnector_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
)

func TestConnectToAWS(t *testing.T) {
	ac, _ := awsconnector.New()

	acType := reflect.TypeOf(ac.AwsConfig)

	if acType.Name() != "Config" {
		t.Errorf("AWS connector failed to connect; wanted aws.Config type, received %s", acType.Name())
	}
}

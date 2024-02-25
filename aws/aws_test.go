package aws_test

import (
	"fmt"
	"reflect"
	"testing"

	awscontroller "github.com/magneticstain/ip-2-cloudresource/aws"
)

func awsControllerFactory() awscontroller.AWSController {
	ac, _ := awscontroller.New()

	return ac
}

func TestSearchAWS(t *testing.T) {
	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"cloudfront", "1.1.1.1"},
		{"ec2", "1.1.1.1"},
		{"elbv1", "1.1.1.1"},
		{"ELBv1", "1.1.1.1"},
		{"elbv2", "1.1.1.1"},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		ac := awsControllerFactory()

		t.Run(testName, func(t *testing.T) {
			res, _ := ac.SearchAWSSvc(td.ipAddr, td.cloudSvc, false)

			resType := reflect.TypeOf(res)
			expectedType := "bool"
			if resType.Name() != expectedType {
				t.Errorf("AWS resource search failed; expected %s after search, received %s", expectedType, resType.Name())
			}
		})
	}
}

func TestSearchAWS_UnknownCloudSvc(t *testing.T) {
	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"magic_svc", "1.1.1.1"},   // known invalid use case
		{"cloudfront-", "1.1.1.1"}, // known valid use case with one character addition to make it invalid
		{"iam", "1.1.1.1"},         // valid AWS service, but not one that would ever interact with IP addresses
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		ac := awsControllerFactory()

		t.Run(testName, func(t *testing.T) {
			_, err := ac.SearchAWSSvc(td.ipAddr, td.cloudSvc, false)
			if err == nil {
				t.Errorf("Error was expected, but not seen, when performing general search; using %s for unknown cloud service key", td.cloudSvc)
			}
		})
	}
}

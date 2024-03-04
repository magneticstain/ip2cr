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

func TestFetchOrgAcctIds(t *testing.T) {
	var tests = []struct {
		orgSearchOrgUnitID, orgSearchXaccountRoleARN string
	}{
		{"abcde", "arn:aws:iam::123456789012:role/valid_role"},
		{"ou-", "arn:aws:bad::123456:role/invalid_role"},
		{"ou-a", "arn:aws:bad::123456:role/invalid_role"},
		{"ou-1234567890abcde", "arn:aws:iam::123456789012:role/valid_role"},
		{"ou-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "arn:aws:iam::123456789012:role/valid_role"},
		{"ou-zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", "arn:aws:iam::123456789012:role/valid_role"},
	}

	for _, td := range tests {
		testName := td.orgSearchOrgUnitID

		ac := awsControllerFactory()

		t.Run(testName, func(t *testing.T) {
			res, _ := ac.FetchOrgAcctIds(td.orgSearchOrgUnitID, td.orgSearchXaccountRoleARN)

			if len(res) != 0 {
				t.Errorf("AWS Orgs account ID fetch failed; expected 0 results from fetch, received %d", len(res))
			}
		})
	}
}

func TestSearchAWSSvc(t *testing.T) {
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
			expectedType := "Resource"
			if resType.Name() != expectedType {
				t.Errorf("AWS resource search failed; expected %s after search, received %s", expectedType, resType.Name())
			}
		})
	}
}

func TestSearchAWSSvc_UnknownCloudSvc(t *testing.T) {
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

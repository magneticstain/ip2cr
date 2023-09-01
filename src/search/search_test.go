package search_test

import (
	"fmt"
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	generalResource "github.com/magneticstain/ip-2-cloudresource/src/resource"
	"github.com/magneticstain/ip-2-cloudresource/src/search"
)

type TestIpAddr struct {
	ipAddr string
}

func searchFactory() search.Search {
	ac, _ := awsconnector.New()

	search := search.NewSearch(&ac)

	return search
}

func ipFactory() []TestIpAddr {
	var ipData []TestIpAddr

	ipData = append(
		ipData,
		TestIpAddr{"52.4.175.237"}, // cloudfront
		TestIpAddr{"65.8.191.186"}, // ALB
		TestIpAddr{"35.170.192.9"}, // EC2
		TestIpAddr{"2600:1f18:243e:1300:4685:5a7:7c28:c53a"}, // EC2 IPv6
		TestIpAddr{"3.218.196.10"},                           // NLB
		TestIpAddr{"34.205.13.193"},                          // Classic ELB
	)

	return ipData
}

func TestSearchAWS(t *testing.T) {
	search := searchFactory()

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

		matchedResource := generalResource.Resource{}
		t.Run(testName, func(t *testing.T) {
			res, _ := search.SearchAWS(td.cloudSvc, &td.ipAddr, &matchedResource)

			matchedResourceType := reflect.TypeOf(*res)
			expectedType := "Resource"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("AWS resource search failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestSearchAWS_UnknownCloudSvc(t *testing.T) {
	search := searchFactory()

	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"magic_svc", "1.1.1.1"},   // known invalid use case
		{"cloudfront-", "1.1.1.1"}, // known valid use case with one character addition to make it invalid
		{"iam", "1.1.1.1"},         // valid AWS service, but not one that would ever interact with IP addresses
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		matchedResource := generalResource.Resource{}
		t.Run(testName, func(t *testing.T) {
			_, err := search.SearchAWS(td.cloudSvc, &td.ipAddr, &matchedResource)
			if err == nil {
				t.Errorf("Error was expected, but not seen, when performing general search; using %s for unknown cloud service key", td.cloudSvc)
			}
		})
	}
}

func TestStartSearch_NoFuzzing(t *testing.T) {
	search := searchFactory()

	var tests = []struct {
		ipAddr string
	}{
		{"1.1.1.1"},
		{"299.11.906.43"},
		{"2600:9000:24eb:dc00:1:3b80:4f00:21"},
		{"x2600:9000:24eb:XYZ1:1:3b80:4f00:21"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch(&td.ipAddr, false, false, false, "")

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "Resource"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with IP fuzzing disabled has failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_BasicFuzzing(t *testing.T) {
	search := searchFactory()

	var tests = ipFactory()

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch(&td.ipAddr, true, false, false, "")

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "Resource"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with IP fuzzing enabled failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_AdvancedFuzzing(t *testing.T) {
	search := searchFactory()

	var tests = ipFactory()

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch(&td.ipAddr, true, false, false, "")

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "Resource"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with advanced IP fuzzing disabled failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_OrgSearchEnabled(t *testing.T) {
	search := searchFactory()

	var tests = ipFactory()

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch(&td.ipAddr, false, false, true, "ip2cr-org-role")

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "Resource"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with AWS Organizations support enabled has failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

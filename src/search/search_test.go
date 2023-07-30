package search_test

import (
	"fmt"
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	generalResource "github.com/magneticstain/ip2cr/src/resource"
	"github.com/magneticstain/ip2cr/src/search"
)

func searchFactory() search.Search {
	ac := awsconnector.New()

	search := search.NewSearch(&ac)

	return search
}

func TestSearchAWS(t *testing.T) {
	search := searchFactory()

	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"cloudfront", "1.1.1.1"},
		{"elb", "1.1.1.1"},
		{"ELB", "1.1.1.1"},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		matchedResource := generalResource.Resource{}
		t.Run(testName, func(t *testing.T) {
			res, err := search.SearchAWS(td.cloudSvc, &td.ipAddr, &matchedResource)
			if err != nil {
				t.Errorf("Encountered unexpected error when running AWS search; received error %s when searching %s for %s", err, td.cloudSvc, td.ipAddr)
			}

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
		{"cloudfront_", "1.1.1.1"}, // known valid use case with one character addition to make it invalid
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

func TestStartSearch(t *testing.T) {
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
			res := search.StartSearch(&td.ipAddr)

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "Resource"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

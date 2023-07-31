package cloudfront_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	"github.com/magneticstain/ip2cr/src/plugin/cloudfront"
)

func cfpFactory() cloudfront.CloudfrontPlugin {
	ac := awsconnector.New()

	cfp := cloudfront.NewCloudfrontPlugin(&ac)

	return cfp
}

func TestNormalizeCFDistroFQDN(t *testing.T) {
	cfp := cfpFactory()

	var tests = []struct {
		origFQDN, normalizedFQDN string
	}{
		{"1234567890abcd.cloudfront.net.", "1234567890abcd.cloudfront.net"},
		{"1234567890abcd.cloudfront.net", "1234567890abcd.cloudfront.net"},
		{"1234567890abcd.cloudfront...net.", "1234567890abcd.cloudfront...net"}, // function only removes trailing period; everything else can/should be left intact
	}

	for _, td := range tests {
		testName := td.origFQDN

		t.Run(testName, func(t *testing.T) {
			normalizedFQDN := cfp.NormalizeCFDistroFQDN(&td.origFQDN)

			if normalizedFQDN != td.normalizedFQDN {
				t.Errorf("CloudFront distribution domain normalization failed; expected %s, received %s", td.normalizedFQDN, normalizedFQDN)
			}
		})
	}
}

func TestGetResources(t *testing.T) {
	cfp := cfpFactory()

	cfResources := cfp.GetResources()

	expectedType := "DistributionSummary"
	for _, cfDistro := range *cfResources {
		cfDistroType := reflect.TypeOf(cfDistro)
		if cfDistroType.Name() != expectedType {
			t.Errorf("Fetching resources via Cloudfront Plugin failed; wanted %s type, received %s", expectedType, cfDistroType.Name())
		}
	}
}

func TestSearchResources(t *testing.T) {
	cfp := cfpFactory()

	var tests = []struct {
		ipAddr, expectedType string
	}{
		{"1.1.1.1", "DistributionSummary"},
		{"1234.45.9666.1", "DistributionSummary"},
		{"18.161.22.61", "DistributionSummary"},
		{"2600:9000:24eb:dc00:1:3b80:4f00:21", "DistributionSummary"},
		{"x2600:9000:24eb:XYZ1:1:3b80:4f00:21", "DistributionSummary"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			matchedDistro := cfp.SearchResources(&td.ipAddr)
			matchedDistroType := reflect.TypeOf(*matchedDistro)

			if matchedDistroType.Name() != td.expectedType {
				t.Errorf("CloudFront distribution search failed; expected %s after search, received %s", td.expectedType, matchedDistroType.Name())
			}
		})
	}
}
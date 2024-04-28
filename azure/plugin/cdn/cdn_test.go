package cdn_test

import (
	"reflect"
	"testing"

	plugin "github.com/magneticstain/ip-2-cloudresource/azure/plugin/cdn"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

func azcdnPlugFactory() plugin.AzCDNPlugin {
	azcdnPlug := plugin.AzCDNPlugin{}

	return azcdnPlug
}

func TestGetResources(t *testing.T) {
	azcdnPlug := azcdnPlugFactory()

	cdnResources, _ := azcdnPlug.GetResources()

	expectedType := "Resource"
	for _, resource := range cdnResources {
		resourceType := reflect.TypeOf(resource)
		if resourceType.Name() != expectedType {
			t.Errorf("Fetching resources via Azure CDN Plugin failed; wanted %s type, received %s", expectedType, resourceType.Name())
		}
	}
}

func TestSearchResources(t *testing.T) {
	azcdnPlug := azcdnPlugFactory()

	var tests = []struct {
		ipAddr, expectedType string
	}{
		{"1.1.1.1", "Resource"},
		{"1234.45.9666.1", "Resource"},
		{"18.161.22.61", "Resource"},
		{"2600:9000:24eb:dc00:1:3b80:4f00:21", "Resource"},
		{"x2600:9000:24eb:XYZ1:1:3b80:4f00:21", "Resource"},
	}

	var matchingResource generalResource.Resource
	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			matchedCdnEndpoint, _ := azcdnPlug.SearchResources(td.ipAddr, &matchingResource)
			matchedCdnEndpointType := reflect.TypeOf(*matchedCdnEndpoint)

			if matchedCdnEndpointType.Name() != td.expectedType {
				t.Errorf("Azure CDN search failed; expected %s after search, received %s", td.expectedType, matchedCdnEndpointType.Name())
			}
		})
	}
}

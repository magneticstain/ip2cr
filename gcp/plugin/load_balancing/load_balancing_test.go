package load_balancing_test

import (
	"reflect"
	"testing"

	plugin "github.com/magneticstain/ip-2-cloudresource/gcp/plugin/load_balancing"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

func lbPlugFactory() plugin.LoadBalancingPlugin {
	lbPlug := plugin.LoadBalancingPlugin{}

	return lbPlug
}

func TestGetResources(t *testing.T) {
	lbPlug := lbPlugFactory()

	lbResources, _ := lbPlug.GetResources()

	expectedType := "LoadBalancingResource"
	for _, resource := range lbResources {
		resourceType := reflect.TypeOf(resource)
		if resourceType.Name() != expectedType {
			t.Errorf("Fetching resources via GCP Load Balancing Plugin failed; wanted %s type, received %s", expectedType, resourceType.Name())
		}
	}
}

func TestSearchResources(t *testing.T) {
	lbPlug := lbPlugFactory()

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
			matchedInstance, _ := lbPlug.SearchResources(td.ipAddr, &matchingResource)
			matchedInstanceType := reflect.TypeOf(matchedInstance)

			if matchedInstanceType.Name() != td.expectedType {
				t.Errorf("GCP Load Balancing search failed; expected %s after search, received %s", td.expectedType, matchedInstanceType.Name())
			}
		})
	}
}

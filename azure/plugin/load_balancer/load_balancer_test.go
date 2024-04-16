package load_balancer_test

import (
	"reflect"
	"testing"

	plugin "github.com/magneticstain/ip-2-cloudresource/azure/plugin/load_balancer"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

func azlbPlugFactory() plugin.AzLoadBalancerPlugin {
	azlbPlug := plugin.AzLoadBalancerPlugin{}

	return azlbPlug
}

func TestGetResources(t *testing.T) {
	azlbPlug := azlbPlugFactory()

	lbResources, _ := azlbPlug.GetResources()

	expectedType := "Resource"
	for _, resource := range lbResources {
		resourceType := reflect.TypeOf(resource)
		if resourceType.Name() != expectedType {
			t.Errorf("Fetching resources via Azure Load Balancer Plugin failed; wanted %s type, received %s", expectedType, resourceType.Name())
		}
	}
}

func TestSearchResources(t *testing.T) {
	azlbPlug := azlbPlugFactory()

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
			matchedLB, _ := azlbPlug.SearchResources(td.ipAddr, &matchingResource)
			matchedLBType := reflect.TypeOf(*matchedLB)

			if matchedLBType.Name() != td.expectedType {
				t.Errorf("Azure Load Balancer search failed; expected %s after search, received %s", td.expectedType, matchedLBType.Name())
			}
		})
	}
}

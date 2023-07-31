package elb_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	"github.com/magneticstain/ip2cr/src/plugin/elb"
)

func elbpFactory() elb.ELBPlugin {
	ac := awsconnector.New()

	elbp := elb.NewELBPlugin(&ac)

	return elbp
}

func TestGetResources(t *testing.T) {
	elbp := elbpFactory()

	elbResources := elbp.GetResources()

	expectedType := "LoadBalancer"
	for _, elb := range *elbResources {
		elbType := reflect.TypeOf(elb)
		if elbType.Name() != expectedType {
			t.Errorf("Fetching resources via ELB Plugin failed; wanted %s type, received %s", expectedType, elbType.Name())
		}
	}
}

func TestSearchResources(t *testing.T) {
	elbp := elbpFactory()

	var tests = []struct {
		ipAddr, expectedType string
	}{
		{"1.1.1.1", "LoadBalancer"},
		{"1234.45.9666.1", "LoadBalancer"},
		{"18.161.22.61", "LoadBalancer"},
		{"2600:9000:24eb:dc00:1:3b80:4f00:21", "LoadBalancer"},
		{"x2600:9000:24eb:XYZ1:1:3b80:4f00:21", "LoadBalancer"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			matchedELB := elbp.SearchResources(&td.ipAddr)
			matchedELBType := reflect.TypeOf(*matchedELB)

			if matchedELBType.Name() != td.expectedType {
				t.Errorf("ELB search failed; expected %s after search, received %s", td.expectedType, matchedELBType.Name())
			}
		})
	}
}
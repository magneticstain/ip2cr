package plugin_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	plugin "github.com/magneticstain/ip-2-cloudresource/src/plugin/elb"
)

func elbpFactory() plugin.ELBPlugin {
	ac, _ := awsconnector.New()

	elbp := plugin.NewELBPlugin(ac)

	return elbp
}

func TestGetResources(t *testing.T) {
	elbp := elbpFactory()

	elbResources, _ := elbp.GetResources()

	expectedType := "Resource"
	for _, elb := range elbResources {
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
		{"1.1.1.1", "Resource"},
		{"1234.45.9666.1", "Resource"},
		{"18.161.22.61", "Resource"},
		{"2600:9000:24eb:dc00:1:3b80:4f00:21", "Resource"},
		{"x2600:9000:24eb:XYZ1:1:3b80:4f00:21", "Resource"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			matchedELB, _ := elbp.SearchResources(td.ipAddr)
			matchedELBType := reflect.TypeOf(matchedELB)

			if matchedELBType.Name() != td.expectedType {
				t.Errorf("ELB search failed; expected %s after search, received %s", td.expectedType, matchedELBType.Name())
			}
		})
	}
}

func elbv1pFactory() plugin.ELBv1Plugin {
	ac, _ := awsconnector.New()

	elbv1p := plugin.NewELBv1Plugin(ac)

	return elbv1p
}

func TestGetResources_Elbv1(t *testing.T) {
	elbv1p := elbv1pFactory()

	elbResources, _ := elbv1p.GetResources()

	expectedType := "Resource"
	for _, elb := range elbResources {
		elbType := reflect.TypeOf(elb)
		if elbType.Name() != expectedType {
			t.Errorf("Fetching resources via ELBv1 Plugin failed; wanted %s type, received %s", expectedType, elbType.Name())
		}
	}
}

func TestSearchResources_Elbv1(t *testing.T) {
	elbv1p := elbv1pFactory()

	var tests = []struct {
		ipAddr, expectedType string
	}{
		{"1.1.1.1", "Resource"},
		{"1234.45.9666.1", "Resource"},
		{"18.161.22.61", "Resource"},
		{"2600:9000:24eb:dc00:1:3b80:4f00:21", "Resource"},
		{"x2600:9000:24eb:XYZ1:1:3b80:4f00:21", "Resource"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			matchedELB, _ := elbv1p.SearchResources(td.ipAddr)
			matchedELBType := reflect.TypeOf(matchedELB)

			if matchedELBType.Name() != td.expectedType {
				t.Errorf("ELBv1 search failed; expected %s after search, received %s", td.expectedType, matchedELBType.Name())
			}
		})
	}
}

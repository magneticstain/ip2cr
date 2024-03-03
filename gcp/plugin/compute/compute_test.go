package compute_test

import (
	"reflect"
	"testing"

	plugin "github.com/magneticstain/ip-2-cloudresource/gcp/plugin/compute"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

func compPlugFactory() plugin.ComputePlugin {
	compPlug := plugin.ComputePlugin{}

	return compPlug
}

func TestCheckComputeIP_ValidIPv4s(t *testing.T) {
	var tests = []struct {
		ipAddr, tgtIp string
		match         bool
	}{
		{"1.1.1.1", "1.1.1.1", true},
		{"18.161.22.61", "18.161.22.61", true},
		{"192.168.1.1", "192.168.1.1", true},
		{"22.642.22.19", "2222.642.22.1999", false},
	}

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			var computeResource, matchingResource generalResource.Resource

			computeResource.PublicIPv4Addrs = []string{td.ipAddr}

			_, found := plugin.CheckComputeIP(&computeResource, &matchingResource, td.tgtIp, 4)

			if found != td.match {
				t.Errorf("Running Compute IP check for IPv4 addresses failed; IP: %s, Found: %t, Should Be Found?: %t", td.tgtIp, found, td.match)
			}
		})
	}
}

func TestCheckComputeIP_IPv6(t *testing.T) {
	var tests = []struct {
		ipAddr, tgtIp string
		match         bool
	}{
		{"2600:9000:24eb:dc00:1:3b80:4f00:21", "2600:9000:24eb:dc00:1:3b80:4f00:21", true},
		{"2001:db8:3333:4444:5555:6666:7777:8888", "2001:db8:3333:4444:5555:6666:7777:8888", true},
		{"xxxx:yyyy:3333:4444:5555:6666:7777:8888", "1111:2222:3333:4444:5555:6666:7777:8888", false},
	}

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			var computeResource, matchingResource generalResource.Resource

			computeResource.PublicIPv6Addrs = []string{td.ipAddr}

			_, found := plugin.CheckComputeIP(&computeResource, &matchingResource, td.tgtIp, 6)

			if found != td.match {
				t.Errorf("Running Compute IP check for IPv6 addresses failed; IP: %s, Found: %t, Should Be Found?: %t", td.tgtIp, found, td.match)
			}
		})
	}
}

func TestCheckComputeIP_InvalidIPs(t *testing.T) {
	var tests = []struct {
		ipAddr, tgtIp string
		match         bool
	}{
		{"1234.56.78.910", "1234.56.78.910", true},
		{"1111:2222:3333:4444:5555:6666:7777:8888", "xxxx:yyyy:3333:4444:5555:6666:7777:8888", false},
	}

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			var computeResource, matchingResource generalResource.Resource

			computeResource.PublicIPv6Addrs = []string{td.ipAddr}

			_, found := plugin.CheckComputeIP(&computeResource, &matchingResource, td.tgtIp, 6)

			if found != td.match {
				t.Errorf("Running Compute IP check for IPv6 addresses failed; IP: %s, Found: %t, Should Be Found?: %t", td.tgtIp, found, td.match)
			}
		})
	}
}

func TestGetResources(t *testing.T) {
	computePlug := compPlugFactory()

	computeResources, _ := computePlug.GetResources()

	expectedType := "Resource"
	for _, resource := range computeResources {
		resourceType := reflect.TypeOf(resource)
		if resourceType.Name() != expectedType {
			t.Errorf("Fetching resources via GCP Compute Plugin failed; wanted %s type, received %s", expectedType, resourceType.Name())
		}
	}
}

func TestSearchResources(t *testing.T) {
	compPlug := compPlugFactory()

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
			matchedInstance, _ := compPlug.SearchResources(td.ipAddr, &matchingResource)
			matchedInstanceType := reflect.TypeOf(matchedInstance)

			if matchedInstanceType.Name() != td.expectedType {
				t.Errorf("GCP Compute search failed; expected %s after search, received %s", td.expectedType, matchedInstanceType.Name())
			}
		})
	}
}

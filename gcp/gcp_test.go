package plugin_test

import (
	"fmt"
	"reflect"
	"testing"

	gcpcontroller "github.com/magneticstain/ip-2-cloudresource/gcp"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

func gcpControllerFactory() gcpcontroller.GCPController {
	gcpc := gcpcontroller.GCPController{}

	return gcpc
}

func TestSearchGCPSvc(t *testing.T) {
	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"compute", "1.1.1.1"},
		{"load_balancing", "1.1.1.1"},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		ac := gcpControllerFactory()
		resource := generalResource.Resource{}

		t.Run(testName, func(t *testing.T) {
			res, _ := ac.SearchGCPSvc("", td.ipAddr, td.cloudSvc, &resource)

			resType := reflect.TypeOf(res)
			expectedType := "Resource"
			if resType.Name() != expectedType {
				t.Errorf("GCP resource search failed; expected %s after search, received %s", expectedType, resType.Name())
			}
		})
	}
}

func TestSearchGCPSvc_UnknownCloudSvc(t *testing.T) {
	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"magic_svc", "1.1.1.1"}, // known invalid use case
		{"compute-", "1.1.1.1"},  // known valid use case with one character addition to make it invalid
		{"iam", "1.1.1.1"},       // valid GCP service, but not one that would ever interact with IP addresses
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		ac := gcpControllerFactory()
		resource := generalResource.Resource{}

		t.Run(testName, func(t *testing.T) {
			_, err := ac.SearchGCPSvc("", td.ipAddr, td.cloudSvc, &resource)
			if err == nil {
				t.Errorf("Error was expected, but not seen, when performing general GCP search; using %s for unknown cloud service name", td.cloudSvc)
			}
		})
	}
}

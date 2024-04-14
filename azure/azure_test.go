package azure_test

import (
	"fmt"
	"reflect"
	"testing"

	azurecontroller "github.com/magneticstain/ip-2-cloudresource/azure"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

func azureControllerFactory() azurecontroller.AzureController {
	azurec := azurecontroller.AzureController{}

	return azurec
}

func TestSearchAzureSvc(t *testing.T) {
	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"virtual_machines", "1.1.1.1"},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		ac := azureControllerFactory()
		resource := generalResource.Resource{}

		t.Run(testName, func(t *testing.T) {
			res, _ := ac.SearchAzureSvc("", td.ipAddr, td.cloudSvc, &resource)

			resType := reflect.TypeOf(res)
			expectedType := "Resource"
			if resType.Name() != expectedType {
				t.Errorf("Azure resource search failed; expected %s after search, received %s", expectedType, resType.Name())
			}
		})
	}
}

func TestSearchAzureSvc_UnknownCloudSvc(t *testing.T) {
	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"magic_svc", "1.1.1.1"},         // known invalid use case
		{"virtual_machines-", "1.1.1.1"}, // known valid use case with one character addition to make it invalid
		{"iam", "1.1.1.1"},               // valid Azure service, but not one that would ever interact with IP addresses
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		ac := azureControllerFactory()
		resource := generalResource.Resource{}

		t.Run(testName, func(t *testing.T) {
			_, err := ac.SearchAzureSvc("", td.ipAddr, td.cloudSvc, &resource)
			if err == nil {
				t.Errorf("Error was expected, but not seen, when performing general Azure search; using %s for unknown cloud service name", td.cloudSvc)
			}
		})
	}
}

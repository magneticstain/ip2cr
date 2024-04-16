package publicip_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	az_public_ip "github.com/magneticstain/ip-2-cloudresource/azure/public_ip"
)

// func azpipaPlugFactory() az_public_ip.AzPublicIPAddr {
// 	azpipa := az_public_ip.AzPublicIPAddr{}

// 	return azpipa
// }

func TestGetPublicIPAddressProperties(t *testing.T) {
	azCreds := azidentity.DefaultAzureCredential{}

	azPubIPMockID := "this_needs_to_be_set_or_theres_a_panic"
	azPubIPObj := armnetwork.PublicIPAddress{ID: &azPubIPMockID}

	azpipaProp, _ := az_public_ip.GetPublicIPAddressProperties(&azCreds, &azPubIPObj, context.Background())

	expectedType := "PublicIPAddressesClientGetResponse"
	resourceType := reflect.TypeOf(azpipaProp)
	if resourceType.Name() != expectedType {
		t.Errorf("Fetching Azure Public IP Address properties failed; wanted %s type, received %s", expectedType, resourceType.Name())
	}
}

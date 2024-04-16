package publicip

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	log "github.com/sirupsen/logrus"
)

type AzPublicIPAddr struct{}

func GetPublicIPAddressProperties(azureConn *azidentity.DefaultAzureCredential, publicIpData *armnetwork.PublicIPAddress, ctx context.Context) (armnetwork.PublicIPAddressesClientGetResponse, error) {
	var pubIpAddrClient *armnetwork.PublicIPAddressesClient
	var publicIpAddrProps armnetwork.PublicIPAddressesClientGetResponse
	var err error

	log.Debug("fetching public IP address properties directly for ", *publicIpData.ID)

	parsedPublicIpId, err := arm.ParseResourceID(*publicIpData.ID)
	if err != nil {
		return publicIpAddrProps, err
	}

	pubIpAddrClient, err = armnetwork.NewPublicIPAddressesClient(parsedPublicIpId.SubscriptionID, azureConn, nil)
	if err != nil {
		return publicIpAddrProps, err
	}

	// yes, it actually does match by name vs ID
	publicIpAddrProps, err = pubIpAddrClient.Get(ctx, parsedPublicIpId.ResourceGroupName, parsedPublicIpId.Name, nil)

	return publicIpAddrProps, err
}

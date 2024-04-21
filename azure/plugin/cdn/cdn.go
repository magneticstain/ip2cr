package cdn

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn"
	log "github.com/sirupsen/logrus"

	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
	"github.com/magneticstain/ip-2-cloudresource/utils"
)

type AzCDNPlugin struct {
	AzureConn      azidentity.DefaultAzureCredential
	SubscriptionID string
}

func (azcdnp *AzCDNPlugin) GetResources() ([]generalResource.Resource, error) {
	var cdnResources []generalResource.Resource
	var currentResource generalResource.Resource
	var cdnEndpointID, cdnEndpointName *string
	var cdnEndpointStatus string

	afdClientFactory, err := armcdn.NewClientFactory(azcdnp.SubscriptionID, &azcdnp.AzureConn, nil)
	if err != nil {
		return cdnResources, err
	}

	ctx := context.Background()

	// traverse CDNM profiles first
	cdnProfilePager := afdClientFactory.NewProfilesClient().NewListPager(nil)
	for cdnProfilePager.More() {
		nextCDNProfileSet, err := cdnProfilePager.NextPage(ctx)
		if err != nil {
			return cdnResources, err
		}
		cdnProfileSet := nextCDNProfileSet.Value
		log.Debug("found [ ", len(cdnProfileSet), " ] Azure Front Door CDN profiles")

		for _, cdnProfile := range cdnProfileSet {
			parsedProfileData, err := arm.ParseResourceID(*cdnProfile.ID)
			if err != nil {
				return cdnResources, err
			}

			// second, list all endpoints for each profile
			cdnEndpointPager := afdClientFactory.NewAFDEndpointsClient().NewListByProfilePager(parsedProfileData.ResourceGroupName, *cdnProfile.Name, nil)
			for cdnEndpointPager.More() {
				nextCDNEndpointSet, err := cdnEndpointPager.NextPage(ctx)
				if err != nil {
					return cdnResources, err
				}
				cdnEndpointSet := nextCDNEndpointSet.Value
				log.Debug("found [ ", len(cdnEndpointSet), " ] Azure Front Door CDN endpoints")

				for _, cdnEndpoint := range cdnEndpointSet {
					cdnEndpointID = cdnEndpoint.ID
					cdnEndpointName = cdnEndpoint.Name
					cdnEndpointStatus = string(*cdnEndpoint.Properties.EnabledState)

					log.Debug("Azure Front Door CDN endpoint found - ID: ", *cdnEndpointID, ", Name: ", *cdnEndpointName, ", Status: ", cdnEndpointStatus)

					var publicIPv4Addrs, publicIPv6Addrs []string
					cdnFQDN := cdnEndpoint.Properties.HostName

					pubIPAddrData, err := utils.LookupFQDN(*cdnFQDN)
					if err != nil {
						return cdnResources, err
					}

					for _, ipAddr := range pubIPAddrData {
						ipVer, err := utils.DetermineIpAddrVersion(ipAddr.String())
						if err != nil {
							return cdnResources, err
						}

						if ipVer == 4 {
							publicIPv4Addrs = append(publicIPv4Addrs, ipAddr.String())
						} else {
							publicIPv6Addrs = append(publicIPv6Addrs, ipAddr.String())
						}
					}

					currentResource = generalResource.Resource{
						Id:              *cdnEndpointID,
						RID:             *cdnEndpointID,
						AccountID:       azcdnp.SubscriptionID,
						Name:            *cdnEndpointName,
						Status:          cdnEndpointStatus,
						CloudSvc:        "cdn",
						PublicIPv4Addrs: publicIPv4Addrs,
						PublicIPv6Addrs: publicIPv6Addrs,
					}

					cdnResources = append(
						cdnResources,
						currentResource,
					)
				}
			}
		}
	}

	return cdnResources, nil
}

func (azcdnp AzCDNPlugin) SearchResources(tgtIP string, matchingResource *generalResource.Resource) (*generalResource.Resource, error) {
	log.Debug("fetching and searching Azure Front Door CDN resources")

	fetchedResources, err := azcdnp.GetResources()
	if err != nil {
		return matchingResource, err
	}

	for _, cdnResource := range fetchedResources {
		for _, ipAddr := range cdnResource.PublicIPv4Addrs {
			if ipAddr == tgtIP {
				matchingResource = &cdnResource

				log.Debug("IP found as Front Door CDN Endpoint -> ", matchingResource.RID)

				break
			}
		}
	}

	return matchingResource, nil
}

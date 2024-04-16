package load_balancer

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	log "github.com/sirupsen/logrus"

	az_public_ip "github.com/magneticstain/ip-2-cloudresource/azure/public_ip"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type AzLoadBalancerPlugin struct {
	AzureConn      azidentity.DefaultAzureCredential
	SubscriptionID string
}

func (azlbp *AzLoadBalancerPlugin) GetResources() ([]generalResource.Resource, error) {
	var lbResources []generalResource.Resource
	var currentResource generalResource.Resource
	var lbID, lbName *string
	var lbStatus string

	lbClient, err := armnetwork.NewLoadBalancersClient(azlbp.SubscriptionID, &azlbp.AzureConn, nil)
	if err != nil {
		return lbResources, err
	}

	ctx := context.Background()
	lbPager := lbClient.NewListAllPager(nil)
	for lbPager.More() {
		nextLbSet, err := lbPager.NextPage(ctx)
		if err != nil {
			return lbResources, err
		}
		lbSet := nextLbSet.Value
		log.Debug("found [ ", len(lbSet), " ] Azure load balancers")

		for _, azlb := range lbSet {
			lbID = azlb.ID
			lbName = azlb.Name
			lbStatus = string(*azlb.Properties.ProvisioningState)
			if err != nil {
				return lbResources, err
			}

			log.Debug("Azure Load Balancer found - ID: ", *lbID, ", Name: ", *lbName, ", Status: ", lbStatus)

			var publicIPv4Addrs []string
			for _, lb_frontend_config := range azlb.Properties.FrontendIPConfigurations {
				pubIPAddrData := lb_frontend_config.Properties.PublicIPAddress

				publicIP, err := az_public_ip.GetPublicIPAddressProperties(&azlbp.AzureConn, pubIPAddrData, ctx)
				if err != nil {
					return lbResources, err
				}

				publicIPv4Addrs = append(publicIPv4Addrs, *publicIP.Properties.IPAddress)
			}

			currentResource = generalResource.Resource{
				Id:              *lbID,
				RID:             *lbID,
				AccountID:       azlbp.SubscriptionID,
				Name:            *lbName,
				Status:          lbStatus,
				CloudSvc:        "virtual_machines",
				PublicIPv4Addrs: publicIPv4Addrs,
			}

			lbResources = append(
				lbResources,
				currentResource,
			)
		}
	}

	return lbResources, nil
}

func (azlbp AzLoadBalancerPlugin) SearchResources(tgtIP string, matchingResource *generalResource.Resource) (*generalResource.Resource, error) {
	log.Debug("fetching and searching Azure load balancer resources")

	fetchedResources, err := azlbp.GetResources()
	if err != nil {
		return matchingResource, err
	}

	for _, lbResource := range fetchedResources {
		for _, ipAddr := range lbResource.PublicIPv4Addrs {
			if ipAddr == tgtIP {
				matchingResource = &lbResource

				log.Debug("IP found as Load Balancer -> ", matchingResource.RID)

				break
			}
		}
	}

	return matchingResource, nil
}

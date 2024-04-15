package virtual_machines

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	log "github.com/sirupsen/logrus"

	az_public_ip "github.com/magneticstain/ip-2-cloudresource/azure/public_ip"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type AzVirtualMachinePlugin struct {
	AzureConn      azidentity.DefaultAzureCredential
	SubscriptionID string
}

func GetVMStatus(vmClient *armcompute.VirtualMachinesClient, vm *armcompute.VirtualMachine, ctx context.Context) (*string, error) {
	var vmStatus *string

	// DEV NOTE: if you're wondering why vm.Properties.InstanceView is always nil and we need to make two API calls here, this is why: https://github.com/Azure/azure-sdk-for-go/issues/4828

	vmIDData, err := arm.ParseResourceID(*vm.ID)
	if err != nil {
		return vmStatus, err
	}

	instView, err := vmClient.InstanceView(ctx, vmIDData.ResourceGroupName, *vm.Name, nil)
	if err != nil {
		return vmStatus, err
	}

	// we only care about the latest status
	vmStatus = instView.Statuses[len(instView.Statuses)-1].DisplayStatus

	return vmStatus, nil
}

func (azvmp *AzVirtualMachinePlugin) GatherVMPublicIPAddrData(vmInstance *armcompute.VirtualMachine, IpVer armnetwork.IPVersion, ctx context.Context) ([]string, error) {
	var publicIPAddrs []string
	var publicIpVer *armnetwork.IPVersion
	var ipAddr *string

	nicClient, err := armnetwork.NewInterfacesClient(azvmp.SubscriptionID, &azvmp.AzureConn, nil)
	if err != nil {
		return publicIPAddrs, err
	}

	vmNics := vmInstance.Properties.NetworkProfile.NetworkInterfaces
	for _, nicRef := range vmNics {
		// to get the public IP, we need to get the details of each NIC, traverse each IP config, check if it has the PublicIPAddress data present, and if so, grab it from the included properties. If no properties exist (because reasons I guess), then we need to fetch the PublicIPAddress data with the given ID, then parse the data from that

		parsedNicId, err := arm.ParseResourceID(*nicRef.ID)
		if err != nil {
			return publicIPAddrs, err
		}

		nicData, err := nicClient.Get(ctx, parsedNicId.ResourceGroupName, parsedNicId.Name, nil)
		if err != nil {
			return publicIPAddrs, err
		}

		nicIpConfigs := nicData.Properties.IPConfigurations
		for _, nicIpAddrConfig := range nicIpConfigs {
			if nicIpAddrConfig.Properties.PublicIPAddress != nil {
				publicIpData := nicIpAddrConfig.Properties.PublicIPAddress

				if publicIpData.Properties != nil {
					publicIpVer = publicIpData.Properties.PublicIPAddressVersion
					ipAddr = publicIpData.Properties.IPAddress
				} else {
					// for some reason, some (and _only_ some) API calls don't have this info, idk....
					// instead, we need to make _another_ API request to get the public IP data

					log.Debug("public IP address properties not included with ", *vmInstance.ID, " / ", *vmInstance.Name, " public IP address data (this is common) - fetching it more directly from Azure...")

					publicIpProps, err := az_public_ip.GetPublicIPAddressProperties(&azvmp.AzureConn, publicIpData, ctx)
					if err != nil {
						return publicIPAddrs, err
					}

					publicIpVer = publicIpProps.Properties.PublicIPAddressVersion
					ipAddr = publicIpProps.Properties.IPAddress
				}

				if *publicIpVer == IpVer {
					publicIPAddrs = append(publicIPAddrs, *ipAddr)
				}
			} else {
				log.Debug("no public ", IpVer, " address found for NIC [ ", nicData.Name, " ] on VM [ ", vmInstance.Name, " ] in resource group [ ", parsedNicId.ResourceGroupName, " ]")
			}
		}
	}

	return publicIPAddrs, err
}

func (azvmp *AzVirtualMachinePlugin) GetResources() ([]generalResource.Resource, error) {
	var vmResources []generalResource.Resource
	var currentResource generalResource.Resource
	var vmID, vmName *string

	vmClient, err := armcompute.NewVirtualMachinesClient(azvmp.SubscriptionID, &azvmp.AzureConn, nil)
	if err != nil {
		return vmResources, err
	}

	ctx := context.Background()
	vmPager := vmClient.NewListAllPager(nil)
	for vmPager.More() {
		nextVmSet, err := vmPager.NextPage(ctx)
		if err != nil {
			return vmResources, err
		}
		vmSet := nextVmSet.Value
		log.Debug("found [ ", len(vmSet), " ] Azure virtual machines")

		for _, vm := range vmSet {
			vmID = vm.ID
			vmName = vm.Name
			vmStatus, err := GetVMStatus(vmClient, vm, ctx)
			if err != nil {
				return vmResources, err
			}

			log.Debug("Azure VM instance found - ID: ", *vmID, ", Name: ", *vmName, ", Status: ", *vmStatus)

			log.Debug("fetching IPv4 addresses")
			publicIPv4Addrs, err := azvmp.GatherVMPublicIPAddrData(vm, armnetwork.IPVersionIPv4, ctx)
			if err != nil {
				return vmResources, err
			}

			log.Debug("fetching IPv6 addresses")
			publicIPv6Addrs, err := azvmp.GatherVMPublicIPAddrData(vm, armnetwork.IPVersionIPv6, ctx)
			if err != nil {
				return vmResources, err
			}

			currentResource = generalResource.Resource{
				Id:              *vmID,
				RID:             *vmID,
				AccountID:       azvmp.SubscriptionID,
				Name:            *vmName,
				Status:          *vmStatus,
				CloudSvc:        "virtual_machines",
				PublicIPv4Addrs: publicIPv4Addrs,
				PublicIPv6Addrs: publicIPv6Addrs,
			}

			vmResources = append(
				vmResources,
				currentResource,
			)
		}
	}

	return vmResources, nil
}

func (azvmp AzVirtualMachinePlugin) SearchResources(tgtIP string, matchingResource *generalResource.Resource) (*generalResource.Resource, error) {
	log.Debug("fetching and searching Azure virtual machine resources")

	fetchedResources, err := azvmp.GetResources()
	if err != nil {
		return matchingResource, err
	}

	for _, vmResource := range fetchedResources {
		for _, ipAddr := range vmResource.PublicIPv4Addrs {
			if ipAddr == tgtIP {
				matchingResource = &vmResource

				log.Debug("IP found as Virtual Machine -> ", matchingResource.RID)

				break
			}
		}
	}

	return matchingResource, nil
}

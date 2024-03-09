package virtual_machine

import (
	log "github.com/sirupsen/logrus"

	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type AzVirtualMachinePlugin struct {
	ProjectID string
}

func (azvmp AzVirtualMachinePlugin) GetResources() ([]generalResource.Resource, error) {
	var vmResources []generalResource.Resource

	return vmResources, nil
}

func (azvmp AzVirtualMachinePlugin) SearchResources(tgtIP string, matchingResource *generalResource.Resource) (generalResource.Resource, error) {
	log.Debug("fetching and searching virtual machine resources")

	_, err := azvmp.GetResources()
	if err != nil {
		return *matchingResource, err
	}

	// for _, vmResource := range fetchedResources {}

	return *matchingResource, nil
}

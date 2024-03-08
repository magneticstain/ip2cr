package plugin

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type AzureController struct{}

func GetSupportedSvcs() []string {
	return []string{
		"virtual_machines",
	}
}

func (azctrlr *AzureController) SearchAzureSvc(projectID, ipAddr, cloudSvc string, matchingResource *generalResource.Resource) (generalResource.Resource, error) {
	// var err error

	log.Debug("searching ", cloudSvc, " in Azure controller")

	switch cloudSvc {
	case "virtual_machines":
		// comp := compute.ComputePlugin{
		// 	ProjectID: projectID,
		// }
		// _, err = comp.SearchResources(ipAddr, matchingResource)
		// if err != nil {
		// 	return *matchingResource, err
		// }
	default:
		msg := fmt.Sprintf("unknown Azure service provided: '%s'", cloudSvc)

		return *matchingResource, errors.New(msg)
	}

	return *matchingResource, nil
}

package plugin

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/magneticstain/ip-2-cloudresource/gcp/plugin/compute"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type GCPController struct{}

func GetSupportedSvcs() []string {
	return []string{
		"compute",
	}
}

func (gcpctrlr *GCPController) SearchGCPSvc(projectID, ipAddr, cloudSvc string, matchingResource *generalResource.Resource) (generalResource.Resource, error) {
	var err error

	log.Debug("searching ", cloudSvc, " in GCP controller")

	switch cloudSvc {
	case "compute":
		comp := compute.ComputePlugin{
			ProjectID: projectID,
		}
		_, err = comp.SearchResources(ipAddr, matchingResource)
		if err != nil {
			return *matchingResource, err
		}
	default:
		msg := fmt.Sprintf("unknown GCP service provided: '%s'", cloudSvc)

		return *matchingResource, errors.New(msg)
	}

	return *matchingResource, nil
}

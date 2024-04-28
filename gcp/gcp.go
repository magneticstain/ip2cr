package gcp

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/magneticstain/ip-2-cloudresource/gcp/plugin/cloud_sql"
	"github.com/magneticstain/ip-2-cloudresource/gcp/plugin/compute"
	"github.com/magneticstain/ip-2-cloudresource/gcp/plugin/load_balancing"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type GCPController struct{}

func GetSupportedSvcs() []string {
	return []string{
		"compute",
		"load_balancing",
		"cloud_sql",
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
	case "load_balancing":
		lbp := load_balancing.LoadBalancingPlugin{
			ProjectID: projectID,
		}
		_, err = lbp.SearchResources(ipAddr, matchingResource)
		if err != nil {
			return *matchingResource, err
		}
	case "cloud_sql":
		csqlp := cloud_sql.CloudSQLPlugin{
			ProjectID: projectID,
		}
		_, err = csqlp.SearchResources(ipAddr, matchingResource)
		if err != nil {
			return *matchingResource, err
		}
	default:
		msg := fmt.Sprintf("unknown GCP service provided: '%s'", cloudSvc)

		return *matchingResource, errors.New(msg)
	}

	return *matchingResource, nil
}

package cloud_sql

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type CloudSQLPlugin struct {
	ProjectID string
}

func (lbp CloudSQLPlugin) GetResources() ([]generalResource.Resource, error) {
	var lbResources []generalResource.Resource

	return lbResources, nil
}

func (lbp CloudSQLPlugin) SearchResources(tgtIP string, matchingResource *generalResource.Resource) (generalResource.Resource, error) {
	log.Debug("fetching and searching cloudsql resources")

	fetchedResources, err := lbp.GetResources()
	if err != nil {
		return *matchingResource, err
	}

	for _, csqlResource := range fetchedResources {
		for _, ipv4Addr := range csqlResource.PublicIPv4Addrs {
			if ipv4Addr == tgtIP {
				matchingResource.RID = fmt.Sprintf("%s/%s", csqlResource.Id, csqlResource.Name)
				matchingResource.CloudSvc = "cloud_sql"

				break
			}
		}

		for _, ipv6Addr := range csqlResource.PublicIPv6Addrs {
			if ipv6Addr == tgtIP {
				matchingResource.RID = csqlResource.Id
				matchingResource.CloudSvc = "cloud_sql"

				break
			}
		}

		if matchingResource.RID != "" {
			log.Debug("IP found as CloudSQL instance -> ", matchingResource.RID)

			break
		}
	}

	return *matchingResource, nil
}

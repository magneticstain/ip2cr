package cloud_sql

/*
DEV NOTE:
---
As of 03/2024, Google does not have a client library for cloudsql, and instead suggests using the sqladmin package (despite what the summary for the sqladmin package says). This packages appears to be legacy, but apparently isn't, idk....

Either way, if they ever come out with a client library for CloudSQL, it would be best to migrate this package to that instead.

https://pkg.go.dev/google.golang.org/api@v0.167.0/sqladmin/v1
https://cloud.google.com/sql/docs/mysql/admin-api/libraries
*/

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"google.golang.org/api/sqladmin/v1"

	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
	"github.com/magneticstain/ip-2-cloudresource/utils"
)

type CloudSQLPlugin struct {
	ProjectID string
}

func (csqlp CloudSQLPlugin) GetResources() ([]generalResource.Resource, error) {
	var csqlResources []generalResource.Resource

	ctx := context.Background()

	sqlAdminSvc, err := sqladmin.NewService(ctx)
	if err != nil {
		return csqlResources, err
	}

	csqlInstListCall := sqlAdminSvc.Instances.List(csqlp.ProjectID)
	csqlInstListResp, err := csqlInstListCall.Do()
	if err != nil {
		return csqlResources, err
	}

	for _, csqlInstance := range csqlInstListResp.Items {
		instanceName := csqlInstance.Name
		instanceStatus := csqlInstance.State
		instanceIpAddrs := csqlInstance.IpAddresses

		log.Debug("load balancer endpoint found - Name: ", instanceName, ", Status: ", instanceStatus)

		currentResource := generalResource.Resource{
			Name:           instanceName,
			Status:         instanceStatus,
			CloudSvc:       "cloudsql",
			AccountAliases: []string{csqlp.ProjectID},
		}

		for _, ipAddrMap := range instanceIpAddrs {
			ipAddr := ipAddrMap.IpAddress

			ipVer, err := utils.DetermineIpAddrVersion(ipAddr)
			if err != nil {
				return csqlResources, err
			}

			switch ipVer {
			case 4:
				currentResource.PublicIPv4Addrs = append(currentResource.PublicIPv4Addrs, ipAddr)
			case 6:
				currentResource.PublicIPv6Addrs = append(currentResource.PublicIPv6Addrs, ipAddr)
			default:
				return csqlResources, fmt.Errorf("invalid IP version found for GCP CloudSQL instance; IP: %s, Version: IPv%d", ipAddr, ipVer)
			}
		}

		csqlResources = append(
			csqlResources,
			currentResource,
		)
	}

	return csqlResources, nil
}

func (csqlp CloudSQLPlugin) SearchResources(tgtIP string, matchingResource *generalResource.Resource) (generalResource.Resource, error) {
	log.Debug("fetching and searching cloudsql resources")

	fetchedResources, err := csqlp.GetResources()
	if err != nil {
		return *matchingResource, err
	}

	for _, csqlResource := range fetchedResources {
		ridSlug := fmt.Sprintf("%s/%s", csqlResource.Id, csqlResource.Name)

		for _, ipv4Addr := range csqlResource.PublicIPv4Addrs {
			if ipv4Addr == tgtIP {
				matchingResource.RID = ridSlug
				matchingResource.CloudSvc = "cloud_sql"

				break
			}
		}

		for _, ipv6Addr := range csqlResource.PublicIPv6Addrs {
			if ipv6Addr == tgtIP {
				matchingResource.RID = ridSlug
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

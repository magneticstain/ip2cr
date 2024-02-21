package search

import (
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type Search struct {
	IpAddr  string
}

func (search Search) InitSearch(cloudSvc string) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource

	return matchingResource, nil
}

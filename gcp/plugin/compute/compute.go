package compute

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"

	gcpcomputeapi "cloud.google.com/go/compute/apiv1"
	gcpcomputepbapi "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"

	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type ComputePlugin struct {
	ProjectID string
}

func (comp ComputePlugin) GetResources() ([]ComputeResource, error) {
	var computeClient *gcpcomputeapi.InstancesClient
	var instanceList *gcpcomputeapi.InstancesScopedListPairIterator
	var computeResources []ComputeResource

	// REF: https://cloud.google.com/compute/docs/samples/compute-instances-list-all#compute_instances_list_all-go
	ctx := context.Background()

	computeClient, err := gcpcomputeapi.NewInstancesRESTClient(ctx)
	if err != nil {
		return computeResources, err
	}
	defer computeClient.Close()

	req := &gcpcomputepbapi.AggregatedListInstancesRequest{
		Project: comp.ProjectID,
	}

	// normally, we would just return the iterator and allow the search function to iterate through each
	// however, for some reason, passing the iterator back to the search function results in a memory error when trying to read it
	// if someone in the future figures out why, feel free to refactor/patch it
	instanceList = computeClient.AggregatedList(ctx, req)

	for {
		instanceListPair, err := instanceList.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return computeResources, err
		}

		instances := instanceListPair.Value.Instances
		for _, instance := range instances {
			instanceId := strconv.FormatUint(instance.GetId(), 10)
			instanceName := instance.GetName()
			instanceStatus := instance.GetStatus()

			log.Debug("compute instance found - ID: ", instanceId, ", Name: ", instanceName, ", Status: ", instanceStatus)

			currentResource := ComputeResource{
				Id:     instanceId,
				Name:   instanceName,
				Status: instanceStatus,
			}
			computeResources = append(
				computeResources,
				currentResource,
			)
		}
	}

	return computeResources, nil
}

func (comp ComputePlugin) SearchResources(tgtIP string) (generalResource.Resource, error) {
	// var computeResources gcpcomputeapi.InstancesScopedListPairIterator
	var matchingResource generalResource.Resource

	log.Debug("fetching and searching compute resources")

	_, err := comp.GetResources()
	if err != nil {
		return matchingResource, err
	}

	return matchingResource, nil
}

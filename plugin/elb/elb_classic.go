package plugin

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws_connector"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
	"github.com/magneticstain/ip-2-cloudresource/utils"
)

type ELBv1Plugin struct {
	AwsConn        awsconnector.AWSConnector
	NetworkMapping bool
}

func (elbv1p ELBv1Plugin) GetResources() ([]types.LoadBalancerDescription, error) {
	var elbs []types.LoadBalancerDescription

	elbClient := elasticloadbalancing.NewFromConfig(elbv1p.AwsConn.AwsConfig)
	paginator := elasticloadbalancing.NewDescribeLoadBalancersPaginator(elbClient, nil)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return elbs, err
		}

		elbs = append(elbs, output.LoadBalancerDescriptions...)
	}

	return elbs, nil
}

func (elbv1p ELBv1Plugin) SearchResources(tgtIP string) (generalResource.Resource, error) {
	var elbIPAddrs []net.IP
	var matchingResource generalResource.Resource

	elbResources, err := elbv1p.GetResources()
	if err != nil {
		return matchingResource, err
	}

	for _, elb := range elbResources {
		elbIPAddrs, err = utils.LookupFQDN(*elb.DNSName)
		if err != nil {
			return matchingResource, err
		}

		for _, ipAddr := range elbIPAddrs {
			if ipAddr.String() == tgtIP {
				matchingResource.RID = *elb.LoadBalancerName
				matchingResource.CloudSvc = "elbv1"

				if elbv1p.NetworkMapping {
					matchingResource.NetworkMap = append(matchingResource.NetworkMap, *elb.DNSName, *elb.CanonicalHostedZoneNameID, *elb.VPCId, utils.FormatStrSliceAsCSV(elb.AvailabilityZones), utils.FormatStrSliceAsCSV(elb.Subnets))
				}

				log.Debug("IP found as Classic Elastic Load Balancer -> ", matchingResource.RID, " with network info ", matchingResource.NetworkMap)

				break
			}
		}
	}

	return matchingResource, nil
}

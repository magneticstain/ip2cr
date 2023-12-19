package plugin

import (
	"context"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	generalResource "github.com/magneticstain/ip-2-cloudresource/src/resource"
	"github.com/magneticstain/ip-2-cloudresource/src/utils"
)

type ELBPlugin struct {
	AwsConn awsconnector.AWSConnector
}

func (elbp ELBPlugin) GetElbListeners(elbArn string) ([]types.Listener, error) {
	var listeners []types.Listener

	elb_client := elasticloadbalancingv2.NewFromConfig(elbp.AwsConn.AwsConfig)
	paginator := elasticloadbalancingv2.NewDescribeListenersPaginator(elb_client, &elasticloadbalancingv2.DescribeListenersInput{
		LoadBalancerArn: &elbArn,
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return listeners, err
		}

		listeners = append(listeners, output.Listeners...)
	}

	return listeners, nil
}

func (elbp ELBPlugin) GetElbTgts(elbListeners []types.Listener) ([]ELBTarget, error) {
	var elbTgt ELBTarget
	var elbTgts []ELBTarget

	elb_client := elasticloadbalancingv2.NewFromConfig(elbp.AwsConn.AwsConfig)
	for _, listener := range elbListeners {
		elbTgt = ELBTarget{
			ListenerArn: *listener.ListenerArn,
		}

		for _, listnerAction := range listener.DefaultActions {
			elbTgt.TgtGrpArn = *listnerAction.TargetGroupArn

			resp, err := elb_client.DescribeTargetHealth(context.TODO(), &elasticloadbalancingv2.DescribeTargetHealthInput{
				TargetGroupArn: &elbTgt.TgtGrpArn,
			})
			if err != nil {
				return elbTgts, err
			}

			for _, targetHealth := range resp.TargetHealthDescriptions {
				if targetHealth.Target.Id != nil {
					elbTgt.TgtIds = append(elbTgt.TgtIds, *targetHealth.Target.Id)
				}
			}
		}

		elbTgts = append(elbTgts, elbTgt)
	}

	return elbTgts, nil
}

func AddElbAZDataToNetworkMap(matchingResource *generalResource.Resource, AZData []types.AvailabilityZone) {
	var AZSlug string
	var AZDataSet []string

	for _, AZ := range AZData {
		AZSlug = fmt.Sprintf("%s (%s)", *AZ.SubnetId, *AZ.ZoneName)

		AZDataSet = append(AZDataSet, AZSlug)
	}

	matchingResource.NetworkMap = append(matchingResource.NetworkMap, utils.FormatStrSliceAsCSV(AZDataSet))
}

func (elbp ELBPlugin) GetResources() ([]types.LoadBalancer, error) {
	var elbs []types.LoadBalancer

	elb_client := elasticloadbalancingv2.NewFromConfig(elbp.AwsConn.AwsConfig)
	paginator := elasticloadbalancingv2.NewDescribeLoadBalancersPaginator(elb_client, nil)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return elbs, err
		}

		elbs = append(elbs, output.LoadBalancers...)
	}

	return elbs, nil
}

func (elbp ELBPlugin) SearchResources(tgtIP string) (generalResource.Resource, error) {
	var elbIPAddrs []net.IP
	var matchingResource generalResource.Resource
	var elbListners []types.Listener
	var elbTgts []ELBTarget

	elbResources, err := elbp.GetResources()
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
				matchingResource.RID = *elb.LoadBalancerArn

				matchingResource.NetworkMap = append(matchingResource.NetworkMap, *elb.DNSName, *elb.CanonicalHostedZoneId)

				AddElbAZDataToNetworkMap(&matchingResource, elb.AvailabilityZones)

				elbListners, err = elbp.GetElbListeners(*elb.LoadBalancerArn)
				if err != nil {
					return matchingResource, err
				}
				elbTgts, err = elbp.GetElbTgts(elbListners)
				if err != nil {
					return matchingResource, err
				}

				var tgtSlug []string
				for _, tgt := range elbTgts {
					tgtSlug = append(tgtSlug, tgt.ListenerArn, tgt.TgtGrpArn)
					tgtSlug = append(tgtSlug, tgt.TgtIds...)
				}
				matchingResource.NetworkMap = append(matchingResource.NetworkMap, tgtSlug...)

				log.Debug("IP found as Elastic Load Balancer -> ", matchingResource.RID, " with network info ", matchingResource.NetworkMap)

				break
			}
		}
	}

	return matchingResource, nil
}

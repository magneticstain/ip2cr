package plugin

import (
	"context"
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

func NewELBPlugin(awsConn awsconnector.AWSConnector) ELBPlugin {
	elbp := ELBPlugin{AwsConn: awsConn}

	return elbp
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

				log.Debug("IP found as Elastic Load Balancer -> ", matchingResource.RID)

				break
			}
		}
	}

	return matchingResource, nil
}

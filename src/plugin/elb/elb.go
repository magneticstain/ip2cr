package elb

import (
	"context"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	generalPlugin "github.com/magneticstain/ip-2-cloudresource/src/plugin"
	"github.com/magneticstain/ip-2-cloudresource/src/utils"
)

type ELBPlugin struct {
	GenPlugin *generalPlugin.Plugin
	AwsConn   awsconnector.AWSConnector
}

func NewELBPlugin(aws_conn *awsconnector.AWSConnector) ELBPlugin {
	elbp := ELBPlugin{GenPlugin: &generalPlugin.Plugin{}, AwsConn: *aws_conn}

	return elbp
}

func (elbp ELBPlugin) GetResources() (*[]types.LoadBalancer, error) {
	var elbs []types.LoadBalancer

	elb_client := elasticloadbalancingv2.NewFromConfig(elbp.AwsConn.AwsConfig)
	paginator := elasticloadbalancingv2.NewDescribeLoadBalancersPaginator(elb_client, nil)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return &elbs, err
		}

		elbs = append(elbs, output.LoadBalancers...)
	}

	return &elbs, nil
}

func (elbp ELBPlugin) SearchResources(tgtIp *string) (*types.LoadBalancer, error) {
	var elbIpAddrs *[]net.IP
	var matchedELB types.LoadBalancer

	elbResources, err := elbp.GetResources()
	if err != nil {
		return &matchedELB, err
	}

	for _, elb := range *elbResources {
		elbIpAddrs, err = utils.LookupFQDN(elb.DNSName)
		if err != nil {
			return &matchedELB, err
		}

		for _, ipAddr := range *elbIpAddrs {
			if ipAddr.String() == *tgtIp {
				matchedELB = elb
				break
			}
		}
	}

	return &matchedELB, nil
}

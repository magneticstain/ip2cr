package elb

import (
	"context"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	generalPlugin "github.com/magneticstain/ip2cr/src/plugin"
	"github.com/magneticstain/ip2cr/src/utils"
)

type ELBv1Plugin struct {
	GenPlugin *generalPlugin.Plugin
	AwsConn   awsconnector.AWSConnector
}

func NewELBv1Plugin(aws_conn *awsconnector.AWSConnector) ELBv1Plugin {
	elbv1p := ELBv1Plugin{GenPlugin: &generalPlugin.Plugin{}, AwsConn: *aws_conn}

	return elbv1p
}

func (elbv1p ELBv1Plugin) GetResources() (*[]types.LoadBalancerDescription, error) {
	var elbs []types.LoadBalancerDescription

	elb_client := elasticloadbalancing.NewFromConfig(elbv1p.AwsConn.AwsConfig)
	paginator := elasticloadbalancing.NewDescribeLoadBalancersPaginator(elb_client, nil)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return &elbs, err
		}

		elbs = append(elbs, output.LoadBalancerDescriptions...)
	}

	return &elbs, nil
}

func (elbv1p ELBv1Plugin) SearchResources(tgt_ip *string) (*types.LoadBalancerDescription, error) {
	var elbIpAddrs *[]net.IP
	var matchedELB types.LoadBalancerDescription

	elbResources, err := elbv1p.GetResources()
	if err != nil {
		return &matchedELB, err
	}

	for _, elb := range *elbResources {
		elbIpAddrs, err = utils.LookupFQDN(elb.DNSName)
		if err != nil {
			return &matchedELB, err
		}

		for _, ipAddr := range *elbIpAddrs {
			if ipAddr.String() == *tgt_ip {
				matchedELB = elb
				break
			}
		}
	}

	return &matchedELB, nil
}

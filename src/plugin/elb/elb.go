package elb

import (
	"context"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	generalPlugin "github.com/magneticstain/ip2cr/src/plugin"
	"github.com/magneticstain/ip2cr/src/utils"
)

type ELBPlugin struct {
	GenPlugin *generalPlugin.Plugin
	AwsConn   awsconnector.AWSConnector
}

func NewELBPlugin(aws_conn *awsconnector.AWSConnector) ELBPlugin {
	elbp := ELBPlugin{GenPlugin: &generalPlugin.Plugin{}, AwsConn: *aws_conn}

	return elbp
}

func (elbp ELBPlugin) GetResources() *[]types.LoadBalancer {
	var elbs []types.LoadBalancer

	elb_client := elasticloadbalancingv2.NewFromConfig(elbp.AwsConn.AwsConfig)
	paginator := elasticloadbalancingv2.NewDescribeLoadBalancersPaginator(elb_client, nil)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Error("error when running ELB search :: [ ", err, " ]")
		}

		elbs = append(elbs, output.LoadBalancers...)
	}

	return &elbs
}

func (elbp ELBPlugin) SearchResources(tgt_ip *string) *types.LoadBalancer {
	var elbIpAddrs *[]net.IP
	var matchedELB types.LoadBalancer

	elbResources := elbp.GetResources()

	for _, elb := range *elbResources {
		elbIpAddrs = utils.LookupFQDN(elb.DNSName)

		for _, ipAddr := range *elbIpAddrs {
			if ipAddr.String() == *tgt_ip {
				matchedELB = elb
			}
		}
	}

	return &matchedELB
}

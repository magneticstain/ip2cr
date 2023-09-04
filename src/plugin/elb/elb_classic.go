package elb

import (
	"context"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	"github.com/magneticstain/ip-2-cloudresource/src/utils"
)

type ELBv1Plugin struct {
	AwsConn awsconnector.AWSConnector
}

func NewELBv1Plugin(awsConn *awsconnector.AWSConnector) ELBv1Plugin {
	elbv1p := ELBv1Plugin{AwsConn: *awsConn}

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

func (elbv1p ELBv1Plugin) SearchResources(tgtIP *string) (*types.LoadBalancerDescription, error) {
	var elbIPAddrs *[]net.IP
	var matchedELB types.LoadBalancerDescription

	elbResources, err := elbv1p.GetResources()
	if err != nil {
		return &matchedELB, err
	}

	for _, elb := range *elbResources {
		elbIPAddrs, err = utils.LookupFQDN(elb.DNSName)
		if err != nil {
			return &matchedELB, err
		}

		for _, ipAddr := range *elbIPAddrs {
			if ipAddr.String() == *tgtIP {
				matchedELB = elb
				break
			}
		}
	}

	return &matchedELB, nil
}

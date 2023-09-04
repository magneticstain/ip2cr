package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
)

type EC2Plugin struct {
	AwsConn   awsconnector.AWSConnector
}

func NewEC2Plugin(aws_conn *awsconnector.AWSConnector) EC2Plugin {
	ec2p := EC2Plugin{AwsConn: *aws_conn}

	return ec2p
}

func (ec2p EC2Plugin) GetResources() (*[]types.Reservation, error) {
	var instances []types.Reservation

	ec2Client := ec2.NewFromConfig(ec2p.AwsConn.AwsConfig)
	paginator := ec2.NewDescribeInstancesPaginator(ec2Client, nil)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return &instances, err
		}

		instances = append(instances, output.Reservations...)
	}

	return &instances, nil
}

func (ec2p EC2Plugin) SearchResources(tgtIp *string) (*types.Instance, error) {
	var matchedInstance types.Instance

	ec2Resources, err := ec2p.GetResources()
	if err != nil {
		return &matchedInstance, err
	}

	for _, ec2Reservation := range *ec2Resources {
		// unpack instances from reservation
		for _, instance := range ec2Reservation.Instances {
			publicIPv4Addr := instance.PublicIpAddress
			IPv6Addr := instance.Ipv6Address

			if *publicIPv4Addr == *tgtIp || *IPv6Addr == *tgtIp {
				matchedInstance = instance
				break
			}
		}
	}

	return &matchedInstance, nil
}

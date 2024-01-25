package plugin

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws_connector"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type EC2Plugin struct {
	AwsConn        awsconnector.AWSConnector
	NetworkMapping bool
}

func (ec2p EC2Plugin) GetResources() ([]types.Reservation, error) {
	var instances []types.Reservation

	ec2Client := ec2.NewFromConfig(ec2p.AwsConn.AwsConfig)
	paginator := ec2.NewDescribeInstancesPaginator(ec2Client, nil)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return instances, err
		}

		instances = append(instances, output.Reservations...)
	}

	return instances, nil
}

func (ec2p EC2Plugin) SearchResources(tgtIP string) (generalResource.Resource, error) {
	var matchingResource generalResource.Resource

	ec2Resources, err := ec2p.GetResources()
	if err != nil {
		return matchingResource, err
	}

	var publicIPv4Addr, IPv6Addr string
	for _, ec2Reservation := range ec2Resources {
		// unpack instances from reservation
		for _, instance := range ec2Reservation.Instances {
			evalAddrPtr := func(addrPtr *string) string {
				addr := ""
				if addrPtr != nil {
					addr = *addrPtr
				}

				return addr
			}
			publicIPv4Addr = evalAddrPtr(instance.PublicIpAddress)
			IPv6Addr = evalAddrPtr(instance.Ipv6Address)

			if publicIPv4Addr == tgtIP || IPv6Addr == tgtIP {
				matchingResource.RID = *instance.InstanceId // for some reason, the EC2 Instance object doesn't contain the ARN of the instance :/
				matchingResource.CloudSvc = "ec2"

				if ec2p.NetworkMapping {
					matchingResource.NetworkMap = append(matchingResource.NetworkMap, *instance.VpcId, *instance.SubnetId, *instance.InstanceId)
				}

				log.Debug("IP found as EC2 instance -> ", matchingResource.RID, " with network info ", matchingResource.NetworkMap)

				break
			}
		}
	}

	return matchingResource, nil
}

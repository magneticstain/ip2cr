package ec2_test

import (
	"reflect"
	"testing"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	"github.com/magneticstain/ip2cr/src/plugin/ec2"
)

func ec2pFactory() ec2.EC2Plugin {
	ac, _ := awsconnector.New()

	ec2p := ec2.NewEC2Plugin(&ac)

	return ec2p
}

func TestGetResources(t *testing.T) {
	ec2p := ec2pFactory()

	ec2Resources, _ := ec2p.GetResources()

	expectedType := "Reservation"
	for _, instance := range *ec2Resources {
		ec2Type := reflect.TypeOf(instance)
		if ec2Type.Name() != expectedType {
			t.Errorf("Fetching resources via EC2 Plugin failed; wanted %s type, received %s", expectedType, ec2Type.Name())
		}
	}
}

func TestSearchResources(t *testing.T) {
	ec2p := ec2pFactory()

	var tests = []struct {
		ipAddr, expectedType string
	}{
		{"1.1.1.1", "Instance"},
		{"1234.45.9666.1", "Instance"},
		{"18.161.22.61", "Instance"},
		{"2600:9000:24eb:dc00:1:3b80:4f00:21", "Instance"},
		{"x2600:9000:24eb:XYZ1:1:3b80:4f00:21", "Instance"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		t.Run(testName, func(t *testing.T) {
			matchedInstance, _ := ec2p.SearchResources(&td.ipAddr)
			matchedInstanceType := reflect.TypeOf(*matchedInstance)

			if matchedInstanceType.Name() != td.expectedType {
				t.Errorf("EC2 search failed; expected %s after search, received %s", td.expectedType, matchedInstanceType.Name())
			}
		})
	}
}

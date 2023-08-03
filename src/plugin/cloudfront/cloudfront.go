package cloudfront

import (
	"context"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	generalPlugin "github.com/magneticstain/ip2cr/src/plugin"
	"github.com/magneticstain/ip2cr/src/utils"
)

type CloudfrontPlugin struct {
	GenPlugin *generalPlugin.Plugin
	AwsConn   awsconnector.AWSConnector
}

func NewCloudfrontPlugin(aws_conn *awsconnector.AWSConnector) CloudfrontPlugin {
	cfp := CloudfrontPlugin{GenPlugin: &generalPlugin.Plugin{}, AwsConn: *aws_conn}

	return cfp
}

func (cfp CloudfrontPlugin) NormalizeCFDistroFQDN(fqdn *string) string {
	// CloudFront currently returns a `.` appended to the fqdn
	// we'll need to get rid of it so that it can be lookup up properly

	return strings.TrimSuffix(*fqdn, ".")
}

func (cfp CloudfrontPlugin) GetResources() (*[]types.DistributionSummary, error) {
	var distros []types.DistributionSummary

	cf_client := cloudfront.NewFromConfig(cfp.AwsConn.AwsConfig)
	paginator := cloudfront.NewListDistributionsPaginator(cf_client, &cloudfront.ListDistributionsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return &distros, err
		}

		distros = append(distros, output.DistributionList.Items...)
	}

	return &distros, nil
}

func (cfp CloudfrontPlugin) SearchResources(tgt_ip *string) (*types.DistributionSummary, error) {
	var cfDistroFQDN string
	var cfIpAddrs *[]net.IP
	var matchedDistro types.DistributionSummary

	cfResources, err := cfp.GetResources()
	if err != nil {
		return &matchedDistro, err
	}

	for _, cfDistro := range *cfResources {
		cfDistroFQDN = cfp.NormalizeCFDistroFQDN(cfDistro.DomainName)
		cfIpAddrs, err = utils.LookupFQDN(&cfDistroFQDN)
		if err != nil {
			return &matchedDistro, err
		}

		for _, ipAddr := range *cfIpAddrs {
			if ipAddr.String() == *tgt_ip {
				matchedDistro = cfDistro
			}
		}
	}

	return &matchedDistro, nil
}

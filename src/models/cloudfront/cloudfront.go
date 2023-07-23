package cloudfront

import (
	"context"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip2cr/src/models/aws_connector"
	"github.com/magneticstain/ip2cr/src/utils"
)

type CloudfrontPlugin struct {
	AwsConn awsconnector.AWSConnector
}

func NewCloudfrontPlugin(aws_conn *awsconnector.AWSConnector) CloudfrontPlugin {
	cfp := CloudfrontPlugin{AwsConn: *aws_conn}

	return cfp
}

func (cfp CloudfrontPlugin) NormalizeCFDistroFQDN(fqdn *string) string {
	// CloudFront currently returns a `.` appended to the fqdn
	// we'll need to get rid of it so that it can be lookup up properly

	return strings.TrimSuffix(*fqdn, ".")
}

func (cfp CloudfrontPlugin) GetResources() *[]types.DistributionSummary {
	var distros []types.DistributionSummary

	cf_client := cloudfront.NewFromConfig(cfp.AwsConn.AwsConfig)
	paginator := cloudfront.NewListDistributionsPaginator(cf_client, &cloudfront.ListDistributionsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Error("error when running CloudFront search :: [ ", err, " ]")
		}

		distros = append(distros, output.DistributionList.Items...)
	}

	return &distros
}

func (cfp CloudfrontPlugin) SearchResources(tgt_ip *string) *types.DistributionSummary {
	var cfDistroFQDN string
	var cfIpAddrs *[]net.IP
	var matchedDistro types.DistributionSummary

	cf_resources := cfp.GetResources()

	for _, cfDistro := range *cf_resources {
		cfDistroFQDN = cfp.NormalizeCFDistroFQDN(cfDistro.DomainName)
		cfIpAddrs = utils.LookupFQDN(&cfDistroFQDN)

		for _, ipAddr := range *cfIpAddrs {
			if ipAddr.String() == *tgt_ip {
				matchedDistro = cfDistro
			}
		}
	}

	return &matchedDistro
}

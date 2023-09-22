package plugin

import (
	"context"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	generalResource "github.com/magneticstain/ip-2-cloudresource/src/resource"
	"github.com/magneticstain/ip-2-cloudresource/src/utils"
)

type CloudfrontPlugin struct {
	AwsConn awsconnector.AWSConnector
}

func NewCloudfrontPlugin(awsConn *awsconnector.AWSConnector) CloudfrontPlugin {
	cfp := CloudfrontPlugin{AwsConn: *awsConn}

	return cfp
}

func NormalizeCFDistroFQDN(fqdn *string) string {
	// CloudFront currently returns a `.` appended to the fqdn
	// we'll need to get rid of it so that it can be lookup up properly

	return strings.TrimSuffix(*fqdn, ".")
}

func (cfp CloudfrontPlugin) GetResources() (*[]types.DistributionSummary, error) {
	var distros []types.DistributionSummary

	cfClient := cloudfront.NewFromConfig(cfp.AwsConn.AwsConfig)
	paginator := cloudfront.NewListDistributionsPaginator(cfClient, &cloudfront.ListDistributionsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return &distros, err
		}

		distros = append(distros, output.DistributionList.Items...)
	}

	return &distros, nil
}

func (cfp CloudfrontPlugin) SearchResources(tgtIP *string) (*generalResource.Resource, error) {
	var cfDistroFQDN string
	var cfIPAddrs *[]net.IP
	var matchingResource generalResource.Resource

	cfResources, err := cfp.GetResources()
	if err != nil {
		return &matchingResource, err
	}

	for _, cfDistro := range *cfResources {
		cfDistroFQDN = NormalizeCFDistroFQDN(cfDistro.DomainName)
		cfIPAddrs, err = utils.LookupFQDN(&cfDistroFQDN)
		if err != nil {
			return &matchingResource, err
		}

		for _, ipAddr := range *cfIPAddrs {
			if ipAddr.String() == *tgtIP {
				matchingResource.RID = *cfDistro.ARN
				log.Debug("IP found as CloudFront distribution -> ", matchingResource.RID)
			}
		}
	}

	return &matchingResource, nil
}

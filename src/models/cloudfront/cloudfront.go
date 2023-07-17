package cloudfront

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip2cr/src/models/aws_connector"
)

type CloudfrontPlugin struct {
	AwsConn awsconnector.AWSConnector
}

func NewCloudfrontPlugin(aws_conn *awsconnector.AWSConnector) CloudfrontPlugin {
	cfp := CloudfrontPlugin{AwsConn: *aws_conn}

	return cfp
}

func (cfp CloudfrontPlugin) SearchResources() []types.DistributionSummary {
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

	return distros
}

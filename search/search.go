package search

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rollbar/rollbar-go"
	log "github.com/sirupsen/logrus"

	awscontroller "github.com/magneticstain/ip-2-cloudresource/aws"
	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws/aws_connector"
	iamp "github.com/magneticstain/ip-2-cloudresource/aws/plugin/iam"
	ipfuzzing "github.com/magneticstain/ip-2-cloudresource/aws/svc/ip_fuzzing"
	azurecontroller "github.com/magneticstain/ip-2-cloudresource/azure"
	gcpcontroller "github.com/magneticstain/ip-2-cloudresource/gcp"
	generalResource "github.com/magneticstain/ip-2-cloudresource/resource"
)

type Search struct {
	AWSCtrlr                   awscontroller.AWSController
	AzureCtrlr                 azurecontroller.AzureController
	CloudSvcs                  []string
	GCPCtrlr                   gcpcontroller.GCPController
	MatchedResource            generalResource.Resource
	IpAddr, Platform, TenantID string
}

func (search *Search) connectToPlatform() (bool, error) {
	// generate a connection to the specified platform via plugin
	// GCP does not require a connector as it uses ADC ( https://cloud.google.com/docs/authentication/application-default-credentials / https://archive.is/tSqC2 )

	switch search.Platform {
	case "aws":
		ac, err := awscontroller.New()
		if err != nil {
			return false, err
		}

		search.AWSCtrlr = ac
	case "azure":
		azc, err := azurecontroller.New()
		if err != nil {
			return false, err
		}

		search.AzureCtrlr = azc
	}

	return true, nil
}

func (search Search) ReconcileCloudSvcParam(cloudSvc string) []string {
	var cloudSvcs []string

	if cloudSvc == "all" {
		switch search.Platform {
		case "aws":
			cloudSvcs = awscontroller.GetSupportedSvcs()
		case "azure":
			cloudSvcs = azurecontroller.GetSupportedSvcs()
		case "gcp":
			cloudSvcs = gcpcontroller.GetSupportedSvcs()
		}
	} else if strings.Contains(cloudSvc, ",") {
		// csv provided, split the values into a slice
		cloudSvcs = strings.Split(cloudSvc, ",")
	} else {
		// assume single service
		cloudSvcs = []string{cloudSvc}
	}

	return cloudSvcs
}

func (search Search) RunIPFuzzing(doAdvIPFuzzing bool) ([]string, error) {
	var svcSet []string
	var fuzzedSvc string
	var err error

	fuzzedSvc, err = ipfuzzing.FuzzIP(search.IpAddr, doAdvIPFuzzing)
	if err != nil {
		return svcSet, err
	}

	// normalize service name to lowercase
	fuzzedSvc = strings.ToLower(fuzzedSvc)

	if fuzzedSvc == "" || fuzzedSvc == "unknown" {
		log.Info("could not determine service via IP fuzzing")
		return svcSet, err
	}

	log.Info("IP fuzzing determined the associated cloud service is: ", fuzzedSvc)
	svcSet = append(svcSet, fuzzedSvc)

	// all ELBs act within EC2 infrastructure, so we will need to add the elb services as well if that's the case
	if fuzzedSvc == "ec2" {
		svcSet = append(svcSet, "elbv1", "elbv2")
	}

	return svcSet, err
}

func (search Search) doAccountLevelSearch(acctID string, doNetMapping bool) (generalResource.Resource, error) {
	var acctAliases []string
	var matchingResource generalResource.Resource
	var err error

	if acctID != "current" && search.Platform == "aws" {
		// resolve account's aliases
		iamp := iamp.IAMPlugin{AwsConn: search.AWSCtrlr.PrincipalAWSConn}
		acctAliases, err = iamp.GetResources()
		if err != nil {
			return matchingResource, err
		}

		log.Info("starting resource search in AWS account: ", acctID, " ", acctAliases)
	} else {
		log.Info("starting resource search in current account")
	}

	for _, svc := range search.CloudSvcs {
		switch search.Platform {
		case "aws":
			matchingResource, err = search.AWSCtrlr.SearchAWSSvc(search.IpAddr, svc, doNetMapping)
		case "azure":
			matchingResource, err = search.AzureCtrlr.SearchAzureSvc(search.TenantID, search.IpAddr, svc, &matchingResource)
		case "gcp":
			matchingResource, err = search.GCPCtrlr.SearchGCPSvc(search.TenantID, search.IpAddr, svc, &matchingResource)
		default:
			errorMsg := fmt.Sprintf("%s is not a supported platform for searching", search.Platform)
			return matchingResource, errors.New(errorMsg)
		}

		if err != nil {
			return matchingResource, err
		} else if matchingResource.RID != "" {
			// resource was found
			matchingResource.AccountID = acctID
			matchingResource.AccountAliases = acctAliases

			break
		}
	}

	return matchingResource, nil
}

func (search Search) runSearchWorker(matchingResourceBuffer chan<- generalResource.Resource, acctID string, orgSearchRoleName string, doNetMapping bool, wg *sync.WaitGroup) {
	defer wg.Done()

	// org support is only available for AWS at this time
	if acctID != "current" && search.Platform == "aws" {
		// replace connector with assumed role connector before running rest of logic
		acctRoleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", acctID, orgSearchRoleName)
		ac, err := awsconnector.NewAWSConnectorAssumeRole(acctRoleArn, aws.Config{})
		if err != nil {
			log.Error("error when assuming role for account search worker: ", err)
			return
		}

		search.AWSCtrlr.PrincipalAWSConn = ac
	}

	resultResource, err := search.doAccountLevelSearch(acctID, doNetMapping)
	if err != nil {
		log.Error("error when running search within account search worker: ", err)
	} else if resultResource.RID != "" {
		matchingResourceBuffer <- resultResource
		return
	}
}

func (search *Search) initSearchWorkers(acctsToSearch []string, orgSearchRoleName string, doNetMapping bool) bool {
	log.Info("beginning resource gathering")

	matchingResourceBuffer := make(chan generalResource.Resource, 1)
	var wg sync.WaitGroup

	for _, acctID := range acctsToSearch {
		wg.Add(1)
		go rollbar.WrapAndWait(
			search.runSearchWorker,
			matchingResourceBuffer,
			acctID,
			orgSearchRoleName,
			doNetMapping,
			&wg,
		)
	}

	go func() {
		wg.Wait()
		close(matchingResourceBuffer)
	}()

	resultResource, found := <-matchingResourceBuffer
	if found {
		search.MatchedResource = resultResource
	}

	return found
}

func (search *Search) StartSearch(cloudSvc string, doIPFuzzing bool, doAdvIPFuzzing bool, doOrgSearch bool, orgSearchXaccountRoleARN string, orgSearchRoleName string, orgSearchOrgUnitID string, doNetMapping bool) (bool, error) {
	var resourceFound bool
	var err error

	_, err = search.connectToPlatform()
	if err != nil {
		log.Fatal("error when connecting to ", search.Platform, ": ", err)
	}

	// TODO: move this to init function
	search.CloudSvcs = search.ReconcileCloudSvcParam(cloudSvc)

	if doIPFuzzing || doAdvIPFuzzing {
		search.CloudSvcs, err = search.RunIPFuzzing(doAdvIPFuzzing)
		if err != nil {
			return resourceFound, err
		}
	}

	var acctsToSearch []string
	if doOrgSearch {
		log.Info("starting org account enumeration")

		acctsToSearch, err = search.AWSCtrlr.FetchOrgAcctIds(orgSearchOrgUnitID, orgSearchXaccountRoleARN)
		if err != nil {
			return resourceFound, err
		}
	} else {
		acctsToSearch = append(acctsToSearch, "current")
	}

	resourceFound = search.initSearchWorkers(acctsToSearch, orgSearchRoleName, doNetMapping)

	return resourceFound, nil
}

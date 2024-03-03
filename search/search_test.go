package search_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"

	awscontroller "github.com/magneticstain/ip-2-cloudresource/aws"
	awsconnector "github.com/magneticstain/ip-2-cloudresource/aws/aws_connector"
	gcpcontroller "github.com/magneticstain/ip-2-cloudresource/gcp"
	"github.com/magneticstain/ip-2-cloudresource/search"
	"golang.org/x/exp/slices"
)

type TestIPAddr struct {
	ipAddr string
}

func searchFactory(ipAddr string) search.Search {
	ac, _ := awscontroller.New()
	gcpc := gcpcontroller.GCPController{}

	search := search.Search{
		AWSCtrlr: ac,
		GCPCtrlr: gcpc,
		IpAddr:   ipAddr,
	}

	return search
}

func ipFactory() []TestIPAddr {
	var ipData []TestIPAddr

	ipData = append(
		ipData,
		TestIPAddr{"52.4.175.237"},  // CloudFront
		TestIPAddr{"65.8.191.186"},  // ALB
		TestIPAddr{"35.170.192.9"},  // EC2
		TestIPAddr{"3.218.196.10"},  // NLB
		TestIPAddr{"34.205.13.193"}, // Classic ELB
		TestIPAddr{"2600:1f18:243e:1300:4685:5a7:7c28:c53a"}, // EC2 IPv6
	)

	return ipData
}

func ipFuzzingCloudSvcsFactory() []string {
	cloudSvcs := []string{
		"cloudfront",
		"ec2",
		"elbv1",
		"elbv2",
		"unknown",
	}

	return cloudSvcs
}

func TestReconcileCloudSvcParam_AllSvcs(t *testing.T) {
	var tests = []struct {
		platform, cloudSvc  string
		expectedCloudSvcSet []string
	}{
		{"aws", "all", []string{
			"cloudfront",
			"ec2",
			"elbv1",
			"elbv2",
		}},
		{"gcp", "all", []string{
			"compute",
			"load_balancing",
		}},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.platform, td.cloudSvc)

		search := searchFactory("")

		t.Run(testName, func(t *testing.T) {
			search.Platform = td.platform
			res := search.ReconcileCloudSvcParam(td.cloudSvc)

			for _, svc := range res {
				if !slices.Contains(td.expectedCloudSvcSet, svc) {
					t.Errorf("Cloud service reconciliation failed due to unexpected service; returned: %s", svc)
				}
			}
		})
	}
}

func TestReconcileCloudSvcParam_InvalidSvcs(t *testing.T) {
	var tests = []struct {
		platform, cloudSvc  string
		expectedCloudSvcSet []string
	}{
		{"gcp", "all", []string{
			"cloudfront",
			"ec2",
			"elbv1",
			"elbv2",
		}},
		{"aws", "all", []string{
			"compute",
		}},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.platform, td.cloudSvc)

		search := searchFactory("")

		t.Run(testName, func(t *testing.T) {
			search.Platform = td.platform
			res := search.ReconcileCloudSvcParam(td.cloudSvc)

			for _, svc := range res {
				if slices.Contains(td.expectedCloudSvcSet, svc) {
					t.Errorf("Cloud service reconciliation failed due to unexpected service; returned: %s", svc)
				}
			}
		})
	}
}

func TestReconcileCloudSvcParam_CsvSvcs(t *testing.T) {
	var tests = []struct {
		svcCsv              string
		expectedCloudSvcSet []string
	}{
		{"ec2,elbv1,elbv2", []string{
			"ec2",
			"elbv1",
			"elbv2",
		}},
	}

	for _, td := range tests {
		testName := td.svcCsv

		search := searchFactory("")

		t.Run(testName, func(t *testing.T) {

			res := search.ReconcileCloudSvcParam(td.svcCsv)

			for i, svc := range res {
				if svc != td.expectedCloudSvcSet[i] {
					t.Errorf("Invalid service returned from service reconciliation via CSV slug; expected: %s, returned: %s", td.expectedCloudSvcSet[i], svc)
				}
			}
		})
	}
}

func TestReconcileCloudSvcParam_SingleSvc(t *testing.T) {
	var tests = []struct {
		cloudSvc string
	}{
		{"ec2"},
	}

	for _, td := range tests {
		testName := td.cloudSvc

		search := searchFactory("")

		t.Run(testName, func(t *testing.T) {

			res := search.ReconcileCloudSvcParam(td.cloudSvc)

			if res[0] != td.cloudSvc {
				t.Errorf("Invalid service returned from service reconciliation of singular service; expected: %s, returned: %s", td.cloudSvc, res[0])
			}
		})
	}
}

func TestRunIPFuzzing(t *testing.T) {
	var tests = ipFactory()

	validSvcs := ipFuzzingCloudSvcsFactory()
	for _, td := range tests {
		testName := td.ipAddr

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			fuzzedSvcSet, err := search.RunIPFuzzing(false)
			if err != nil {
				t.Errorf("Basic IP fuzzing routine unexpectedly failed; error: %s", err)
			}

			for _, fuzzedSvc := range fuzzedSvcSet {
				if !slices.Contains[[]string, string](validSvcs, fuzzedSvc) {
					t.Errorf("Basic IP fuzzing routine failed; unexpected service was returned: %s", fuzzedSvc)
				}
			}
		})
	}
}

func TestRunIPFuzzing_AdvancedFuzzing(t *testing.T) {
	var tests = ipFactory()

	validSvcs := ipFuzzingCloudSvcsFactory()
	for _, td := range tests {
		testName := td.ipAddr

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			fuzzedSvcSet, err := search.RunIPFuzzing(true)
			if err != nil {
				t.Errorf("Basic IP fuzzing routine unexpectedly failed; error: %s", err)
			}

			for _, fuzzedSvc := range fuzzedSvcSet {
				if !slices.Contains[[]string, string](validSvcs, fuzzedSvc) {
					t.Errorf("Basic IP fuzzing routine failed; unexpected service was returned: %s", fuzzedSvc)
				}
			}
		})
	}
}

func TestStartSearch_CloudSvcs(t *testing.T) {
	var tests = []struct {
		ipAddr, cloudSvc string
	}{
		{"1.1.1.1", "cloudfront"},
		{"299.11.906.43", "elbv1,elbv2"},
		{"2600:9000:24eb:dc00:1:3b80:4f00:21", "ec2"},
		{"x2600:9000:24eb:XYZ1:1:3b80:4f00:21", "not_a_svc"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch(td.cloudSvc, false, false, false, "", "", "", false)

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "bool"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with IP fuzzing disabled has failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_NoFuzzing(t *testing.T) {
	var tests = []struct {
		ipAddr string
	}{
		{"1.1.1.1"},
		{"299.11.906.43"},
		{"2600:9000:24eb:dc00:1:3b80:4f00:21"},
		{"x2600:9000:24eb:XYZ1:1:3b80:4f00:21"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch("all", false, false, false, "", "", "", false)

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "bool"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with IP fuzzing disabled has failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_BasicFuzzing(t *testing.T) {
	var tests = ipFactory()

	for _, td := range tests {
		testName := td.ipAddr

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch("all", true, false, false, "", "", "", false)

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "bool"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with IP fuzzing enabled failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_AdvancedFuzzing(t *testing.T) {
	var tests = ipFactory()

	for _, td := range tests {
		testName := td.ipAddr

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch("all", true, false, false, "", "", "", false)

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "bool"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with advanced IP fuzzing disabled failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_OrgSearchEnabled(t *testing.T) {
	var tests = ipFactory()

	for _, td := range tests {
		testName := td.ipAddr

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch("all", false, false, true, "", "ip2cr-org-role", "", false)

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "bool"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with AWS Organizations support enabled has failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_OrgSearchEnabled_XaccountSvcRole(t *testing.T) {
	var tests = []struct {
		orgXaccountRoleARN string
	}{
		{"arn:aws:iam::123456789012:role/valid_role"},
		{"arn:aws:bad::123456:role/invalid_role"},
		{"arn:aws:iam::123456789012:user/invalid_user"}, // only roles are supported by this function
	}

	for _, td := range tests {
		testName := td.orgXaccountRoleARN

		t.Run(testName, func(t *testing.T) {
			ac, _ := awsconnector.NewAWSConnectorAssumeRole(td.orgXaccountRoleARN, aws.Config{})

			acType := reflect.TypeOf(ac.AwsConfig)

			if acType.Name() != "Config" {
				t.Errorf("AWS connector failed to connect when assuming xaccount service role; wanted aws.Config type, received %s", acType.Name())
			}
		})
	}
}

func TestStartSearch_OrgSearchEnabled_TargetOUID_ParentOrgID(t *testing.T) {
	// REF: https://docs.aws.amazon.com/organizations/latest/APIReference/API_Organization.html#organizations-Type-Organization-Id

	var tests = []struct {
		orgID, ipAddr string
	}{
		{"o-0000000000", "1.1.1.1"},
		{"o-9999999999", "1.1.1.1"},
		{"o-1234567890abcde", "1.1.1.1"},
		{"o-00000000000000000000000000000000", "1.1.1.1"},
		{"o-99999999999999999999999999999999", "1.1.1.1"},
	}

	for _, td := range tests {
		testName := td.orgID

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch("all", false, false, true, "", "ip2cr-org-role", td.orgID, false)

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "bool"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with AWS Organizations support enabled has failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_OrgSearchEnabled_TargetOUID_ChildOUID(t *testing.T) {
	// REF: https://docs.aws.amazon.com/organizations/latest/APIReference/API_OrganizationalUnit.html#organizations-Type-OrganizationalUnit-Id

	var tests = []struct {
		OUID, ipAddr string
	}{
		{"ou-aaaa", "1.1.1.1"},
		{"ou-zzzz", "1.1.1.1"},
		{"ou-1234567890abcde", "1.1.1.1"},
		{"ou-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "1.1.1.1"},
		{"ou-zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", "1.1.1.1"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch("all", false, false, true, "", "ip2cr-org-role", td.OUID, false)

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "bool"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with AWS Organizations support enabled has failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

func TestStartSearch_OrgSearchEnabled_TargetOUID_InvalidID(t *testing.T) {
	var tests = []struct {
		OUID, ipAddr string
	}{
		{"abcde", "1.1.1.1"},
		{"ou-", "1.1.1.1"},
		{"ou-a", "1.1.1.1"},
		{"ou-zzz", "1.1.1.1"},
		{"ou-zzzABCd", "1.1.1.1"},
		{"ou-1234567890abcde", "1.1.1.1"},
		{"ou-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "1.1.1.1"},
		{"ou-zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", "1.1.1.1"},
	}

	for _, td := range tests {
		testName := td.ipAddr

		search := searchFactory(td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			res, _ := search.StartSearch("all", false, false, true, "", "ip2cr-org-role", td.OUID, false)

			matchedResourceType := reflect.TypeOf(res)
			expectedType := "bool"
			if matchedResourceType.Name() != expectedType {
				t.Errorf("Overall search with AWS Organizations support enabled has failed; expected %s after search, received %s", expectedType, matchedResourceType.Name())
			}
		})
	}
}

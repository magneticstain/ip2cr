package ipfuzzing_test

import (
	"fmt"
	"testing"

	"golang.org/x/exp/slices" // Update to the stable `slices` package once 1.12 becomes oldstable ( Issue #112 )

	ipfuzzing "github.com/magneticstain/ip-2-cloudresource/src/svc/ip_fuzzing"
)

func GetValidCloudSvcs(includeUnknownSvc bool) *[]string {
	validSvcs := []string{
		"CLOUDFRONT",
		"EC2",
		"ELBv1",
		"ELBv2",
	}

	if includeUnknownSvc {
		validSvcs = append(validSvcs, "UNKNOWN")
	}

	return &validSvcs
}

func TestMapFQDNToSvc(t *testing.T) {
	var tests = []struct {
		cloudSvc, fqdn string
	}{
		{"CLOUDFRONT", "server-65-8-191-186.bos50.r.cloudfront.net."},
		{"EC2", "ec2-35-170-192-9.compute-1.amazonaws.com."},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.fqdn)

		t.Run(testName, func(t *testing.T) {
			mappedSvc, err := ipfuzzing.MapFQDNToSvc(&td.fqdn)
			if err != nil {
				t.Errorf("unexpected error received when attempting to match %s service: %s", td.cloudSvc, err)
			}

			if *mappedSvc == "" || *mappedSvc != td.cloudSvc {
				t.Errorf("failed to map FQDN to service; EXPECTED SVC: %s , MAPPED SVC: %s , FQDN: %s", td.cloudSvc, *mappedSvc, td.fqdn)
			}
		})
	}
}

func TestMapFQDNToSvc_InvalidSvcs(t *testing.T) {
	var tests = []struct {
		cloudSvc, fqdn string
	}{
		{"cf", "server-65-8-191-186.bos50.r.cloudfront.net."},
		{"EC2----", "ec2-35-170-192-9.compute-1.amazonaws.com."},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.fqdn)

		t.Run(testName, func(t *testing.T) {
			mappedSvc, err := ipfuzzing.MapFQDNToSvc(&td.fqdn)
			if err != nil {
				t.Errorf("unexpected error received when attempting to match invalid service %s: %s", td.cloudSvc, err)
			}

			if *mappedSvc == td.cloudSvc {
				t.Errorf("expected error when mapping FQDN to invalid service, but was successful; EXPECTED SVC: %s , MAPPED SVC: %s , FQDN: %s", td.cloudSvc, *mappedSvc, td.fqdn)
			}
		})
	}
}

func TestRunAdvancedFuzzing(t *testing.T) {
	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"CLOUDFRONT", "65.8.191.186"},
		{"EC2", "52.4.175.237"},  // ALB
		{"EC2", "35.170.192.9"},  // EC2 - IPv4
		{"EC2", "3.218.196.10"},  // NLB
		{"EC2", "34.205.13.193"}, // Classic ELB
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			fuzzedSvc, err := ipfuzzing.RunAdvancedFuzzing(&td.ipAddr)
			if err != nil {
				t.Errorf("unexpected error received when attempting to fuzz %s service using advanced fuzzing: %s", td.cloudSvc, err)
			}

			if *fuzzedSvc != td.cloudSvc {
				t.Errorf("failed to fuzz service using advanced IP fuzzing; EXPECTED SVC: %s , FUZZED SVC: %s , IP: %s", td.cloudSvc, *fuzzedSvc, td.ipAddr)
			}
		})
	}
}

func TestRunAdvancedFuzzing_InvalidIPs(t *testing.T) {
	var tests = []struct {
		cloudSvc, ipAddr string
	}{
		{"CLOUDFRONT", "555.5.5.555"},
		{"EC2", "999.888.77777.1"},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.cloudSvc, td.ipAddr)

		t.Run(testName, func(t *testing.T) {
			_, err := ipfuzzing.RunAdvancedFuzzing(&td.ipAddr)
			if err == nil {
				t.Errorf("expected error when performing advanced IP fuzzing, but didn't")
			}
		})
	}
}

func TestFuzzIP(t *testing.T) {
	var tests = []struct {
		ipAddr        string
		useAdvFuzzing bool
	}{
		{"1.1.1.1", false},
		{"35.170.192.9", false},
		{"2600:1f18:243e:1300:4685:5a7:7c28:c53a", false},
		{"1.1.1.1", true},
		{"35.170.192.9", true},
		{"2600:1f18:243e:1300:4685:5a7:7c28:c53a", true},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%t", td.ipAddr, td.useAdvFuzzing)
		validSvcs := GetValidCloudSvcs(true)

		t.Run(testName, func(t *testing.T) {
			svcName, err := ipfuzzing.FuzzIP(&td.ipAddr, td.useAdvFuzzing)
			if err != nil {
				t.Errorf("unexpected error received when attempting to fuzz %s IP using general fuzzing: %s", td.ipAddr, err)
			}

			if !slices.Contains[[]string, string](*validSvcs, *svcName) {
				t.Errorf("unexpected service name when performing IP fuzzing tests; received %s", *svcName)
			}
		})
	}
}

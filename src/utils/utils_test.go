package utils_test

import (
	"fmt"
	"testing"

	"github.com/magneticstain/ip2cr/src/utils"
)

func TestReverseDNSLookup(t *testing.T) {
	var tests = []struct {
		ipAddr, fqdn    string
		expectedVerdict bool
	}{
		{"1.1.1.1", "one.one.one.one.", true},
		{"8.8.8.8", "dns.google.", true},
		{"1.1.1.1", "google.com.", false},
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.fqdn, td.ipAddr)
		t.Run(testName, func(t *testing.T) {
			fqdns, _ := utils.ReverseDNSLookup(&td.ipAddr)

			fqdnFound := false
			for _, fqdn := range fqdns {
				if fqdn == td.fqdn {
					fqdnFound = true
					break
				}
			}

			if fqdnFound != td.expectedVerdict {
				t.Errorf("reverse IP lookup failed; expected %s to be %t FQDN for %s", td.fqdn, td.expectedVerdict, td.ipAddr)
			}
		})
	}
}

func TestLookupFQDN(t *testing.T) {
	var tests = []struct {
		fqdn, ipAddr    string
		expectedVerdict bool
	}{
		{"example.com", "93.184.216.34", true},                       // accurate IPv4 lookup
		{"example.com", "1.1.1.1", false},                            // inaccurate IPv4 lookup
		{"example.com", "2606:2800:220:1:248:1893:25c8:1946", true},  // accurate IPv6 lookup
		{"example.com", "2600:9000:24eb:3a00:1:3b80:4f00:21", false}, // inaccurate IPv6 lookup
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.fqdn, td.ipAddr)
		t.Run(testName, func(t *testing.T) {
			ipAddrs, _ := utils.LookupFQDN(&td.fqdn)

			ipFound := false
			for _, ipAddr := range *ipAddrs {
				if ipAddr.String() == td.ipAddr {
					ipFound = true
					break
				}
			}

			if ipFound != td.expectedVerdict {
				t.Errorf("FQDN lookup failed; expected %s to be %t IP address for %s", td.ipAddr, td.expectedVerdict, td.fqdn)
			}
		})
	}
}

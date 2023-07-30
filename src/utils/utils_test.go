package utils_test

import (
	"fmt"
	"testing"

	"github.com/magneticstain/ip2cr/src/utils"
)

func TestLookupFQDN(t *testing.T) {
	var tests = []struct {
		fqdn, ipAddr    string
		expectedVerdict bool
	}{
		{"example.com", "93.184.216.34", true}, // valid lookup
		{"example.com", "1.1.1.1", false},      // invalid lookup
	}

	for _, td := range tests {
		testName := fmt.Sprintf("%s_%s", td.fqdn, td.ipAddr)
		t.Run(testName, func(t *testing.T) {
			ipAddrs := utils.LookupFQDN(&td.fqdn)

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

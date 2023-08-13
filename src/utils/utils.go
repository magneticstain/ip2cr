package utils

import (
	"net"
)

func ReverseDNSLookup(ipAddr *string) ([]string, error) {
	return net.LookupAddr(*ipAddr)
}

func LookupFQDN(fqdn *string) (*[]net.IP, error) {
	var ipAddrs []net.IP

	ipAddrs, err := net.LookupIP(*fqdn)

	return &ipAddrs, err
}

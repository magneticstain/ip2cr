package utils

import (
	"net"
)

func ReverseDNSLookup(ipAddr *string) ([]string, error) {
	// NOTE: IPv6 addresses are not supported (see https://datatracker.ietf.org/doc/html/rfc8501)
	return net.LookupAddr(*ipAddr)
}

func LookupFQDN(fqdn *string) (*[]net.IP, error) {
	var ipAddrs []net.IP

	ipAddrs, err := net.LookupIP(*fqdn)

	return &ipAddrs, err
}

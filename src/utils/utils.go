package utils

import (
	"net"
)

func LookupFQDN(fqdn *string) (*[]net.IP, error) {
	var ipAddrs []net.IP

	ipAddrs, err := net.LookupIP(*fqdn)

	return &ipAddrs, err
}

package utils

import (
	"net"

	log "github.com/sirupsen/logrus"
)

func LookupFQDN(fqdn *string) *[]net.IP {
	var ipAddrs []net.IP

	ipAddrs, err := net.LookupIP(*fqdn)
	if err != nil {
		log.Error("failed to lookup IP of CloudFront distribution :: [ FQDN: ", *fqdn, " // ERR: ", err, " ]")
	}

	return &ipAddrs
}

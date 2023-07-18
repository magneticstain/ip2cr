package utils

import (
	"net"

	log "github.com/sirupsen/logrus"
)

func lookupFQDN(fqdn string) string {
	ipAddr, err := net.LookupCNAME(fqdn)
	if err != nil {
		log.Error("failed to lookup IP of CloudFront distribution :: [ FQDN: ", fqdn, " // ERR: ", err, " ]")
		return ""
	}

	return ipAddr
}

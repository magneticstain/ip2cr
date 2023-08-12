package ipfuzzing

import (
	"net"

	log "github.com/sirupsen/logrus"

	awsipprefix "github.com/magneticstain/ip2cr/src/svc/ip_fuzzing/models/aws_ip_prefix"
)

func FuzzIP(ipAddr *string) (*string, error) {
	var cloudSvc *string

	awsIpSet, err := FetchIpRanges()
	if err != nil {
		return cloudSvc, err
	}
	log.Debug("AWS public IP dataset => ", awsIpSet)

	// AWS divides their prefixes by IP version, so we should determine that first to reduce the number of checks needed
	// Here, we're checking the IP version and then converting the prefixes for the given version to generic prefixes
	var ipPrefixSet *[]awsipprefix.GenericAWSPrefix
	parsedIpAddr := net.ParseIP(*ipAddr)
	if parsedIpAddr.To4() != nil {
		// IPv4
		ipPrefixSet, err = ConvertIpPrefixesToGeneric(&awsIpSet.Prefixes, nil)
	} else {
		// IPv6
		ipPrefixSet, err = ConvertIpPrefixesToGeneric(nil, &awsIpSet.IPv6Prefixes)
	}

	if err != nil {
		log.Error("not able to convert versioned IP prefix groups to generic")
	} else {
		log.Debug("IP prefix set reduced by version successfully")
	}

	fuzzedSvc, err := ResolveIpAddrToCloudSvc(ipAddr, ipPrefixSet)
	if err != nil {
		return cloudSvc, err
	}

	if *fuzzedSvc == "AMAZON" {
		// AWS's generic service name for ranges
		normalizedSvcName := "UNKNOWN"
		cloudSvc = &normalizedSvcName
	} else {
		cloudSvc = fuzzedSvc
	}

	return cloudSvc, nil
}

package ipfuzzing

import (
	"net"
	"regexp"

	log "github.com/sirupsen/logrus"

	awsfqdnregexmap "github.com/magneticstain/ip-2-cloudresource/src/svc/ip_fuzzing/models/aws_fqdn_regex_map"
	awsipprefix "github.com/magneticstain/ip-2-cloudresource/src/svc/ip_fuzzing/models/aws_ip_prefix"
	"github.com/magneticstain/ip-2-cloudresource/src/utils"
)

func MapFQDNToSvc(fqdn *string) (*string, error) {
	var re *regexp.Regexp
	var svcName *string

	regexMap := awsfqdnregexmap.GetRegexMap()
	for svc, regex := range regexMap {
		// check if fqdn matches the associated regex; if so, we have our service
		re = regexp.MustCompile(regex)

		if re.MatchString(*fqdn) {
			svcName = &svc
			break
		}
	}

	return svcName, nil
}

func RunAdvancedFuzzing(ipAddr *string) (*string, error) {
	// perform a reverse DNS lookup on the IP and then use heuristics to try to determine the associated service
	var cloudSvc *string

	reverseLookupResult, err := utils.ReverseDNSLookup(ipAddr)
	if err != nil {
		return cloudSvc, err
	} else {
		log.Debug("reverse DNS lookup for IP [ ", *ipAddr, " ] resolves to [ ", reverseLookupResult, " ]")
	}

	var svcName *string
	for _, fqdn := range reverseLookupResult {
		svcName, err = MapFQDNToSvc(&fqdn)
		if err != nil {
			return cloudSvc, nil
		}

		if svcName != nil {
			// service was found!
			cloudSvc = svcName

			log.Debug("advanced fuzzing identified the service as [ ", *svcName, " ]")

			// we assume that the first match is the true match; we can adjust this if real-world results don't match this presumption
			break
		}
	}

	return cloudSvc, nil
}

func FuzzIP(ipAddr *string, attemptAdvancedFuzzing bool) (*string, error) {
	var cloudSvc *string

	awsIpSet, err := FetchIpRanges()
	if err != nil {
		return cloudSvc, err
	}
	log.Debug("AWS public IP dataset loaded")

	// AWS divides their prefixes by IP version, so we should determine that first to reduce the number of checks needed
	// Here, we're checking the IP version and then converting the prefixes for the given version to generic prefixes
	var ipPrefixSet *[]awsipprefix.GenericAWSPrefix
	parsedIPAddr := net.ParseIP(*ipAddr)
	parsedIPAddrV4 := parsedIPAddr.To4()
	if parsedIPAddrV4 != nil {
		// IPv4
		ipPrefixSet, err = ConvertIpPrefixesToGeneric(&awsIpSet.Prefixes, nil)
	} else {
		// IPv6
		ipPrefixSet, err = ConvertIpPrefixesToGeneric(nil, &awsIpSet.IPv6Prefixes)
	}
	if err != nil {
		log.Error("not able to convert versioned IP prefix groups to generic; [ ERR: ", err, " ]")
	} else {
		log.Debug("IP prefix set reduced by version successfully")
	}

	fuzzedSvc, err := ResolveIPAddrToCloudSvc(ipAddr, ipPrefixSet)
	if err != nil {
		return cloudSvc, err
	}
	// if AWS IP range scanning doesn't work, we can try advanced fuxxing, which uses reverse DNS and heuristics to try to determine the service
	// NOTE: this only works for IPv4 at this time as AWS doesn't appear to have PTR records setup for their IPv6 prefixes
	if parsedIPAddrV4 == nil {
		log.Debug("skipping advanced fuzzing since IPv6 is not supported by this feature")
	} else if *fuzzedSvc == "" || *fuzzedSvc == "AMAZON" {
		log.Debug("basic IP fuzzing failed to determine cloud service")

		if attemptAdvancedFuzzing {
			log.Debug("starting advanced IP fuzzing")
			advFuzzResult, err := RunAdvancedFuzzing(ipAddr)
			if err != nil {
				return cloudSvc, err
			}

			if advFuzzResult != nil {
				return advFuzzResult, nil
			}
		}
	} else {
		// cloud service was found
		log.Debug("basic IP fuzzing determined the IP belongs to the ", *fuzzedSvc, " service")
		return fuzzedSvc, nil
	}

	if *fuzzedSvc == "AMAZON" || *fuzzedSvc == "" {
		// AWS's generic service name for ranges
		normalizedSvcName := "UNKNOWN"
		cloudSvc = &normalizedSvcName
	} else {
		cloudSvc = fuzzedSvc
	}

	return cloudSvc, nil
}

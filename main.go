package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	searchctlr "github.com/magneticstain/ip2cr/src/controllers/search"
	awsconnector "github.com/magneticstain/ip2cr/src/models/aws_connector"
)

func main() {
	ipAddr := flag.String("ipaddr", "127.0.0.1", "IP address to search for")
	cloudSvc := flag.String("svc", "all", "Specific cloud service to search")
	flag.Parse()

	log.Info("starting IP-2-CloudResource...")

	log.Info("generating AWS connection...")
	ac := awsconnector.New()

	log.Info("searching for IP ", *ipAddr, " in ", *cloudSvc, " service(s)")
	matchedResource := searchctlr.StartSearch(&ac, ipAddr)

	if matchedResource.RID != "" {
		log.Info("resource found -> [ ", matchedResource.RID, " ]")
	} else {
		log.Info("resource not found :( better luck next time!")
	}
}

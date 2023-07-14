package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	searchctlr "github.com/magneticstain/ip2cr/src/controllers/search"
	awsconnector "github.com/magneticstain/ip2cr/src/models/aws_connector"
)

func initLogging() {
	// initialize all Logrus configs and prefs
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	ipAddr := flag.String("ipaddr", "127.0.0.1", "IP address to search for")
	cloudSvc := flag.String("svc", "all", "Specific cloud service to search")
	flag.Parse()

	initLogging()

	log.Info("Starting IP-2-CloudResource...")

	log.Info("Checking AWS connection...")
	ac := awsconnector.AWSConnector{}

	log.Info("Searching for IP ", *ipAddr, " in ", *cloudSvc, " service(s)")
	searchctlr.StartSearch(&ac)
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip2cr/src/aws_connector"
	"github.com/magneticstain/ip2cr/src/search"
)

func main() {
	silent := flag.Bool("silent", false, "If enabled, only output the results")
	ipAddr := flag.String("ipaddr", "127.0.0.1", "IP address to search for")
	cloudSvc := flag.String("svc", "all", "Specific cloud service to search")
	jsonOutput := flag.Bool("json", false, "Outputs results in JSON format; implies usage of --silent flag")
	flag.Parse()

	if *jsonOutput {
		*silent = true
	}

	if *silent {
		log.SetOutput(io.Discard)
	}

	log.Info("starting IP-2-CloudResource")

	log.Debug("generating AWS connection")
	ac, err := awsconnector.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("searching for IP ", *ipAddr, " in ", *cloudSvc, " service(s)")
	searchCtlr := search.NewSearch(&ac)
	matchedResource, err := searchCtlr.StartSearch(ipAddr)
	if err != nil {
		log.Fatal("failed to run search :: [ ERR: ", err, " ]")
	}

	if matchedResource.RID != "" {
		if !*silent {
			log.Info("resource found -> [ ", matchedResource.RID, " ]")
		} else {
			if *jsonOutput {
				output, err := json.Marshal(matchedResource)
				if err != nil {
					errMap := map[string]error{"error": err}
					errMapJSON, _ := json.Marshal(errMap)

					fmt.Printf("%s\n", errMapJSON)
				} else {
					fmt.Printf("%s\n", output)
				}
			} else {
				// plaintext
				fmt.Println(matchedResource.RID)
			}
		}
	} else {
		log.Info("resource not found :( better luck next time!")
	}
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/rollbar/rollbar-go"
	log "github.com/sirupsen/logrus"

	awsconnector "github.com/magneticstain/ip-2-cloudresource/src/aws_connector"
	"github.com/magneticstain/ip-2-cloudresource/src/resource"
	"github.com/magneticstain/ip-2-cloudresource/src/search"
)

func initRollbar() {
	rollbar.SetToken("98a9cbd56b164657ab447d79eac9b258")
	rollbar.SetCaptureIp(rollbar.CaptureIpAnonymize)
	rollbar.SetServerRoot("github.com/magneticstain/ip-2-cloudresource")
	rollbar.SetCodeVersion("v1.0.0")
	rollbar.SetEnvironment("development")
}

func OutputResults(matchedResource *resource.Resource, silent *bool, jsonOutput *bool) {
	acctAliasFmted := strings.Join(matchedResource.AccountAliases, ", ")

	if !*silent {
		if matchedResource.RID != "" {
			var acctStr string
			if matchedResource.AccountId == "current" {
				acctStr = "current account"
			} else {
				acctStr = fmt.Sprintf("account [ %s ( %s ) ]", matchedResource.AccountId, acctAliasFmted)
			}

			log.Info("resource found -> [ ", matchedResource.RID, " ] in ", acctStr)
		} else {
			log.Info("resource not found :( better luck next time!")
		}
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
			if matchedResource.RID != "" {
				fmt.Println(matchedResource.RID)
				fmt.Printf("%s (%s)", matchedResource.AccountId, acctAliasFmted)
			} else {
				fmt.Println("not found")
			}
		}
	}
}

func RunCloudSearch(ipAddr *string, cloudSvc *string, ipFuzzing *bool, advIpFuzzing *bool, orgSearch *bool, orgSearchRoleName *string, silent *bool, jsonOutput *bool) {
	// cloud connections
	log.Debug("generating AWS connection")
	ac, err := awsconnector.New()
	if err != nil {
		log.Fatal(err)
	}

	// search
	log.Info("searching for IP ", *ipAddr, " in ", *cloudSvc, " service(s)")
	searchCtlr := search.NewSearch(&ac)
	matchedResource, err := searchCtlr.StartSearch(ipAddr, *ipFuzzing, *advIpFuzzing, *orgSearch, *orgSearchRoleName)
	if err != nil {
		log.Fatal("failed to run search :: [ ERR: ", err, " ]")
	}

	// output
	OutputResults(&matchedResource, silent, jsonOutput)
}

func main() {
	// CLI param parsing
	silent := flag.Bool("silent", false, "If enabled, only output the results")
	ipAddr := flag.String("ipaddr", "127.0.0.1", "IP address to search for")
	cloudSvc := flag.String("svc", "all", "Specific cloud service to search")
	ipFuzzing := flag.Bool("ip-fuzzing", true, "Toggle the IP fuzzing feature to evaluate the IP and help optimize search (not recommended for small accounts)")
	advIpFuzzing := flag.Bool("adv-ip-fuzzing", true, "Toggle the advanced IP fuzzing feature to perform a more intensive heuristics evaluation to fuzz the service (not recommended for IPv6 addresses)")
	orgSearch := flag.Bool("org-search", false, "Search through all child accounts of the organization for resources, as well as target account (target account should be parent account)")
	orgSearchRoleName := flag.String("org-search-role-name", "ip2cr", "The name of the role in each child account of an AWS Organization to assume when performing a search")
	jsonOutput := flag.Bool("json", false, "Outputs results in JSON format; implies usage of --silent flag")
	verboseOutput := flag.Bool("verbose", false, "Outputs all logs, from debug level to critical")
	flag.Parse()

	if *jsonOutput {
		*silent = true
	}
	if *silent {
		log.SetOutput(io.Discard)
	}
	if *verboseOutput {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("starting IP-2-CloudResource")

	initRollbar()

	rollbar.WrapAndWait(RunCloudSearch, ipAddr, cloudSvc, ipFuzzing, advIpFuzzing, orgSearch, orgSearchRoleName, silent, jsonOutput)

	rollbar.Close()
}

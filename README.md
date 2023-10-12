# ip-2-cloudresource

[![Build and Test - Main Branch](https://github.com/magneticstain/ip-2-cloudresource/actions/workflows/build_main.yml/badge.svg)](https://github.com/magneticstain/ip-2-cloudresource/actions/workflows/build_main.yml)
[![Build and Test - Release Branches](https://github.com/magneticstain/ip-2-cloudresource/actions/workflows/build_release.yml/badge.svg)](https://github.com/magneticstain/ip-2-cloudresource/actions/workflows/build_release.yml)
[![Release](https://github.com/magneticstain/ip-2-cloudresource/actions/workflows/release.yml/badge.svg)](https://github.com/magneticstain/ip-2-cloudresource/actions/workflows/release.yml)

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/5137ec7cf2d14a9c9fc3eac1cd37e0d3)](https://app.codacy.com/gh/magneticstain/ip-2-cloudresource/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

![GitHub commits since latest release (by SemVer including pre-releases)](https://img.shields.io/github/commits-since/magneticstain/ip-2-cloudresource/v1.0.1)
![GitHub issues](https://img.shields.io/github/issues/magneticstain/ip-2-cloudresource)
![GitHub pull requests](https://img.shields.io/github/issues-pr/magneticstain/ip-2-cloudresource)
![GitHub all releases](https://img.shields.io/github/downloads/magneticstain/ip-2-cloudresource/total)

## Summary

IP-2-CloudResource (IP2CR) is a tool used for correlating a cloud IP address with its associated resources. It focuses on providing as much context to the user as possible, as fast as possible.

### Disclaimer

I created this project mainly to learn Go. It should be fine for a cloud admin running this on their workstation, trying to identify a resource in the AWS account. But I wouldn't necessarily integrate it with my production monitoring system/lunar lander.

## Features

- Built for speed and efficiency while only generating a small resource footprint
- Supports finding IPs for:
  - CloudFront
  - ALBs & NLBs (and probably GLBs, but hasn't been tested yet)
  - Classic ELBs
  - EC2 instances with public IP addresses
- Support for searching through accounts within an AWS Organization
- IPv6 support
- JSON output to easily integrate with scripts

### Roadmap

#### 2023

- [X] EC2 support ( [Issue #11](https://github.com/magneticstain/ip-2-cloudresource/issues/11) )
- [X] Classic ELB support ( [Issue #29](https://github.com/magneticstain/ip-2-cloudresource/issues/29) )
- [X] JSON output ( [Issue #37](https://github.com/magneticstain/ip-2-cloudresource/issues/37) )
- [X] IP service fuzzing (perform a reverse DNS lookup to identify the services to search, leading to faster results)  ( [Issue #39](https://github.com/magneticstain/ip-2-cloudresource/issues/39) )
- [X] Support for installing using Homebrew ( [Issue #77](https://github.com/magneticstain/ip-2-cloudresource/issues/77) )
- [X] AWS Organizations support ( [Issue #38](https://github.com/magneticstain/ip-2-cloudresource/issues/38) )
- [ ] Network path calculation ( [Issue #44](https://github.com/magneticstain/ip-2-cloudresource/issues/44) )
- [ ] Add Support For Concurrent Account-Based Resource Searches When Running With AWS Org Support ( [Issue #141](https://github.com/magneticstain/ip-2-cloudresource/issues/141) )
- [ ] Improve Unit Tests to Integrate With tf-ip2cr and Use Real Data ( [Issue #142](https://github.com/magneticstain/ip-2-cloudresource/issues/142) )

## Prerequisites

### OS

- Linux
- MacOS

Windows should probably work, but I'm not able to test it at this time.

### Go

IP2CR supports running on n-1 minor versions of Golang, aka [stable and old-stable](https://go.dev/dl/#stable).

## Install

### Homebrew

The easiest way to install IP2CR is to use [Homebrew](https://brew.sh). With homebrew installed, run the following to install IP2CR:

```bash
brew tap magneticstain/ip2cr
brew install ip2cr
```

Once the install completes, you should be able to run the `ip2cr` command from your path.

### Go Install

If homebrew isn't an option, you can also install using `go install`:

```bash
go install github.com/magneticstain/ip-2-cloudresource@latest
```

### Binary

Finally, if all else fails or you prefer directly managing IP2CR, to install the IP2CR binary, download the version that matches your OS from the latest release in the release page. The link for that page can be found to the right of this page.

After downloading, extract to your preferred app directory. E.g.:

```bash
mkdir -p /opt/ip2cr
tar -xvzf ip2cr_Linux_x86_64_v0.0.3.tar.gz -C /opt/ip2cr/
cd /opt/ip2cr/
```

## Usage

After installing, `cd` to the app directory if not already there and run the `ip2cr` binary.

```bash
> ./ip2cr --help
Usage of ./ip2cr:
  -adv-ip-fuzzing
    	Toggle the advanced IP fuzzing feature to perform a more intensive heuristics evaluation to fuzz the service (not recommended for IPv6 addresses) (default true)
  -ip-fuzzing
    	Toggle the IP fuzzing feature to evaluate the IP and help optimize search (not recommended for small accounts) (default true)
  -ipaddr string
    	IP address to search for (default "127.0.0.1")
  -json
    	Outputs results in JSON format; implies usage of --silent flag
  -org-search
    	Search through all child accounts of the organization for resources, as well as target account (target account should be parent account)
  -org-search-role-name string
    	The name of the role in each child account of an AWS Organization to assume when performing a search (default "ip2cr")
  -silent
    	If enabled, only output the results
  -svc string
    	Specific cloud service to search (default "all")
  -verbose
    	Outputs all logs, from debug level to critical
```

## Testing/Demo

You can use the Terraform plans provided here to generate sample resources in AWS for testing.

<https://github.com/magneticstain/tf-ip2cr>

## Support, Feature Requests, and General Community Discussion

The `Discusions` module of this repository has been setup as a place to get support, request new features, and facilitate any general discorse related to IP2CR.

If you are having an issue when using IP2CR, or just need general help, you should start here as opposed to creating an Issue. Any Issues created for support purposes will be closed.

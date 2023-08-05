# ip2cr

[![Build and Test](https://github.com/magneticstain/ip2cr/actions/workflows/build.yml/badge.svg)](https://github.com/magneticstain/ip2cr/actions/workflows/build.yml)
[![Release](https://github.com/magneticstain/ip2cr/actions/workflows/release.yml/badge.svg)](https://github.com/magneticstain/ip2cr/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/magneticstain/ip2cr/branch/main/graph/badge.svg?token=YI5A0BA12D)](https://codecov.io/gh/magneticstain/ip2cr)

## Summary

IP-2-CloudResource (IP2CR) is a tool used for correlating a cloud IP address with its associated resources. It focuses on providing as much context to the user as possible, as fast as possible.

### Disclaimer

I created this project mainly to learn Go. It should be fine for a cloud admin running this on their workstation, trying to identify a resource in the AWS account. But I wouldn't necessarily integrate it with my production monitoring system.

## Features

- Built for speed and efficiency while only generating a small resource footprint
- Supports finding IPs for:
  - CloudFront
  - ALBs & NLBs (and probably GLBs, but hasn't been tested yet)
  - Classic ELBs
  - EC2 instances with public IP addresses
- IPv6 support

### Roadmap

- [X] EC2 support ( [Issue #11](https://github.com/magneticstain/ip2cr/issues/11) )
- [X] Classic ELB support ( [Issue #29](https://github.com/magneticstain/ip2cr/issues/29) )
- [ ] JSON output ( [Issue #37](https://github.com/magneticstain/ip2cr/issues/37) )
- [ ] AWS Organizations support ( [Issue #38](https://github.com/magneticstain/ip2cr/issues/38) )
- [ ] IP service fuzzing (perform a reverse DNS lookup to identify the services to search, leading to faster results)  ( [Issue #39](https://github.com/magneticstain/ip2cr/issues/39) )

## Prerequisites

### OS

- Linux
- MacOS

Windows should probably work, but I'm not able to test it at this time.

### Go

IP2CR supports running on n-1 minor versions of Golang, aka [stable and old-stable](https://go.dev/dl/#stable).

## Testing/Demo

You can use the Terraform plans provided here to generate sample resources in AWS for testing.

<https://github.com/magneticstain/tf-ip2cr>

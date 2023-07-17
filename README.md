# ip2cr

## Summary

IP-2-CloudResource: a tool used for correlating a cloud IP address with its associated resources.

IP2CR focuses on providing as much context to the user as possible, as fast as possible.

### Disclaimer

I created this project mainly to learn Go. It should be fine for a cloud admin running this on their workstation, trying to identify a resource in the AWS account. But I wouldn't necessarily integrate it with my production monitoring system.

## Prerequisites

## Go

Tested on Go v1.20.5

## Development

### Testing

You can use the Terraform plans provided here to generate sample resources in AWS for testing.

<https://github.com/magneticstain/tf-ip2cr>

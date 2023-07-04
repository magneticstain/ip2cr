package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("Starting IP-2-CloudResource...")

	ipAddr := flag.String("ipaddr", "127.0.0.1", "IP address to search for")
	flag.Parse()

	fmt.Println("Using IP", *ipAddr)
}

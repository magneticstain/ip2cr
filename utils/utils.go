package utils

import (
	"net"
	"strings"

	"github.com/rollbar/rollbar-go"
)

func InitRollbar(appVer string) {
	rollbar.SetToken("98a9cbd56b164657ab447d79eac9b258")
	rollbar.SetCaptureIp(rollbar.CaptureIpAnonymize)
	rollbar.SetServerHost("anonymous")
	rollbar.SetServerRoot("github.com/magneticstain/ip-2-cloudresource")
	rollbar.SetCodeVersion(appVer)
	rollbar.SetEnvironment("development")
}

func ReverseDNSLookup(ipAddr string) ([]string, error) {
	// NOTE: IPv6 addresses are not supported (see https://datatracker.ietf.org/doc/html/rfc8501)
	return net.LookupAddr(ipAddr)
}

func LookupFQDN(fqdn string) ([]net.IP, error) {
	var ipAddrs []net.IP

	ipAddrs, err := net.LookupIP(fqdn)

	return ipAddrs, err
}

func FormatStrSliceAsCSV(strs []string) string {
	formattedStr := "[" + strings.Join(strs, ",") + "]"

	return formattedStr
}

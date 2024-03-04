package utils

import (
	"errors"
	"net"
	"strings"

	"github.com/rollbar/rollbar-go"
)

func InitRollbar(appEnv, appVer string) {
	rollbar.SetToken("98a9cbd56b164657ab447d79eac9b258")
	rollbar.SetCaptureIp(rollbar.CaptureIpAnonymize)
	rollbar.SetServerHost("anonymous")
	rollbar.SetServerRoot("github.com/magneticstain/ip-2-cloudresource")
	rollbar.SetEnvironment(appEnv)
	rollbar.SetCodeVersion(appVer)
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

func DetermineIpAddrVersion(ipAddr string) (int, error) {
	var ipVer int

	parsedIPAddr := net.ParseIP(ipAddr)
	if parsedIPAddr == nil {
		return ipVer, errors.New("invalid IP provided")
	}

	parsedIPAddrV4 := parsedIPAddr.To4()
	if parsedIPAddrV4 != nil {
		ipVer = 4
	} else {
		ipVer = 6
	}

	return ipVer, nil
}

func FormatStrSliceAsCSV(strs []string) string {
	formattedStr := "[" + strings.Join(strs, ",") + "]"

	return formattedStr
}

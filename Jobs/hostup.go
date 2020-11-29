package jobs

import (
	//"fmt"
	"github.com/bachittle/ping-go/pinger"
	"github.com/bachittle/ping-go/utils"
	"net"
	"strings"
)

// HostUp checks if a given host is online with ICMP echo packets.
// uses an external library (made for the purpose of this assignment) called ping-go
// check it out here: https://pkg.go.dev/github.com/bachittle/ping-go
// pass a regular domain name or CIDR name
// EX:
// - jobs.HostUp("google.com")		// hostname DNS lookup
// - jobs.HostUp("192.168.2.1/24")  // local CIDR lookup
//
// returns a string of hostnames that are online and offline.
func HostUp(hostname string) (online []string, offline []string, err error) {
	p := pinger.NewPinger()
	var IPs []net.IP
	if strings.Contains(hostname, "/24") {
		// do CIDR translation to multiple IPs
		IPs, err = utils.GetIPv4CIDR(hostname)
	} else {
		var ip net.IP
		ip, err = utils.GetIPv4(hostname)
		if err != nil {
			return
		}
		IPs = append(IPs, ip)
	}
	for i := 0; i < len(IPs); i++ {
		_, err = p.SetDst(IPs[0])
		if err != nil {
			return
		}
		p.SetAmt(1)
		err = p.Ping()
		if err != nil {
			return
		}
	}
	return
}

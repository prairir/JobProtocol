package jobs

import (
	"fmt"
	"github.com/bachittle/ping-go/pinger"
	"github.com/bachittle/ping-go/utils"
	"io"
	"net"
	"os"
	"strings"
	//"sync"
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
func HostUp(hostname string, w io.Writer) (online []string, offline []string, err error) {
	if w == nil {
		w = os.Stdout
	}
	var IPs []net.IP
	if strings.Contains(hostname, "/") {
		// do CIDR translation to multiple IPs
		IPs, err = utils.GetIPv4CIDR(hostname)
		if err != nil {
			return
		}
	} else {
		var ip net.IP
		ip, err = utils.GetIPv4(hostname)
		if err != nil {
			return
		}
		IPs = append(IPs, ip)
	}

	for _, ip := range IPs {
		fmt.Fprintln(w, "pinging ip", ip)
		p := pinger.NewPinger()
		_, err = p.SetDst(ip)
		if err != nil {
			return
		}
		p.SetAmt(1)
		p.Ping()
		// code is bad but i dont know how to fix without a massive refactor
		_, err = p.Pong(20)
		if err != nil {
			offline = append(offline, ip.String())
		} else {
			online = append(online, ip.String())
		}
	}
	return
}

package jobs

import (
	"fmt"
	"github.com/bachittle/ping-go-v2"
	"io"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
	var ip net.IP
	fmt.Println(hostname)

	con := ping2.Controller{
		SrcIP: net.IP{192, 168, 50, 77},
	}

	CIDRnum := 0
	if n := strings.Index(hostname, "/"); n != -1 {
		// do CIDR translation to multiple IPs
		CIDRnum, err = strconv.Atoi(hostname[n+1:])
		if err != nil {
			return
		}
		ip = net.ParseIP(hostname[:n]).To4()
		con.DstIPs = append(con.DstIPs, ping2.CustomIP{IP: ip, Subnet: &CIDRnum})
	} else {
		//var ip net.IP
		ip = net.ParseIP(hostname).To4()
		con.DstIPs = append(con.DstIPs, ping2.CustomIP{IP: ip, Subnet: nil})
	}
	fmt.Println("ip:", ip, "cidr:", CIDRnum)
	con.Init()
	fmt.Println("starting...")
	dict := con.SendAndRecv(5 * time.Second)
	fmt.Println("stopping...")

	offlineMap := make(map[string]bool)

	// i know this is bad but using GenerateIPs I was getting issues...
	// would skip IPs unless if I checked it with fmt print, kind of like schroedingers cat
	var tmpIP string
	for ip := range ping2.GenerateIPs(con.DstIPs) {
		if ip.String() != tmpIP {
			tmpIP = ip.String()

		} else {
			if ip[3] == 255 {
				ip[2]++
			} else {
				ip[3]++
			}
			tmpIP = ip.String()
		}
		if dict[ip.String()] == false {
			offlineMap[ip.String()] = true
		}
	}
	fmt.Println("n:", len(offlineMap))
	fmt.Println("m:", len(dict))

	for ip := range dict {
		online = append(online, ip)
	}
	for ip := range offlineMap {
		offline = append(offline, ip)
	}
	sort.Strings(online)
	sort.Strings(offline)
	return
}

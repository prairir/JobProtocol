package jobs

import (
	"errors"
	"fmt"
	"github.com/bachittle/ping-go/pinger"
	"github.com/bachittle/ping-go/utils"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"io"
	"net"
	"os"
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

	// returns error, or nil if success
	for i := 0; i < len(IPs); i++ {
		chanMsg := make(chan *icmp.Message, 1)
		chanErr := make(chan error, 1)
		p := pinger.NewPinger()
		_, err = p.SetDst(IPs[i])
		if err != nil {
			return
		}
		// timeout goroutine
		go func() {
			msg, err := p.PingOne(nil)
			if err != nil {
				chanErr <- err
			} else {
				chanMsg <- msg
			}
		}()
		// timeout select checking
		var msg *icmp.Message
		select {
		case res := <-chanMsg:
			msg = res
		case res := <-chanErr:
			err = res
		case <-time.After(30 * time.Millisecond):
			err = errors.New("timeout")
		}
		if err != nil || msg.Type != ipv4.ICMPTypeEchoReply {
			offline = append(offline, IPs[i].String())
		} else {
			fmt.Fprintln(w, IPs[i], msg, err)
			online = append(online, IPs[i].String())
		}
	}
	return
}

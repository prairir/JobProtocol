package jobs

import (
	"testing"
)

func TestTraceroute(t *testing.T) {
	IPs := []string{
		"192.168.50.1",
		"google.com",
	}

	for _, ip := range IPs {
		t.Run(ip, func(t *testing.T) {
			ips := Traceroute("Wi-Fi", ip)
			t.Log(ips)
		})
	}
}

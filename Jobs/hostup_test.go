package jobs

import (
	"fmt"
	"strings"
	"testing"
)

func TestHostUp(t *testing.T) {
	// normal network test
	tests := []string{
		"192.168.50.1", // set this to your default gateway
		"google.com",
	}
	// CIDR network test (last byte, 255 hosts)
	testsLen := len(tests) // keep old length to for-loop over
	for i := 0; i < testsLen; i++ {
		tests = append(tests, fmt.Sprint(tests[i], "/24"))
	}
	t.Log("testing array:", tests)
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			var numHosts int
			if strings.Contains(test, "/24") {
				numHosts = 255
			} else {
				numHosts = 1
			}
			online, offline, err := HostUp(test)
			if err != nil {
				t.Error(err)
			}
			t.Log("online:", online, "offline:", offline)
			if len(online)+len(offline) != numHosts {
				t.Error("Invalid number of online/offline hosts:", len(online), "+", len(offline), "!=", numHosts)
			}
		})
	}
}

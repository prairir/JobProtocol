package jobs

import (
	"fmt"
	"math"
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
	CIDRnum := 24
	for i := 0; i < testsLen; i++ {
		tests = append(tests, fmt.Sprint(tests[i], "/", CIDRnum))
	}
	t.Log("testing array:", tests)
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			var numHosts int
			if strings.Contains(test, fmt.Sprint("/", CIDRnum)) {
				numHosts = int(math.Pow(2, float64(32-CIDRnum)))
			} else {
				numHosts = 1
			}
			online, offline, err := HostUp(test, nil)
			if err != nil {
				if err.Error() == "timeout" {
					t.Log("timeout error, no worries. ")
				} else {
					t.Error(fmt.Sprint("ERROR: ", err))
				}
			}
			t.Log("online hosts:", online)
			t.Log("online hosts length:", len(online))
			if len(online)+len(offline) != numHosts {
				t.Error("Invalid number of online/offline hosts:", len(online), "+", len(offline), "!=", numHosts)
			}
		})
	}
}

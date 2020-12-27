package jobs

import (
	"testing"
	"time"
)

func TestNeighbours(t *testing.T) {
	test := make(map[string]interface{})
	test["192.168.50.77"] = []string{"c4:9d:ed:28:de:af"}
	test["192.168.50.1"] = []string{"70:4d:7b:5a:57:78"}
	sameLAN, report := Neighbours(test, 10*time.Second)
	t.Log("sameLAN:", sameLAN)
	t.Log("report:", report)
}

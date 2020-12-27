package jobs

import (
	"testing"
)

func TestGetMACstr(t *testing.T) {
	mac, err := GetMACstr()
	t.Log(mac, err)
}

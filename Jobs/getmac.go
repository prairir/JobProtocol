// GetMAC gets the current seekers MAC address. This is useful for the neighbours code

package jobs

import (
	"net"
)

// GetMACstr gets the hardware address of this computer
// returns the interface associated with the mac (the string), and the mac associated
func GetMACstr() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	macsSet := make(map[string]bool)
	for _, iface := range ifaces {
		macStr := iface.HardwareAddr.String()
		if macStr != "" {
			macsSet[macStr] = true
		}
	}
	var macsArr []string
	for key := range macsSet {
		macsArr = append(macsArr, key)
	}
	return macsArr, nil
}

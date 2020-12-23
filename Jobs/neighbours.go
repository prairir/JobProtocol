package jobs

import (
	"bytes"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
	"time"
)

//Neighbours checks for job seeker neighbours within the same LAN
// it checks if ip addresses in an ip addr array are floating around the LAN
func Neighbours(duration time.Duration) (sameLAN [][]byte, report []map[string][]byte) {

	go func() {
		// opens packet souce on an interface
		if handle, err := pcap.OpenLive("\\Device\\NPF_{2F557FE1-6FE0-4B4F-8C12-3B40FC5C87A6}", 1600, true, duration); err != nil {
			panic(err)
		} else {
			// deserialize / decode -> turn bytes from eth0 into packet
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			//iterate through packets
			for packet := range packetSource.Packets() {
				m := make(map[string][]byte)
				// decode ethernet and IPv4 layers
				// checks for linklayer
				if eth := packet.Layer(layers.LayerTypeEthernet); eth != nil {
					//extracts srctMAC and dstMAC
					eth, ok := eth.(*layers.Ethernet)
					if !ok {
						panic("invalid ethernet packet")
					}
					src, dst := eth.LinkFlow().Endpoints()
					// adds src to []byte
					m["mac_src"] = src.Raw()
					m["mac_dst"] = dst.Raw()
				}
				// checks for IPv4 layer
				if ip4 := packet.Layer(layers.LayerTypeIPv4); ip4 != nil {
					// extracts end points, srcIP and dstIP
					ip4, ok := ip4.(*layers.IPv4)
					if ok {
						src, dst := ip4.NetworkFlow().Endpoints()
						// adds src to []byte
						m["ip4_src"] = src.Raw()
						m["ip4_dst"] = dst.Raw()
					}
				}
				// compare to array of net.IP
				ifaces, err := net.Interfaces()
				if err != nil {
					panic(err)
				}
				for _, i := range ifaces {
					addrs, err := i.Addrs()
					if err != nil {
						panic(err)
					}
					// handle err
					for _, addr := range addrs {
						var ip net.IP
						switch v := addr.(type) {
						case *net.IPNet:
							ip = v.IP

							if bytes.Equal(ip, m["ip4_src"]) {
								sameLAN = append(sameLAN, ip)
							} else if bytes.Equal(ip, m["ip4_dst"]) {
								sameLAN = append(sameLAN, ip)
							}
							break
						}
					}

					report = append(report, m)
				}
			}
		}
	}()

	time.Sleep(duration)
	return
}

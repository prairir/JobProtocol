package jobs

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"time"
)

func Neighbours(duration time.Duration, addr_net []byte) []map[string][]byte {

	var report []map[string][]byte
	var sameLAN []map[string][]byte

	go func() {
		// opens packet souce on an interface
		if handle, err := pcap.OpenLive("eth0", 1600, true, duration); err != nil {
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
				
				for addr := range addr_net {
				    if addr = m["ip4_src"] {
				        sameLAN = append(sameLAN, m["ip4_src"])
				    } 
				    else if  addr = m["ip4_dst"] {
				        sameLAN = append(sameLAN, m["ip4_dst"])
				    }
				}

				report = append(report, m)
			}
		}
	}()

	time.Sleep(duration)
	return report
	return sameLAN
}

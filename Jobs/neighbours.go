package jobs

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	//"strconv"
	"strings"
	"sync"
	"time"
)

//Neighbours checks for job seeker neighbours within the same LAN
// it checks if ip addresses in an ip addr array are floating around the LAN
// returns job seekers within the same LAN, as well as a report of all the packets for future reference.
//
// macMap and sameLAN structure:
// {
//		"192.168.50.1" : [
// 			"12:34:56:78:9A:BC",
// 			"34:56:78:9A:BC:DE"
// 		]
// }
func Neighbours(macMap map[string]interface{}, duration time.Duration) (sameLAN map[string]map[string]int, report map[string]map[string]int) {

	report = make(map[string]map[string]int)
	sameLAN = make(map[string]map[string]int)
	var doneChans []chan struct{}
	var mut sync.Mutex
	// opens packet souce on an interface
	ifaces, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}
	for i, iface := range ifaces {
		handle, err := pcap.OpenLive(iface.Name, 1600, true, duration)
		if err != nil {
			panic(err)
		}
		doneChans = append(doneChans, make(chan struct{}))
		// deserialize / decode -> turn bytes from eth0 into packet
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		//iterate through packets
		go func(index int) {
			for packet := range packetSource.Packets() {
				select {
				case <-doneChans[index]:
					return
				default:
				}
				// decode ethernet and IPv4 layers
				// checks for linklayer
				var macSrc net.HardwareAddr
				var macDst net.HardwareAddr
				if eth := packet.Layer(layers.LayerTypeEthernet); eth != nil {
					//extracts srctMAC and dstMAC
					eth, ok := eth.(*layers.Ethernet)
					if !ok {
						panic("invalid ethernet packet")
					}
					src, dst := eth.LinkFlow().Endpoints()
					macSrc = src.Raw()
					macDst = dst.Raw()
					log.Println("macSrc:", macSrc)
					log.Println("macDst:", macDst)
					// adds src to []byte
				}
				// checks for IPv4 layer
				if ip4 := packet.Layer(layers.LayerTypeIPv4); ip4 != nil {
					// extracts end points, srcIP and dstIP
					ip4, ok := ip4.(*layers.IPv4)
					if ok {
						ipSrc, ipDst := ip4.NetworkFlow().Endpoints()
						// adds src to []byte
						if macSrc != nil && macDst != nil {
							if ipSrc.Raw() != nil && ipDst.Raw() != nil {
								mut.Lock()
								if report[ipSrc.String()] == nil {
									report[ipSrc.String()] = make(map[string]int)
								}
								if report[ipDst.String()] == nil {
									report[ipDst.String()] = make(map[string]int)
								}
								report[ipSrc.String()][macSrc.String()]++
								report[ipDst.String()][macDst.String()]++
								mut.Unlock()
								for ip, macs := range macMap {
									// values in json that may not be IPs...
									if ip == "duration" {
										continue
									}
									//var port int
									n := strings.Index(ip, ":")
									if n != -1 {
										ip = ip[:n]
										/*
											port, err = strconv.Atoi(ip[n+1:])
											if err != nil {
												panic(err)
											}
										*/
									}
									macs, ok := macs.([]string)
									if !ok {
										continue
									}
									for _, mac := range macs {
										if report[ip][mac] > 0 {
											if sameLAN[ip] == nil {
												sameLAN[ip] = make(map[string]int)
											}
											sameLAN[ip][mac]++
										}
									}
								}
							}
						}
					}
				}
				// compare to array of net.IP
				/*
					ifaces, err := net.Interfaces()
					if err != nil {
						panic(err)
					}
					for _, iface := range ifaces {
						addrs, err := iface.Addrs()
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
									mut.Lock()
									sameLAN = append(sameLAN, ip)
									mut.Unlock()
								} else if bytes.Equal(ip, m["ip4_dst"]) {
									mut.Lock()
									sameLAN = append(sameLAN, ip)
									mut.Unlock()
								}
								break
							}
						}
					}
				*/
			}
		}(i)
	}

	time.Sleep(duration)
	// signify closure of all goroutines
	for _, ch := range doneChans {
		close(ch) // closing channel will "read" from it on loop
	}
	return
}

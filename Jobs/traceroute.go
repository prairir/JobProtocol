package jobs

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/prairir/JobProtocol/Globals"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

// Traceroute traces the route of a given IP address string, then returns a list of IPs that are the route that
// was taken to that IP
func Traceroute(ifaceName string, ipStr string) []net.IP {
	ci := GetIfaces(ifaceName)
	fmt.Println(ci)
	ci.Init()
	ci.StartReading()
	tries := 0
Loop:
	for i := 0; i < 255; i++ {
		buf, err := ci.ICMPReqPacket(ipStr)
		if err != nil {
			panic(err)
		}
		fmt.Println("seq:", ci.Seq)
		ci.WriteData(buf)

		select {
		case t := <-ci.GotType:
			fmt.Println("t:", t)
			fmt.Println("results:", ci.ResultIPs)
			ci.Seq++
			if t == 0 {
				break Loop
			}
		case <-time.After(200 * time.Millisecond):
			fmt.Println("timeout")
			tries++
			if tries == 2 {
				tries = 0
				ci.Seq++
			}
		}
	}
	return ci.ResultIPs
}

// CustomIface is a custom interface type used in this library
// since pcap and net use different interface types
type CustomIface struct {
	PcapName     string
	NetName      string
	HardwareAddr net.HardwareAddr
	IPAddr       net.IP
	IPNet        net.IPNet
	ResultIPs    []net.IP

	dstIP   net.IP
	mut     sync.Mutex
	handle  *pcap.Handle
	stopper chan int
	GotType chan uint8
	gotMAC  chan string
	id      uint16 // an id that all icmp packets will use
	Seq     uint16 // an id that all icmp packets will use
}

// GetGateway gets the gateway IP and MAC.
func (c *CustomIface) GetGateway() (net.IP, net.HardwareAddr, error) {
	IP := c.IPAddr.To4()
	IP[3] = 1
	// Set up all the layers' fields we can.
	eth := layers.Ethernet{
		SrcMAC:       c.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	fmt.Println("IP:", IP)
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(c.HardwareAddr),
		SourceProtAddress: []byte(c.IPAddr),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(IP),
	}
	fmt.Println(arp)
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	gopacket.SerializeLayers(buf, opts, &eth, &arp)
	if err := c.handle.WritePacketData(buf.Bytes()); err != nil {
		return nil, nil, err
	}
	fmt.Println("sent arp")
	s := <-c.gotMAC
	m, err := net.ParseMAC(s)
	return IP, m, err
}

// Init sets default initializers
func (c *CustomIface) Init() {
	c.id = uint16(rand.Intn(65535))
	c.Seq = 1

	c.stopper = make(chan int)
	c.GotType = make(chan uint8, 1)
}

// GetIfaces gets ifaces
func GetIfaces(name string) *CustomIface {
	var ci CustomIface
	netIface, err := net.InterfaceByName(name)
	if err != nil {
		panic(err)
	}
	pcapIfaces, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}
	fmt.Println("iface:", netIface)
	netAddrs, _ := netIface.Addrs()
	netNet := netAddrs[1].String()
	n := strings.Index(netNet, "/")
	netAddr := netNet[:n]
	fmt.Println("address:", netAddr)
	fmt.Println("pcap Ifaces")
	for i, pcapIface := range pcapIfaces {
		pcapAddrs := pcapIface.Addresses
		if len(pcapAddrs) > 1 {
			pcapAddr := pcapAddrs[1].IP
			pcapAddrStr := pcapAddr.String()
			if pcapAddrStr == netAddr {
				fmt.Println(i, ":", pcapIface)
				ci.PcapName = pcapIface.Name
				ci.NetName = netIface.Name
				ci.HardwareAddr = netIface.HardwareAddr
				ip, network, err := net.ParseCIDR(netNet)
				if err != nil {
					panic(err)
				}
				ci.IPAddr = ip
				ci.IPNet = *network
			}
		}
	}
	return &ci
}

// ICMPReqPacket writes an ICMP packet to an ipv4 string (no CIDR)
func (c *CustomIface) ICMPReqPacket(ipStr string) ([]byte, error) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	// Set up all the layers' fields we can.
	//_, dstMac, err := c.GetGateway()
	macStr, err := globals.MACString()
	if err != nil {
		return nil, err
	}
	dstMac, err := net.ParseMAC(macStr)
	if err != nil {
		return nil, err
	}
	eth := layers.Ethernet{
		SrcMAC: c.HardwareAddr,
		//DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		DstMAC:       dstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}

	dstIPs, err := net.LookupIP(ipStr)
	if err != nil {
		return nil, err
	}
	fmt.Println(dstIPs)
	c.dstIP = dstIPs[0]
	//dstIP := net.ParseIP(ipStr).To4()
	ip4 := layers.IPv4{
		SrcIP:    c.IPAddr,
		DstIP:    c.dstIP,
		Version:  4,
		TTL:      uint8(c.Seq % 255),
		Protocol: layers.IPProtocolICMPv4,
	}

	icmp4 := layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       c.id,
		Seq:      c.Seq,
	}
	err = gopacket.SerializeLayers(buf, opts,
		&eth,
		&ip4,
		&icmp4,
	)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteData writes buffer data on the wire
func (c *CustomIface) WriteData(buf []byte) error {

	c.handle.WritePacketData(buf)
	return nil
}

// StartReading starts reading ICMP packets and keeps the handle in state if you want to use it. Close it by running Close()
func (c *CustomIface) StartReading() error {
	handle, err := pcap.OpenLive(c.PcapName, 65535, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	c.handle = handle
	go func() {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		//fmt.Println("starting packet receiver...")
		for packet := range packetSource.Packets() {
			if c.stopper != nil {
				// in the loop, if stopChan is ever given any values, break the loop.
				select {
				case <-c.stopper:
					break
				default:
				}
			}
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if icmpLayer != nil && ipLayer != nil {
				icmpLayer, ok := icmpLayer.(*layers.ICMPv4)
				if !ok {
					panic(ok)
				}
				ipLayer, ok := ipLayer.(*layers.IPv4)
				if !ok {
					panic(ok)
				}
				fmt.Println(icmpLayer)
				//fmt.Println("type:", icmpLayer.TypeCode.Type())
				if (icmpLayer.TypeCode.Type() != 0 && icmpLayer.TypeCode.Type() != 11) || icmpLayer.TypeCode.Code() != 0 {
					fmt.Println("not a valid response")
					fmt.Println(icmpLayer.TypeCode.Type(), icmpLayer.TypeCode.Code())
					continue
					/*
						} else if icmpLayer.Id != c.id {
							fmt.Println("not a valid response")
							fmt.Println(icmpLayer.Id, "!=", c.id)
							continue
						}
					*/
				} else if ipLayer.DstIP.String() != c.IPAddr.String() {
					fmt.Println("not a valid response")
					fmt.Println(ipLayer.DstIP, "!=", c.IPAddr)
					continue
				}
				fmt.Println("valid match!")
				c.mut.Lock()
				c.ResultIPs = append(c.ResultIPs, ipLayer.SrcIP)
				c.mut.Unlock()
				c.GotType <- icmpLayer.TypeCode.Type()
				if icmpLayer.TypeCode.Type() == 0 {
					fmt.Println("got zero!")
					return
				}
			} else if arpLayer != nil {
				arpLayer, ok := icmpLayer.(*layers.ARP)
				if !ok {
					//panic(err)
					fmt.Println("error")
				}
				fmt.Println("ARP:", arpLayer)
				c.gotMAC <- "test"
			}
		}
	}()
	return nil
}

// Close closes the handler and the asynchronous goroutine
func (c *CustomIface) Close() {
	c.stopper <- 1
}

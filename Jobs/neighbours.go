package main

import (

    "fmt"
    "time"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
)
   

func []byte neighbours(duration time.duration){

  var add_MAC [6]byte
  var add_IP  [4]byte
  var report  [10]byte

// opens packet souce on an interface
if handle, err := pcap.OpenLive("eth0", 1600, true, duration); err != nil {
  panic(err)
} else {

// deserialize / decode -> turn bytes from eth0 into packet
  packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
  //iterate through packets
  for packet := range packetSource.Packets() {
        // decode ethernet and IPv4 layers
        // checks for linklayer
        if eth := packet.Layers(layers.LayerTypeEthernet); eth != nil {
            //extracts srctMAC and dstMAC
            src, dst := eth.NetworkFlow().EndPoints()
            // adds src to []byte
            add_MAC = src
        }
          // checks for IPv4 layer
        if ip4 := packet.Layers(layers.LayerTypeIPv4); ip4 != nil {
          // extracts end points, srcIP and dstIP
          src, dst :=  ip4.NetworkFlow().EndPoints()
          // adds src to []byte
          add_IP := src
        }
        
     report := append(add_MAC, add_IP...)
     
     return report
  }
}

}


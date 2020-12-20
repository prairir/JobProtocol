package main

import (

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
)
   

func main(){

  /*var add_MAC [6]byte
  var add_IP  [4]byte
  var eth layers.Ethernet
  var ip4 layers.IPv4
  var ip6 layers.IPv6
  var tcp layers.TCP*/
  //parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp)
  
  decoded := []gopacket.LayerType{} //provides array of layers to terate through

// opens packet souce on an interface
if handle, err := pcap.OpenLive("eth0", 1600, true, pcap.BlockForever); err != nil {
  panic(err)
} else {

// deserialize / decode -> turn bytes from eth0 into packet
  packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
  //iterate through packets
  for packet := range packetSource.Packets() {
    
  
      // do stuff to packets
     
  }
}

}

 /* // searches layers
     for _, layerType := range decoded {
      switch layerType {
        case layers.LayerTypeEthernet:
          fmt.Println("    Eth ", eth.SrcMAC, eth.DstMAC)
        case layers.LayerTypeIPv4:
          fmt.Println("    IP4 ", ip4.SrcIP, ip4.DstIP)   
      }
      
      //serialize a layer
      buf := gopacket.NewSerializeBuffer()
      opts := gopacket.SerializeOptions{}  // See SerializeOptions for more details.
      err := ip.SerializeTo(buf, opts)
      if err != nil { panic(err) }
      fmt.Println(buf.Bytes())  // prints out a byte slice containing the serialized IPv4 layer.*/
      

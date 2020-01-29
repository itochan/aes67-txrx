package aes67

import (
	"fmt"
	"log"
	"net"

	"github.com/pion/rtp"
)

type Receiver struct {
	localIP  net.IP
	remoteIP net.IP
}

func NewReceiver(localIP net.IP, remoteIP net.IP) *Receiver {
	return &Receiver{localIP: localIP, remoteIP: remoteIP}
}

func (receiver Receiver) Receive() {
	udpAddr := &net.UDPAddr{IP: receiver.remoteIP, Port: aes67Port}
	var err error
	connect, err = net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	receivePacket()

	defer connect.Close()
}

func receivePacket() {
	buffer := make([]byte, 156)
	packet := &rtp.Packet{}

	var cnt int
	for {
		_, err := connect.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		cnt++
		packet.Unmarshal(buffer)
		fmt.Printf("Remote receiver got %d packets\n", cnt)
		fmt.Print(packet.String())
	}
}

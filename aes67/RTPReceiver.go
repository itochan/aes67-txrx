package aes67

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/pion/rtp"
)

type Receiver struct {
	senderIP         net.IP
	MulticastAddress net.IPNet
}

func NewReceiver(senderIP net.IP, multicastAddress net.IPNet) *Sender {
	return &Sender{senderIP: senderIP, MulticastAddress: multicastAddress}
}

func (sender Sender) Receive() {
	udpAddr, _ := net.ResolveUDPAddr("udp", "239.69.128.213:"+strconv.Itoa(aes67Port))
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

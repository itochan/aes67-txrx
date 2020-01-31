package aes67

import (
	"log"
	"net"

	"github.com/pion/rtp"
)

var connectRx *net.UDPConn

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
	connectRx, err = net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal("net.ListenMulticastUDP()", err)
	}

	receivePacket()

	defer connectRx.Close()
}

func receivePacket() {
	buffer := make([]byte, 156)
	packet := &rtp.Packet{}

	var cnt int
	for {
		_, err := connectRx.Read(buffer)
		if err != nil {
			log.Fatal("connect.Read()", err)
		}
		cnt++
		packet.Unmarshal(buffer)
		RxCh <- packet.SequenceNumber
		// fmt.Printf("Remote receiver got %d packets\n", cnt)
		// fmt.Print(packet.String())
	}
}

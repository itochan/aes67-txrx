package aes67

import (
	"fmt"
	"log"
	"net"

	"github.com/pion/rtp"
)

type RxTx struct {
	// Rx
	localIP  net.IP
	remoteIP net.IP

	// Tx
	senderIP         net.IP
	MulticastAddress net.IPNet
}

func NewRxTx(localIP net.IP, remoteIP net.IP, senderIP net.IP, multicastAddress net.IPNet) *RxTx {
	return &RxTx{localIP: localIP, remoteIP: remoteIP, senderIP: senderIP, MulticastAddress: multicastAddress}
}

func (rxtx RxTx) ReceiveAndSend() {
	var err error

	udpAddr := &net.UDPAddr{IP: rxtx.remoteIP, Port: aes67Port}
	connectRx, err = net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal("net.ListenMulticastUDP()", err)
	}

	dialer := net.Dialer{
		LocalAddr: &net.UDPAddr{IP: rxtx.senderIP, Port: aes67Port + 50000},
	}
	destinationAddr := net.UDPAddr{IP: rxtx.MulticastAddress.IP, Port: aes67Port}
	connectTx, err = dialer.Dial("udp", destinationAddr.String())
	if err != nil {
		log.Fatal("Dial", err)
	}

	receiveAndSendPacket()

	defer connectRx.Close()
	defer connectTx.Close()
}

func receiveAndSendPacket() {
	// const RTPHeader = 12
	// const PCM24bit48kHz = 144

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
		fmt.Printf("Remote receiver got %d packets\n", cnt)
		fmt.Print(packet.String())

		sendPacket(packet)
	}
}

package aes67

import (
	"fmt"
	"net"

	"github.com/itochan/GoRTP/src/net/rtp"
)

type Receiver struct {
	senderIP         net.IP
	MulticastAddress net.IPNet
}

func NewReceiver(senderIP net.IP, multicastAddress net.IPNet) *Sender {
	return &Sender{senderIP: senderIP, MulticastAddress: multicastAddress}
}

func (sender Sender) Receive() {
	local := &net.IPAddr{IP: sender.senderIP}
	// transmitAddr, _ := net.ResolveIPAddr("ip", sender.MulticastAddress.IP.String())
	transmitAddr, _ := net.ResolveIPAddr("ip", "239.69.128.194")

	tpLocal, _ := rtp.NewTransportUDP(local, aes67Port, localZone)

	rsLocal = rtp.NewSession(tpLocal, tpLocal)

	rsLocal.AddRemote(&rtp.Address{transmitAddr.IP, aes67Port, aes67Port + 1, remoteZone})

	strLocalIdx, _ := rsLocal.NewSsrcStreamOut(&rtp.Address{local.IP, aes67Port, aes67Port + 1, localZone}, 0, 0)
	rsLocal.SsrcStreamOutForIndex(strLocalIdx).SetPayloadType(0)

	rsLocal.StartSession()
	receivePacket()
	defer rsLocal.CloseRecv()
}

func receivePacket() {
	// Create and store the data receive channel.
	dataReceiver := rsLocal.CreateDataReceiveChan()
	var cnt int

	for {
		select {
		case rp := <-dataReceiver: // just get a packet - maybe we add some tests later
			if (cnt % 50) == 0 {
				fmt.Printf("Remote receiver got %d packets\n", cnt)
			}
			cnt++
			rp.FreePacket()
		}
	}
}

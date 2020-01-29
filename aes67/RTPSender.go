package aes67

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
)

type Sender struct {
	senderIP         net.IP
	MulticastAddress net.IPNet
}

func NewSender(senderIP net.IP, multicastAddress net.IPNet) *Sender {
	return &Sender{senderIP: senderIP, MulticastAddress: multicastAddress}
}

func (sender Sender) Play(transmitFile string) {
	playFile(transmitFile)
}

func playFile(transmitFile string) {
	const PCM24bit48kHz = 144
	buf := make([]byte, PCM24bit48kHz)

	file, _ := ioutil.ReadFile(transmitFile)
	reader := bytes.NewReader(file)

	packetizer := rtp.NewPacketizer(1452, 97, 0xC1E0F3FB, &codecs.G722Payloader{}, rtp.NewRandomSequencer(), 90000)

	const tickTime = 1 * time.Millisecond
	t := time.NewTicker(tickTime)

	start := time.Now()
	for {
		n, _ := reader.Read(buf)
		if n == 0 {
			break
		}
		packet := packetizer.Packetize(buf, 480)
		select {
		case <-t.C:
			fmt.Print(packet[0].String())
		}
	}
	elapsed := time.Since(start)
	log.Printf("Send RTP Packet %s", elapsed)
}

func sendPacket(payload []byte, stamp uint32) {
	rp := rsLocal.NewDataPacket(stamp)
	rp.SetPayload(payload[:])
	rp.SetPayloadType(PCM24)
	rsLocal.WriteData(rp)
	rp.FreePacket()
}

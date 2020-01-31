package aes67

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
)

var connectTx net.Conn

type Sender struct {
	senderIP         net.IP
	MulticastAddress net.IPNet
}

func NewSender(senderIP net.IP, multicastAddress net.IPNet) *Sender {
	return &Sender{senderIP: senderIP, MulticastAddress: multicastAddress}
}

func (sender Sender) Play(transmitFile string) {
	dialer := net.Dialer{
		LocalAddr: &net.UDPAddr{IP: sender.senderIP, Port: aes67Port + 50000},
	}
	destinationAddr := net.UDPAddr{IP: sender.MulticastAddress.IP, Port: aes67Port}
	var err error
	connectTx, err = dialer.Dial("udp", destinationAddr.String())
	if err != nil {
		log.Fatal("Dial", err)
	}

	playFile(transmitFile)

	defer connectTx.Close()
}

func playFile(transmitFile string) {
	const RTPHeader = 12
	const PCM24bit48kHz = 144

	buf := make([]byte, PCM24bit48kHz)

	file, _ := ioutil.ReadFile(transmitFile)
	reader := bytes.NewReader(file)

	packetizer := rtp.NewPacketizer(RTPHeader+PCM24bit48kHz, PCM24, 0xC1E0F3FB, &codecs.G722Payloader{}, rtp.NewFixedSequencer(1), 90000)

	const tickTime = 1 * time.Millisecond
	t := time.NewTicker(tickTime)

	start := time.Now()
	for {
		n, _ := reader.Read(buf)
		if n == 0 {
			break
		}
		packet := packetizer.Packetize(buf, 48)
		select {
		case <-t.C:
			go sendPacket(packet[0])
		}
	}
	elapsed := time.Since(start)
	log.Printf("Sent RTP Packet %s", elapsed)
}

func sendPacket(packet *rtp.Packet) {
	bytes, _ := packet.Marshal()
	_, err := connectTx.Write(bytes)
	TxCh <- packet.SequenceNumber
	if err != nil {
		log.Print(err)
	}
}

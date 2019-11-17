package aes67

import (
	"github.com/wernerd/GoRTP/src/net/rtp"
	"net"
	"os"
	"time"
)

const (
	aes67Port = 5004
	PCM24     = 0x61
)

var (
	rsLocal    *rtp.Session
	localZone  = ""
	remoteZone = ""
)

type Sender struct {
	senderIP         net.IP
	MulticastAddress net.IPNet
}

func NewSender(senderIP net.IP, multicastAddress net.IPNet) *Sender {
	return &Sender{senderIP: senderIP, MulticastAddress: multicastAddress}
}

func (sender Sender) Play(transmitFile string) {
	local := &net.IPAddr{IP: sender.senderIP}
	transmitAddr, _ := net.ResolveIPAddr("ip", sender.MulticastAddress.IP.String())

	tpLocal, _ := rtp.NewTransportUDP(local, aes67Port, localZone)

	rsLocal = rtp.NewSession(tpLocal, tpLocal)

	rsLocal.AddRemote(&rtp.Address{transmitAddr.IP, aes67Port, aes67Port + 1, remoteZone})

	strLocalIdx, _ := rsLocal.NewSsrcStreamOut(&rtp.Address{local.IP, aes67Port, aes67Port + 1, localZone}, 1020304, 4711)
	rsLocal.SsrcStreamOutForIndex(strLocalIdx).SetPayloadType(0)

	rsLocal.StartSession()
	playFile(transmitFile)
	defer rsLocal.CloseRecv()
}

func playFile(transmitFile string) {
	var cnt int
	stamp := uint32(0)

	const PCM24bit48kHz = 288
	buf := make([]byte, PCM24bit48kHz)

	file, _ := os.Open(transmitFile)
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			file.Close()
			break
		}

		rp := rsLocal.NewDataPacket(stamp)
		rp.SetPayload(buf[:])
		rp.SetPayloadType(PCM24)
		rsLocal.WriteData(rp)
		rp.FreePacket()
		cnt++
		stamp += 48
		time.Sleep(1440 * time.Microsecond)
	}
}

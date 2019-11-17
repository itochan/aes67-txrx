package aes67

import (
	"github.com/wernerd/GoRTP/src/net/rtp"
	"net"
	"strconv"
	"time"
)

const (
	aes67Port = 5004
	PCM24     = 0x61
)

var rsLocal *rtp.Session
var rsRemote *rtp.Session

var localPay [160]byte
var remotePay [160]byte

var stop bool
var stopLocalRecv chan bool
var stopRemoteRecv chan bool
var stopLocalCtrl chan bool
var stopRemoteCtrl chan bool

var eventNamesNew = []string{"NewStreamData", "NewStreamCtrl"}
var eventNamesRtcp = []string{"SR", "RR", "SDES", "BYE"}

var localZone = ""
var remoteZone = ""

type Sender struct {
	senderIP         net.IP
	MulticastAddress net.IPNet
}

func NewSender(senderIP net.IP, multicastAddress net.IPNet) *Sender {
	return &Sender{senderIP: senderIP, MulticastAddress: multicastAddress}
}

func (sender Sender) PlayFile(transmitFile string) {
	initialize()
	localZone := ""
	remoteZone := ""

	local := &net.IPAddr{IP: sender.senderIP}
	//local, _ := net.ResolveIPAddr("ip", "127.0.0.1")
	transmitAddr, _ := net.ResolveIPAddr("ip", sender.MulticastAddress.IP.String())

	tpLocal, _ := rtp.NewTransportUDP(local, aes67Port, localZone)

	rsLocal = rtp.NewSession(tpLocal, tpLocal)

	rsLocal.AddRemote(&rtp.Address{transmitAddr.IP, aes67Port, aes67Port + 1, remoteZone})

	strLocalIdx, _ := rsLocal.NewSsrcStreamOut(&rtp.Address{local.IP, aes67Port, aes67Port + 1, localZone}, 1020304, 4711)
	rsLocal.SsrcStreamOutForIndex(strLocalIdx).SetPayloadType(0)

	//fmt.Println("RTP session opened")
	//
	//rsLocal.NewDataPacket()
	//_, _ = tpLocal.WriteDataTo(sendAddr)
	//tpLocal.CloseWrite()
	//
	//fmt.Println("Closing session...")
	//defer rsLocal.CloseSession()

	rsLocal.StartSession()

	sendLocalToRemote()
	defer rsLocal.CloseRecv()
}

func initialize() {
	var localPay [160]byte
	for i := range localPay {
		localPay[i] = byte(i)
	}
}

func sendLocalToRemote() {
	var cnt int
	stamp := uint32(0)
	for {
		rp := rsLocal.NewDataPacket(stamp)
		rp.SetPayload(localPay[:])
		rp.SetPayloadType(PCM24)
		rsLocal.WriteData(rp)
		rp.Print(strconv.Itoa(cnt))
		rp.FreePacket()
		//if (cnt % 50) == 0 {
		//	fmt.Printf("Local sent %d packets\n", cnt)
		//}
		cnt++
		stamp += 160
		time.Sleep(20e6)
	}
}

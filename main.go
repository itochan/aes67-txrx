package main

import (
	"flag"
	"log"
	"net"
	"sync"

	"github.com/itochan/aes67-txrx/aes67"
	"github.com/itochan/aes67-txrx/sap"
)

var (
	mode          = flag.String("m", "transmit", "Mode")
	interfaceName = flag.String("i", "", "Network interface")
	transmitFile  = flag.String("f", "", "Transmit File")
	address       = flag.String("a", "", "Receive address")
)

func main() {
	flag.Parse()

	sap := sap.NewSAP(*interfaceName)
	switch *mode {
	case "receive":
		log.Printf("IP: %s", net.ParseIP(*address))
		r := aes67.NewReceiver(sap.HostAddress, net.ParseIP(*address))
		r.Receive()
		break
	case "transmit":
		sap.AnnounceSAP()
		s := aes67.NewSender(sap.HostAddress, sap.MulticastAddress)
		s.Play(*transmitFile)
		break
	case "txrx":
		sap.AnnounceSAP()

		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			log.Printf("Start Transmitter")
			s := aes67.NewSender(sap.HostAddress, sap.MulticastAddress)
			s.Play(*transmitFile)
			wg.Done()
		}()
		go func() {
			log.Printf("IP: %s", net.ParseIP(*address))
			log.Printf("Start Receiver")
			r := aes67.NewReceiver(sap.HostAddress, net.ParseIP(*address))
			r.Receive()
			wg.Done()
		}()
		wg.Wait()

		// start := time.Now()
		// elapsed := time.Since(start)
		// log.Printf("Sent RTP Packet %s", elapsed)
		break
	case "rxtx":
		rxtx := aes67.NewRxTx(sap.HostAddress, net.ParseIP(*address), sap.HostAddress, sap.MulticastAddress)
		rxtx.ReceiveAndSend()
		break
	}
}

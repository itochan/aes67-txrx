package main

import (
	"flag"
	"net"

	"github.com/itochan/aes67-txrx/sap"
	"github.com/itochan/aes67-txrx/aes67"
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
		r := aes67.NewReceiver(sap.HostAddress, net.ParseIP(*address))
		r.Receive()
		break
	case "transmit":
		sap.AnnounceSAP()
		s := aes67.NewSender(sap.HostAddress, sap.MulticastAddress)
		s.Play(*transmitFile)
		break
	}
}

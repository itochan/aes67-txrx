package main

import (
	"flag"
	"github.com/itochan/aes67-transmitter/aes67"
	"github.com/itochan/aes67-transmitter/sap"
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
		r := aes67.NewReceiver(sap.HostAddress, sap.MulticastAddress)
		r.Receive()
		break
	case "transmit":
		sap.AnnounceSAP()
		s := aes67.NewSender(sap.HostAddress, sap.MulticastAddress)
		s.Play(*transmitFile)
		break
	}
}

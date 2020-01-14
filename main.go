package main

import (
	"flag"
	"github.com/itochan/aes67-transmitter/aes67"
	"github.com/itochan/aes67-transmitter/sap"
)

var (
	interfaceName = flag.String("i", "", "Network interface")
	transmitFile  = flag.String("f", "", "Transmit File")
)

func main() {
	flag.Parse()

	sap := sap.NewSAP(*interfaceName)
	sap.AnnounceSAP()

	r := aes67.NewSender(sap.HostAddress, sap.MulticastAddress)
	r.Play(*transmitFile)
}

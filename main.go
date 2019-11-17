package main

import (
	"flag"
	"github.com/itochan/aes67-transmitter/sap"
)

var (
	interfaceName = flag.String("i", "", "Network interface")
)

func main() {
	flag.Parse()

	sap := sap.NewSAP(*interfaceName)
	sap.AnnounceSAP()
}

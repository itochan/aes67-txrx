package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/itochan/aes67-txrx/aes67"
	"github.com/itochan/aes67-txrx/sap"
)

var (
	mode          = flag.String("m", "transmit", "Mode")
	interfaceName = flag.String("i", "", "Network interface")
	transmitFile  = flag.String("f", "", "Transmit File")
	address       = flag.String("a", "", "Receive address")
)

type TxRxLog struct {
	txTime time.Time
	rxTime time.Time
	rtt    time.Duration
}

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

		chEnd := make(chan struct{})
		go func() {
			log.Printf("Start Transmitter")
			s := aes67.NewSender(sap.HostAddress, sap.MulticastAddress)
			s.Play(*transmitFile)
			time.Sleep(20 * time.Millisecond)
			close(chEnd)
		}()
		go func() {
			log.Printf("IP: %s", net.ParseIP(*address))
			log.Printf("Start Receiver")
			r := aes67.NewReceiver(sap.HostAddress, net.ParseIP(*address))
			r.Receive()
		}()

		txRxLog := make([]TxRxLog, 30000)

		for {
			select {
			case txSequenceNo := <-aes67.TxCh:
				txRxLog[txSequenceNo-1].txTime = time.Now()
			case rxSequenceNo := <-aes67.RxCh:
				txRxLog[rxSequenceNo-1].rxTime = time.Now()
				txRxLog[rxSequenceNo-1].rtt = txRxLog[rxSequenceNo-1].rxTime.Sub(txRxLog[rxSequenceNo-1].txTime)
			case <-chEnd:
				file, err := os.Create(fmt.Sprintf("%s.csv", time.Now()))
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				for i, log := range txRxLog {
					line := fmt.Sprintf("%d,%d\n", i+1, log.rtt/time.Nanosecond)
					file.WriteString(line)
				}
				return
			}
		}
	case "rxtx":
		rxtx := aes67.NewRxTx(sap.HostAddress, net.ParseIP(*address), sap.HostAddress, sap.MulticastAddress)
		chEnd := make(chan struct{})
		go func() {
			rxtx.ReceiveAndSend()
			close(chEnd)
		}()

		for {
			select {
			case <-aes67.TxCh:
				break
			case <-aes67.RxCh:
				break
			case <-chEnd:
				return
			}
		}
	}
}

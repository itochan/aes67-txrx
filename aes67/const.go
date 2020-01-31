package aes67

const (
	aes67Port = 5004
	PCM24     = 0x61
)

var (
	localZone  = ""
	remoteZone = ""

	TxCh = make(chan uint16)
	RxCh = make(chan uint16)
)

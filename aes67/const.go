package aes67

import "github.com/itochan/GoRTP/src/net/rtp"

const (
	aes67Port = 5004
	PCM24     = 0x61
)

var (
	rsLocal    *rtp.Session
	localZone  = ""
	remoteZone = ""
)

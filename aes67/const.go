package aes67

import "net"

const (
	aes67Port = 5004
	PCM24     = 0x61
)

var (
	connect    net.Conn
	localZone  = ""
	remoteZone = ""
)

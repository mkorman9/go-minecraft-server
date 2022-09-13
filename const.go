package main

import "time"

var (
	SegmentBits = 0x7F
	ContinueBit = 0x80

	ProtocolName          = "1.19"
	ProtocolVersion       = 759
	MaxPacketSize         = 2097151
	ServerKeyLength       = 1024
	VerifyTokenLength     = 8
	KeepAliveSendInterval = 5 * time.Second
)

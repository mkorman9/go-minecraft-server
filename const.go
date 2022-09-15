package main

import "time"

var (
	ProtocolName          = "1.19"
	ProtocolVersion       = 759
	ServerKeyLength       = 1024
	VerifyTokenLength     = 16
	KeepAliveSendInterval = 5 * time.Second
)

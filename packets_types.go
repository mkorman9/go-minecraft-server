package main

type HandshakeType = int

const (
	HandshakeTypeStatus = 1
	HandshakeTypeLogin  = 2
)

type SystemChatMessageType = int

const (
	SystemChatMessageTypeChat     = 0
	SystemChatMessageTypeSystem   = 1
	SystemChatMessageTypeGameInfo = 2
)

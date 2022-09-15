package main

const (
	HandshakeTypeStatus = 1
	HandshakeTypeLogin  = 2
)

const (
	SystemChatMessageTypeChat     = 0
	SystemChatMessageTypeSystem   = 1
	SystemChatMessageTypeGameInfo = 2
)

type GameMode = byte

const (
	GameModeSurvival  byte = 0
	GameModeCreative  byte = 1
	GameModeAdventure byte = 2
	GameModeUnknown   byte = 255
)

type EntityAction = int

const (
	EntityActionStartSneaking         = 0
	EntityActionStopSneaking          = 1
	EntityActionLeaveBed              = 2
	EntityActionStartSprinting        = 3
	EntityActionStopSprinting         = 4
	EntityActionStartJumpWithHorse    = 5
	EntityActionStopJumpWithHorse     = 6
	EntityActionOpenHorseInventory    = 7
	EntityActionStartFlyingWithElytra = 8
)

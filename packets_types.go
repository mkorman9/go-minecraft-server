package main

import "github.com/mkorman9/go-minecraft-server/nbt"

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

type GameMode = byte

const (
	GameModeSurvival  = 0
	GameModeCreative  = 1
	GameModeAdventure = 2
	GameModeUnknown   = 255
)

type SlotData struct {
	Present   bool
	ItemID    int
	ItemCount byte
	NBT       nbt.RawMessage
}

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

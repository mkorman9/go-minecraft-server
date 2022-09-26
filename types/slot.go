package types

import "github.com/mkorman9/go-minecraft-server/nbt"

type SlotData struct {
	Present   bool
	ItemID    int
	ItemCount byte
	NBT       *nbt.RawMessage
}

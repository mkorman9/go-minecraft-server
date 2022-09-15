package packets

const (
	TypeArray = iota
	TypeByte
	TypeBool
	TypeInt16
	TypeInt32
	TypeInt64
	TypeVarInt
	TypeFloat32
	TypeFloat64
	TypeUUID
	TypeVarLong
	TypeByteArray
	TypeString
	TypeNBT
	TypePosition
	TypeSlot
	TypeBitSet
)

type Field struct {
	Type                   int
	Value                  any
	Blueprint              any
	ArrayElementDefinition *PacketDefinition
	ArrayLengthOption      ArrayLengthOption
}

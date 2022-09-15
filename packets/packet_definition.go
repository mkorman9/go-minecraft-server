package packets

import (
	"io"
)

type PacketDefinition struct {
	PacketID int
	Fields   []Field

	namesMapping map[string]int
}

func Packet(opts ...PacketOpt) *PacketDefinition {
	packet := PacketDefinition{
		PacketID: -1,
	}

	for _, opt := range opts {
		opt(&packet)
	}

	return &packet
}

func (pd *PacketDefinition) AddField(name string, fieldType int) {
	pd.Fields = append(pd.Fields, Field{Type: fieldType})
	pd.namesMapping[name] = len(pd.Fields) - 1
}

func (pd *PacketDefinition) SetBlueprint(name string, blueprint any) {
	if i, ok := pd.namesMapping[name]; ok {
		pd.Fields[i].Blueprint = blueprint
	}
}

func (pd *PacketDefinition) New() *PacketData {
	return &PacketData{
		PacketID:     pd.PacketID,
		Fields:       pd.Fields[:],
		namesMapping: pd.namesMapping,
	}
}

func (pd *PacketDefinition) Read(reader io.Reader) (PacketData, error) {
	packet := pd.New()

	for _, field := range packet.Fields {
		var err error

		switch field.Type {
		case TypeArray:
			var length int
			if field.ArrayLengthOption == ArrayLengthAppend {
				length, err = readVarInt(reader)
				if err != nil {
					return PacketData{}, err
				}
			}

			elements := make(ArrayValue, length)
			for i := 0; i < length; i++ {
				element, err := field.ArrayElementDefinition.Read(reader)
				if err != nil {
					return PacketData{}, err
				}

				elements[i] = element
			}
			field.Value = elements
		case TypeByte:
			var value byte
			value, err = readByte(reader)
			field.Value = value
		case TypeBool:
			var value bool
			value, err = readBool(reader)
			field.Value = value
		case TypeInt16:
			var value int16
			value, err = readInt16(reader)
			field.Value = value
		case TypeInt32:
			var value int32
			value, err = readInt32(reader)
			field.Value = value
		case TypeInt64:
			var value int64
			value, err = readInt64(reader)
			field.Value = value
		case TypeVarInt:
			var value int
			value, err = readVarInt(reader)
			field.Value = value
		case TypeFloat32:
			var value float32
			value, err = readFloat32(reader)
			field.Value = value
		case TypeFloat64:
			var value float64
			value, err = readFloat64(reader)
			field.Value = value
		case TypeUUID:
			var value UUID
			value, err = readUUID(reader)
			field.Value = value
		case TypeVarLong:
			var value int64
			value, err = readVarLong(reader)
			field.Value = value
		case TypeByteArray:
			var value []byte
			value, err = readByteArray(reader)
			field.Value = value
		case TypeString:
			var value string
			value, err = readString(reader)
			field.Value = value
		case TypeNBT:
			var value any
			value, err = readNBT(reader, field.Blueprint)
			field.Value = value
		case TypePosition:
			var value Position
			value, err = readPosition(reader)
			field.Value = value
		case TypeSlot:
			var value SlotData
			value, err = readSlot(reader)
			field.Value = value
		case TypeBitSet:
			var value BitSet
			value, err = readBitSet(reader)
			field.Value = value
		}

		if err != nil {
			return PacketData{}, err
		}
	}

	return PacketData{}, nil
}

func (pd *PacketDefinition) specifyArrayOptions(name string, elementDefinition *PacketDefinition, lengthOptions ArrayLengthOption) {
	if i, ok := pd.namesMapping[name]; ok {
		pd.Fields[i].ArrayElementDefinition = elementDefinition
		pd.Fields[i].ArrayLengthOption = lengthOptions
	}
}

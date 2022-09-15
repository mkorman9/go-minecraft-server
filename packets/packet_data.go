package packets

import (
	"io"
)

type PacketData struct {
	PacketID int
	Fields   []*Field

	namesMapping map[string]int
}

func (pd *PacketData) Set(name string, value any) *PacketData {
	if i, ok := pd.namesMapping[name]; ok {
		pd.Fields[i].Value = value
	}

	return pd
}

func (pd *PacketData) SetArray(name string, value ConvertedArrayValue) *PacketData {
	if i, ok := pd.namesMapping[name]; ok {
		pd.Fields[i].Value = value(pd.Fields[i].ArrayElementDefinition)
	}

	return pd
}

func (pd *PacketData) WriteTo(writer io.Writer) (int64, error) {
	if pd.PacketID != -1 {
		err := writeVarInt(writer, pd.PacketID)
		if err != nil {
			return 0, err
		}
	}

	for _, field := range pd.Fields {
		var err error

		cancelField := false
		for _, opt := range field.FieldOptions {
			if !opt(pd) {
				cancelField = true
			}
		}
		if cancelField {
			continue
		}

		switch field.Type {
		case TypeArray:
			if field.Value == nil {
				if field.ArrayLengthOption == ArrayLengthAppend {
					err := writeVarInt(writer, 0)
					if err != nil {
						return 0, err
					}
				}

				continue
			}

			array := field.Value.(ArrayValue)

			if field.ArrayLengthOption == ArrayLengthAppend {
				err := writeVarInt(writer, len(array))
				if err != nil {
					return 0, err
				}
			}

			for _, element := range array {
				_, err := element.WriteTo(writer)
				if err != nil {
					return 0, err
				}
			}
		case TypeByte:
			err = writeByte(writer, field.Value.(byte))
		case TypeBool:
			err = writeBool(writer, field.Value.(bool))
		case TypeInt16:
			err = writeInt16(writer, field.Value.(int16))
		case TypeInt32:
			err = writeInt32(writer, field.Value.(int32))
		case TypeInt64:
			err = writeInt64(writer, field.Value.(int64))
		case TypeVarInt:
			err = writeVarInt(writer, field.Value.(int))
		case TypeFloat32:
			err = writeFloat32(writer, field.Value.(float32))
		case TypeFloat64:
			err = writeFloat64(writer, field.Value.(float64))
		case TypeUUID:
			err = writeUUID(writer, field.Value.(UUID))
		case TypeVarLong:
			err = writeVarLong(writer, field.Value.(int64))
		case TypeByteArray:
			err = writeByteArray(writer, field.Value.([]byte))
		case TypeString:
			err = writeString(writer, field.Value.(string))
		case TypeNBT:
			err = writeNBT(writer, field.Value)
		case TypePosition:
			err = writePosition(writer, field.Value.(*Position))
		case TypeSlot:
			err = writeSlot(writer, field.Value.(*SlotData))
		case TypeBitSet:
			err = writeBitSet(writer, field.Value.(*BitSet))
		}

		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}

func (pd *PacketData) Array(name string) ArrayValue {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeByte {
			return pd.Fields[i].Value.(ArrayValue)
		}
	}

	return nil
}

func (pd *PacketData) Byte(name string) byte {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeByte {
			return pd.Fields[i].Value.(byte)
		}
	}

	return 0
}

func (pd *PacketData) Bool(name string) bool {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeBool {
			return pd.Fields[i].Value.(bool)
		}
	}

	return false
}

func (pd *PacketData) Int16(name string) int16 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeInt16 {
			return pd.Fields[i].Value.(int16)
		}
	}

	return 0
}

func (pd *PacketData) Int32(name string) int32 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeInt32 {
			return pd.Fields[i].Value.(int32)
		}
	}

	return 0
}

func (pd *PacketData) Int64(name string) int64 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeInt64 {
			return pd.Fields[i].Value.(int64)
		}
	}

	return 0
}

func (pd *PacketData) VarInt(name string) int {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeVarInt {
			return pd.Fields[i].Value.(int)
		}
	}

	return 0
}

func (pd *PacketData) Float32(name string) float32 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeFloat32 {
			return pd.Fields[i].Value.(float32)
		}
	}

	return 0
}

func (pd *PacketData) Float64(name string) float64 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeFloat64 {
			return pd.Fields[i].Value.(float64)
		}
	}

	return 0
}

func (pd *PacketData) UUID(name string) UUID {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeUUID {
			return pd.Fields[i].Value.(UUID)
		}
	}

	return UUID{0, 0}
}

func (pd *PacketData) VarLong(name string) int64 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeVarLong {
			return pd.Fields[i].Value.(int64)
		}
	}

	return 0
}

func (pd *PacketData) ByteArray(name string) []byte {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeByteArray {
			return pd.Fields[i].Value.([]byte)
		}
	}

	return nil
}

func (pd *PacketData) String(name string) string {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeString {
			return pd.Fields[i].Value.(string)
		}
	}

	return ""
}

func (pd *PacketData) NBT(name string) any {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeNBT {
			return pd.Fields[i].Value
		}
	}

	return nil
}

func (pd *PacketData) Position(name string) Position {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypePosition {
			return pd.Fields[i].Value.(Position)
		}
	}

	return Position{}
}

func (pd *PacketData) Slot(name string) SlotData {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeSlot {
			return pd.Fields[i].Value.(SlotData)
		}
	}

	return SlotData{}
}

func (pd *PacketData) BitSet(name string) BitSet {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeBitSet {
			return pd.Fields[i].Value.(BitSet)
		}
	}

	return BitSet{}
}

package packets

import (
	"github.com/mkorman9/go-minecraft-server/types"
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
		err := types.WriteVarInt(writer, pd.PacketID)
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
			length := field.ArrayLengthOption(pd)

			if field.Value == nil {
				if length == arrayLengthPrefix {
					err := types.WriteVarInt(writer, 0)
					if err != nil {
						return 0, err
					}
				}

				continue
			}

			array := field.Value.(ArrayValue)

			if length == arrayLengthPrefix {
				err := types.WriteVarInt(writer, len(array))
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
			err = types.WriteByte(writer, field.Value.(byte))
		case TypeBool:
			err = types.WriteBool(writer, field.Value.(bool))
		case TypeInt16:
			err = types.WriteInt16(writer, field.Value.(int16))
		case TypeInt32:
			err = types.WriteInt32(writer, field.Value.(int32))
		case TypeInt64:
			err = types.WriteInt64(writer, field.Value.(int64))
		case TypeVarInt:
			err = types.WriteVarInt(writer, field.Value.(int))
		case TypeFloat32:
			err = types.WriteFloat32(writer, field.Value.(float32))
		case TypeFloat64:
			err = types.WriteFloat64(writer, field.Value.(float64))
		case TypeUUID:
			err = types.WriteUUID(writer, field.Value.(types.UUID))
		case TypeVarLong:
			err = types.WriteVarLong(writer, field.Value.(int64))
		case TypeByteArray:
			err = types.WriteByteArray(writer, field.Value.([]byte))
		case TypeString:
			err = types.WriteString(writer, field.Value.(string))
		case TypeNBT:
			err = types.WriteNBT(writer, field.Value)
		case TypePosition:
			err = types.WritePosition(writer, field.Value.(*types.Position))
		case TypeSlot:
			err = types.WriteSlot(writer, field.Value.(*types.SlotData))
		case TypeBitSet:
			err = types.WriteBitSet(writer, field.Value.(*types.BitSet))
		}

		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}

func (pd *PacketData) Any(name string) any {
	if i, ok := pd.namesMapping[name]; ok {
		return pd.Fields[i].Value
	}

	return nil
}

func (pd *PacketData) Array(name string) ArrayValue {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeArray {
			if value, ok := pd.Fields[i].Value.(ArrayValue); ok {
				return value
			}
		}
	}

	return nil
}

func (pd *PacketData) Byte(name string) byte {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeByte {
			if value, ok := pd.Fields[i].Value.(byte); ok {
				return value
			}
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
			if value, ok := pd.Fields[i].Value.(int16); ok {
				return value
			} else if value, ok := pd.Fields[i].Value.(byte); ok {
				return int16(value)
			}
		}
	}

	return 0
}

func (pd *PacketData) Int32(name string) int32 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeInt32 {
			if value, ok := pd.Fields[i].Value.(int32); ok {
				return value
			} else if value, ok := pd.Fields[i].Value.(int16); ok {
				return int32(value)
			} else if value, ok := pd.Fields[i].Value.(byte); ok {
				return int32(value)
			}
		}
	}

	return 0
}

func (pd *PacketData) Int64(name string) int64 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeInt64 {
			if value, ok := pd.Fields[i].Value.(int64); ok {
				return value
			} else if value, ok := pd.Fields[i].Value.(int32); ok {
				return int64(value)
			} else if value, ok := pd.Fields[i].Value.(int16); ok {
				return int64(value)
			}
		}
	}

	return 0
}

func (pd *PacketData) VarInt(name string) int {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeVarInt {
			if value, ok := pd.Fields[i].Value.(int); ok {
				return value
			}
		}
	}

	return 0
}

func (pd *PacketData) Float32(name string) float32 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeFloat32 {
			if value, ok := pd.Fields[i].Value.(float32); ok {
				return value
			}
		}
	}

	return 0
}

func (pd *PacketData) Float64(name string) float64 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeFloat64 {
			if value, ok := pd.Fields[i].Value.(float64); ok {
				return value
			} else if value, ok := pd.Fields[i].Value.(float32); ok {
				return float64(value)
			}
		}
	}

	return 0
}

func (pd *PacketData) UUID(name string) types.UUID {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeUUID {
			if value, ok := pd.Fields[i].Value.(types.UUID); ok {
				return value
			}
		}
	}

	return types.UUID{0, 0}
}

func (pd *PacketData) VarLong(name string) int64 {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeVarLong {
			if value, ok := pd.Fields[i].Value.(int64); ok {
				return value
			}
		}
	}

	return 0
}

func (pd *PacketData) ByteArray(name string) []byte {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeByteArray {
			if value, ok := pd.Fields[i].Value.([]byte); ok {
				return value
			} else if value, ok := pd.Fields[i].Value.(string); ok {
				return []byte(value)
			}
		}
	}

	return nil
}

func (pd *PacketData) String(name string) string {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeString {
			if value, ok := pd.Fields[i].Value.(string); ok {
				return value
			} else if value, ok := pd.Fields[i].Value.([]byte); ok {
				return string(value)
			}
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

func (pd *PacketData) Position(name string) *types.Position {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypePosition {
			if value, ok := pd.Fields[i].Value.(*types.Position); ok {
				return value
			}
		}
	}

	return nil
}

func (pd *PacketData) Slot(name string) *types.SlotData {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeSlot {
			if value, ok := pd.Fields[i].Value.(*types.SlotData); ok {
				return value
			}
		}
	}

	return nil
}

func (pd *PacketData) BitSet(name string) *types.BitSet {
	if i, ok := pd.namesMapping[name]; ok {
		if pd.Fields[i].Type == TypeBitSet {
			if value, ok := pd.Fields[i].Value.(*types.BitSet); ok {
				return value
			}
		}
	}

	return nil
}

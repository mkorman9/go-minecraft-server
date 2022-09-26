package packets

import (
	"errors"
	"github.com/mkorman9/go-minecraft-server/types"
	"io"
)

type PacketDefinition struct {
	PacketID int
	Fields   []*Field

	namesMapping map[string]int
}

func Packet(opts ...PacketOpt) *PacketDefinition {
	packet := PacketDefinition{
		PacketID:     -1,
		namesMapping: make(map[string]int),
	}

	for _, opt := range opts {
		opt(&packet)
	}

	return &packet
}

func (pd *PacketDefinition) AddField(name string, fieldType int) {
	pd.Fields = append(pd.Fields, &Field{Type: fieldType})
	pd.namesMapping[name] = len(pd.Fields) - 1
}

func (pd *PacketDefinition) SetBlueprint(name string, blueprint any) {
	if i, ok := pd.namesMapping[name]; ok {
		pd.Fields[i].Blueprint = blueprint
	}
}

func (pd *PacketDefinition) New() *PacketData {
	fields := make([]*Field, len(pd.Fields))
	copy(fields, pd.Fields)

	namesMapping := make(map[string]int)
	for name, field := range pd.namesMapping {
		namesMapping[name] = field
	}

	return &PacketData{
		PacketID:     pd.PacketID,
		Fields:       fields,
		namesMapping: namesMapping,
	}
}

func (pd *PacketDefinition) Read(reader io.Reader) (*PacketData, error) {
	packet := pd.New()

	for _, field := range packet.Fields {
		var err error

		cancelField := false
		for _, opt := range field.FieldOptions {
			if !opt(packet) {
				cancelField = true
			}
		}
		if cancelField {
			continue
		}

		switch field.Type {
		case TypeArray:
			length := field.ArrayLengthOption(packet)

			if length == arrayLengthPrefix {
				length, err = types.ReadVarInt(reader)
				if err != nil {
					return nil, err
				}
			}

			elements := make(ArrayValue, length)
			for i := 0; i < length; i++ {
				element, err := field.ArrayElementDefinition.Read(reader)
				if err != nil {
					return nil, err
				}

				elements[i] = *element
			}
			field.Value = elements
		case TypeByte:
			var value byte
			value, err = types.ReadByte(reader)
			field.Value = value
		case TypeBool:
			var value bool
			value, err = types.ReadBool(reader)
			field.Value = value
		case TypeInt16:
			var value int16
			value, err = types.ReadInt16(reader)
			field.Value = value
		case TypeInt32:
			var value int32
			value, err = types.ReadInt32(reader)
			field.Value = value
		case TypeInt64:
			var value int64
			value, err = types.ReadInt64(reader)
			field.Value = value
		case TypeVarInt:
			var value int
			value, err = types.ReadVarInt(reader)
			field.Value = value
		case TypeFloat32:
			var value float32
			value, err = types.ReadFloat32(reader)
			field.Value = value
		case TypeFloat64:
			var value float64
			value, err = types.ReadFloat64(reader)
			field.Value = value
		case TypeUUID:
			var value types.UUID
			value, err = types.ReadUUID(reader)
			field.Value = value
		case TypeVarLong:
			var value int64
			value, err = types.ReadVarLong(reader)
			field.Value = value
		case TypeByteArray:
			var value []byte
			value, err = types.ReadByteArray(reader)
			field.Value = value
		case TypeString:
			var value string
			value, err = types.ReadString(reader)
			field.Value = value
		case TypeNBT:
			var value any
			value, err = types.ReadNBT(reader, field.Blueprint)
			field.Value = value
		case TypePosition:
			var value types.Position
			value, err = types.ReadPosition(reader)
			field.Value = value
		case TypeSlot:
			var value types.SlotData
			value, err = types.ReadSlot(reader)
			field.Value = value
		case TypeBitSet:
			var value types.BitSet
			value, err = types.ReadBitSet(reader)
			field.Value = value
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				return packet, nil
			}

			return nil, err
		}
	}

	return packet, nil
}

func (pd *PacketDefinition) specifyArrayOptions(name string, elementDefinition *PacketDefinition, lengthOptions ArrayLengthOption) {
	if i, ok := pd.namesMapping[name]; ok {
		pd.Fields[i].ArrayElementDefinition = elementDefinition
		pd.Fields[i].ArrayLengthOption = lengthOptions
	}
}

func (pd *PacketDefinition) setFieldOpts(name string, opts []PacketFieldOpt) {
	if i, ok := pd.namesMapping[name]; ok {
		pd.Fields[i].FieldOptions = opts
	}
}

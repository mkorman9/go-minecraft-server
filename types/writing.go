package types

import (
	"encoding/binary"
	"github.com/mkorman9/go-minecraft-server/nbt"
	"io"
	"reflect"
)

func WriteByte(writer io.Writer, value byte) error {
	_, err := writer.Write([]byte{value})
	return err
}

func WriteBool(writer io.Writer, value bool) error {
	var b byte
	if value {
		b = 1
	}

	return WriteByte(writer, b)
}

func WriteInt16(writer io.Writer, value int16) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func WriteInt32(writer io.Writer, value int32) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func WriteInt64(writer io.Writer, value int64) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func WriteFloat32(writer io.Writer, value float32) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func WriteFloat64(writer io.Writer, value float64) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func WriteUUID(writer io.Writer, value UUID) error {
	err := binary.Write(writer, binary.BigEndian, value.Upper)
	if err != nil {
		return err
	}

	return binary.Write(writer, binary.BigEndian, value.Lower)
}

func WriteVarInt(writer io.Writer, value int) error {
	for {
		if (value & ^SegmentBits) == 0 {
			err := WriteByte(writer, byte(value))
			if err != nil {
				return err
			}

			break
		}

		err := WriteByte(writer, byte((value&SegmentBits)|ContinueBit))
		if err != nil {
			return err
		}

		value >>= 7
	}

	return nil
}

func WriteVarLong(writer io.Writer, value int64) error {
	for {
		if (value & ^int64(SegmentBits)) == 0 {
			err := WriteByte(writer, byte(value))
			if err != nil {
				return err
			}

			break
		}

		err := WriteByte(writer, byte((value&int64(SegmentBits))|int64(ContinueBit)))
		if err != nil {
			return err
		}

		value >>= 7
	}

	return nil
}

func WriteByteArray(writer io.Writer, value []byte) error {
	err := WriteVarInt(writer, len(value))
	if err != nil {
		return err
	}

	_, err = writer.Write(value)
	return err
}

func WriteString(writer io.Writer, value string) error {
	return WriteByteArray(writer, []byte(value))
}

func WriteNBT(writer io.Writer, obj any) error {
	if obj == nil || reflect.ValueOf(obj).IsNil() {
		return WriteByte(writer, nbt.TagEnd)
	}

	data, err := nbt.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func WritePosition(writer io.Writer, position *Position) error {
	return WriteInt64(writer, position.ToInt64())
}

func WriteSlot(writer io.Writer, slot *SlotData) error {
	err := WriteBool(writer, slot.Present)
	if err != nil {
		return err
	}

	if slot.Present {
		err = WriteVarInt(writer, slot.ItemID)
		if err != nil {
			return err
		}

		err = WriteByte(writer, slot.ItemCount)
		if err != nil {
			return err
		}

		err = WriteNBT(writer, slot.NBT)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteBitSet(writer io.Writer, bitSet *BitSet) error {
	err := WriteVarInt(writer, len(bitSet.v))
	if err != nil {
		return err
	}

	for _, v := range bitSet.v {
		err = WriteInt64(writer, v)
		if err != nil {
			return err
		}
	}

	return nil
}

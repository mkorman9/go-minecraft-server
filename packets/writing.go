package packets

import (
	"encoding/binary"
	"github.com/mkorman9/go-minecraft-server/nbt"
	"io"
	"reflect"
)

func writeByte(writer io.Writer, value byte) error {
	_, err := writer.Write([]byte{value})
	return err
}

func writeBool(writer io.Writer, value bool) error {
	var b byte
	if value {
		b = 1
	}

	return writeByte(writer, b)
}

func writeInt16(writer io.Writer, value int16) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func writeInt32(writer io.Writer, value int32) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func writeInt64(writer io.Writer, value int64) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func writeFloat32(writer io.Writer, value float32) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func writeFloat64(writer io.Writer, value float64) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func writeUUID(writer io.Writer, value UUID) error {
	err := binary.Write(writer, binary.BigEndian, value.Upper)
	if err != nil {
		return err
	}

	return binary.Write(writer, binary.BigEndian, value.Lower)
}

func writeVarInt(writer io.Writer, value int) error {
	for {
		if (value & ^SegmentBits) == 0 {
			err := writeByte(writer, byte(value))
			if err != nil {
				return err
			}

			break
		}

		err := writeByte(writer, byte((value&SegmentBits)|ContinueBit))
		if err != nil {
			return err
		}

		value >>= 7
	}

	return nil
}

func writeVarLong(writer io.Writer, value int64) error {
	for {
		if (value & ^int64(SegmentBits)) == 0 {
			err := writeByte(writer, byte(value))
			if err != nil {
				return err
			}

			break
		}

		err := writeByte(writer, byte((value&int64(SegmentBits))|int64(ContinueBit)))
		if err != nil {
			return err
		}

		value >>= 7
	}

	return nil
}

func writeByteArray(writer io.Writer, value []byte) error {
	err := writeVarInt(writer, len(value))
	if err != nil {
		return err
	}

	_, err = writer.Write(value)
	return err
}

func writeString(writer io.Writer, value string) error {
	return writeByteArray(writer, []byte(value))
}

func writeNBT(writer io.Writer, obj any) error {
	if obj == nil || reflect.ValueOf(obj).IsNil() {
		return writeByte(writer, nbt.TagEnd)
	}

	data, err := nbt.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func writePosition(writer io.Writer, position *Position) error {
	return writeInt64(writer, position.ToInt64())
}

func writeSlot(writer io.Writer, slot *SlotData) error {
	err := writeBool(writer, slot.Present)
	if err != nil {
		return err
	}

	if slot.Present {
		err = writeVarInt(writer, slot.ItemID)
		if err != nil {
			return err
		}

		err = writeByte(writer, slot.ItemCount)
		if err != nil {
			return err
		}

		err = writeNBT(writer, slot.NBT)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeBitSet(writer io.Writer, bitSet *BitSet) error {
	err := writeVarInt(writer, len(bitSet.v))
	if err != nil {
		return err
	}

	for _, v := range bitSet.v {
		err = writeInt64(writer, v)
		if err != nil {
			return err
		}
	}

	return nil
}

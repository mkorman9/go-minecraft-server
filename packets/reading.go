package packets

import (
	"encoding/binary"
	"errors"
	"github.com/mkorman9/go-minecraft-server/nbt"
	"io"
	"math"
	"reflect"
)

func readByte(reader io.Reader) (byte, error) {
	buff := make([]byte, 1)

	_, err := reader.Read(buff)
	if err != nil {
		return 0, err
	}

	return buff[0], nil
}

func readBytes(reader io.Reader, n int) ([]byte, error) {
	buff := make([]byte, n)
	_, err := reader.Read(buff)
	return buff, err
}

func readBool(reader io.Reader) (bool, error) {
	value, err := readByte(reader)
	if err != nil {
		return false, err
	}

	if value > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func readInt16(reader io.Reader) (int16, error) {
	b, err := readBytes(reader, 2)
	if err != nil {
		return 0, err
	}

	return int16(binary.BigEndian.Uint16(b)), nil
}

func readInt32(reader io.Reader) (int32, error) {
	b, err := readBytes(reader, 4)
	if err != nil {
		return 0, err
	}

	return int32(binary.BigEndian.Uint32(b)), nil
}

func readInt64(reader io.Reader) (int64, error) {
	b, err := readBytes(reader, 8)
	if err != nil {
		return 0, err
	}

	return int64(binary.BigEndian.Uint64(b)), nil
}

func readVarInt(reader io.Reader) (int, error) {
	var value int
	var position int

	for {
		currentByte, err := readByte(reader)
		if err != nil {
			return 0, err
		}

		value |= int(currentByte) & SegmentBits << position

		if (int(currentByte) & ContinueBit) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, errors.New("invalid size of VarInt")
		}
	}

	return value, nil
}

func readVarLong(reader io.Reader) (int64, error) {
	var value int64
	var position int64

	for {
		currentByte, err := readByte(reader)
		if err != nil {
			return 0, err
		}

		value |= int64(currentByte) & int64(SegmentBits) << position

		if (int(currentByte) & ContinueBit) == 0 {
			break
		}

		position += 7

		if position >= 64 {
			return 0, errors.New("invalid size of VarLong")
		}
	}

	return value, nil
}

func readFloat32(reader io.Reader) (float32, error) {
	value, err := readInt32(reader)
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(uint32(value)), nil
}

func readFloat64(reader io.Reader) (float64, error) {
	value, err := readInt64(reader)
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(uint64(value)), nil
}

func readUUID(reader io.Reader) (UUID, error) {
	upper, err := readInt64(reader)
	if err != nil {
		return UUID{}, err
	}

	lower, err := readInt64(reader)
	if err != nil {
		return UUID{}, err
	}

	return UUID{Upper: upper, Lower: lower}, nil
}

func readByteArray(reader io.Reader) ([]byte, error) {
	length, err := readVarInt(reader)
	if err != nil {
		return nil, err
	}

	buff := make([]byte, length)
	_, err = reader.Read(buff)
	return buff, err
}

func readString(reader io.Reader) (string, error) {
	b, err := readByteArray(reader)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func readNBT(reader io.Reader, blueprint any) (any, error) {
	obj := reflect.New(reflect.TypeOf(blueprint))

	_, err := nbt.NewDecoder(reader).Decode(&obj)
	if err != nil {
		if errors.Is(err, nbt.ErrEND) {
			return nil, nil
		}

		return nil, err
	}

	return obj, nil
}

func readPosition(reader io.Reader) (Position, error) {
	value, err := readInt64(reader)
	if err != nil {
		return Position{}, err
	}

	return PositionFromInt64(value), nil
}

func readSlot(reader io.Reader) (slot SlotData, err error) {
	slot.Present, err = readBool(reader)
	if err != nil {
		return
	}

	if slot.Present {
		slot.ItemID, err = readVarInt(reader)
		if err != nil {
			return
		}

		slot.ItemCount, err = readByte(reader)
		if err != nil {
			return
		}

		tags, e := readNBT(reader, &nbt.RawMessage{})
		if e != nil {
			err = e
			return
		}

		if tags != nil {
			slot.NBT = tags.(*nbt.RawMessage)
		}
	}

	return
}

func readBitSet(reader io.Reader) (BitSet, error) {
	length, err := readVarInt(reader)
	if err != nil {
		return BitSet{}, err
	}

	var bitSet BitSet
	for i := 0; i < length; i++ {
		value, err := readInt64(reader)
		if err != nil {
			return BitSet{}, err
		}

		bitSet.v = append(bitSet.v, value)
	}

	return bitSet, nil
}

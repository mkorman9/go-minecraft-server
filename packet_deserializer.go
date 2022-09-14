package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/mkorman9/go-minecraft-server/nbt"
	"io"
	"math"
)

type PacketDeserializer struct {
	data   []byte
	cursor int
	err    error
}

func NewPacketDeserializer(data []byte) *PacketDeserializer {
	return &PacketDeserializer{
		data:   data,
		cursor: 0,
	}
}

func (pd *PacketDeserializer) Error() error {
	return pd.err
}

func (pd *PacketDeserializer) BytesLeft() int {
	return len(pd.data) - pd.cursor
}

func (pd *PacketDeserializer) FetchByte() byte {
	if pd.BytesLeft() <= 0 {
		pd.err = errors.New("out of bounds read")
		return 0
	}

	value := pd.data[pd.cursor]
	pd.cursor++
	return value
}

func (pd *PacketDeserializer) FetchBool() bool {
	return pd.FetchByte() > 0
}

func (pd *PacketDeserializer) FetchInt16() int16 {
	return int16(binary.BigEndian.Uint16([]byte{pd.FetchByte(), pd.FetchByte()}))
}

func (pd *PacketDeserializer) FetchInt32() int {
	return int(binary.BigEndian.Uint32([]byte{pd.FetchByte(), pd.FetchByte(), pd.FetchByte(), pd.FetchByte()}))
}

func (pd *PacketDeserializer) FetchInt64() int64 {
	return int64(binary.BigEndian.Uint64([]byte{
		pd.FetchByte(), pd.FetchByte(), pd.FetchByte(), pd.FetchByte(),
		pd.FetchByte(), pd.FetchByte(), pd.FetchByte(), pd.FetchByte(),
	}))
}

func (pd *PacketDeserializer) FetchFloat32() float32 {
	value := pd.FetchInt32()
	return math.Float32frombits(uint32(value))
}

func (pd *PacketDeserializer) FetchFloat64() float64 {
	value := pd.FetchInt64()
	return math.Float64frombits(uint64(value))
}

func (pd *PacketDeserializer) FetchUUID() UUID {
	return UUID{
		Upper: pd.FetchInt64(),
		Lower: pd.FetchInt64(),
	}
}

func (pd *PacketDeserializer) FetchVarInt() int {
	var value int
	var position int
	var currentByte byte

	for {
		currentByte = pd.FetchByte()
		value |= int(currentByte) & SegmentBits << position

		if (int(currentByte) & ContinueBit) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return -1
		}
	}

	return value
}

func (pd *PacketDeserializer) FetchVarLong() int64 {
	var value int64
	var position int64
	var currentByte byte

	for {
		currentByte = pd.FetchByte()
		value |= int64(currentByte) & int64(SegmentBits) << position

		if (int(currentByte) & ContinueBit) == 0 {
			break
		}

		position += 7

		if position >= 64 {
			return -1
		}
	}

	return value
}

func (pd *PacketDeserializer) FetchByteArray() []byte {
	length := pd.FetchVarInt()
	if pd.cursor+length > len(pd.data) {
		pd.err = errors.New("out of bounds read")
		return nil
	}

	value := pd.data[pd.cursor : pd.cursor+length]
	pd.cursor += length
	return value
}

func (pd *PacketDeserializer) FetchString() string {
	return string(pd.FetchByteArray())
}

func (pd *PacketDeserializer) FetchNBT(v any) {
	buff := pd.data[pd.cursor:]
	reader := &io.LimitedReader{R: bytes.NewBuffer(buff), N: math.MaxInt64}

	_, err := nbt.NewDecoder(reader).Decode(v)
	if err != nil {
		if !errors.Is(err, nbt.ErrEND) {
			pd.err = err
			return
		}
	}

	pd.cursor += int(math.MaxInt64 - reader.N)
}

func (pd *PacketDeserializer) FetchPosition() *Position {
	value := pd.FetchInt64()
	return PositionFromInt64(value)
}

func (pd *PacketDeserializer) FetchSlot() *SlotData {
	var slot SlotData

	slot.Present = pd.FetchBool()
	if slot.Present {
		slot.ItemID = pd.FetchVarInt()
		slot.ItemCount = pd.FetchByte()

		slot.NBT = nbt.RawMessage{Type: nbt.TagEnd}
		pd.FetchNBT(&slot.NBT)
	}

	return &slot
}

func (pd *PacketDeserializer) FetchBitSet() *BitSet {
	length := pd.FetchVarInt()

	var bitSet BitSet
	for i := 0; i < length; i++ {
		value := pd.FetchInt64()
		bitSet.v = append(bitSet.v, value)
	}

	return &bitSet
}

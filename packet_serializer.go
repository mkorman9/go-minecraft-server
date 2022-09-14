package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"github.com/mkorman9/go-minecraft-server/nbt"
	"reflect"
)

type PacketSerializer struct {
	buffer               *bytes.Buffer
	err                  error
	compressionThreshold int
}

func NewPackerSerializer(compressionThreshold int) *PacketSerializer {
	return &PacketSerializer{
		buffer:               bytes.NewBuffer(make([]byte, 0)),
		compressionThreshold: compressionThreshold,
	}
}

func (ps *PacketSerializer) Error() error {
	return ps.err
}

func (ps *PacketSerializer) Bytes() []byte {
	data := ps.buffer.Bytes()
	dataSize := len(data)
	finalWriter := &PacketSerializer{buffer: bytes.NewBuffer(make([]byte, 0)), compressionThreshold: -1}

	if ps.compressionThreshold >= 0 {
		if dataSize >= ps.compressionThreshold {
			var zlibBuffer bytes.Buffer

			zlibWriter := zlib.NewWriter(&zlibBuffer)
			_, _ = zlibWriter.Write(ps.buffer.Bytes())
			_ = zlibWriter.Close()

			finalWriter.AppendVarInt(getVarIntSize(dataSize) + zlibBuffer.Len())
			finalWriter.AppendVarInt(dataSize)
			finalWriter.buffer.Write(zlibBuffer.Bytes())
		} else {
			finalWriter.AppendVarInt(dataSize + 1)
			finalWriter.AppendVarInt(0)
			finalWriter.buffer.Write(data)
		}
	} else {
		finalWriter.AppendVarInt(dataSize)
		finalWriter.buffer.Write(data)
	}

	return finalWriter.buffer.Bytes()
}

func (ps *PacketSerializer) AppendByte(value byte) {
	ps.wrapError(
		ps.buffer.WriteByte(value),
	)
}

func (ps *PacketSerializer) AppendBool(value bool) {
	var b byte
	if value {
		b = 1
	}

	ps.AppendByte(b)
}

func (ps *PacketSerializer) AppendInt16(value int16) {
	ps.wrapError(
		binary.Write(ps.buffer, binary.BigEndian, value),
	)
}

func (ps *PacketSerializer) AppendInt32(value int32) {
	ps.wrapError(
		binary.Write(ps.buffer, binary.BigEndian, value),
	)
}

func (ps *PacketSerializer) AppendInt64(value int64) {
	ps.wrapError(
		binary.Write(ps.buffer, binary.BigEndian, value),
	)
}

func (ps *PacketSerializer) AppendFloat32(value float32) {
	ps.wrapError(
		binary.Write(ps.buffer, binary.BigEndian, value),
	)
}

func (ps *PacketSerializer) AppendFloat64(value float64) {
	ps.wrapError(
		binary.Write(ps.buffer, binary.BigEndian, value),
	)
}

func (ps *PacketSerializer) AppendUUID(value UUID) {
	ps.wrapError(
		binary.Write(ps.buffer, binary.BigEndian, value.Upper),
	)
	ps.wrapError(
		binary.Write(ps.buffer, binary.BigEndian, value.Lower),
	)
}

func (ps *PacketSerializer) AppendVarInt(value int) {
	for {
		if (value & ^SegmentBits) == 0 {
			ps.AppendByte(byte(value))
			break
		}

		ps.AppendByte(byte((value & SegmentBits) | ContinueBit))

		value >>= 7
	}
}

func (ps *PacketSerializer) AppendVarLong(value int64) {
	for {
		if (value & ^int64(SegmentBits)) == 0 {
			ps.AppendByte(byte(value))
			break
		}

		ps.AppendByte(byte((value & int64(SegmentBits)) | int64(ContinueBit)))

		value >>= 7
	}
}

func (ps *PacketSerializer) AppendByteArray(value []byte) {
	ps.AppendVarInt(len(value))
	_, err := ps.buffer.Write(value)
	ps.wrapError(err)
}

func (ps *PacketSerializer) AppendString(value string) {
	ps.AppendByteArray([]byte(value))
}

func (ps *PacketSerializer) AppendNBT(obj any) {
	if obj == nil || reflect.ValueOf(obj).IsNil() {
		ps.AppendByte(nbt.TagEnd)
		return
	}

	data, err := nbt.Marshal(obj)
	ps.wrapError(err)

	_, err = ps.buffer.Write(data)
	ps.wrapError(err)
}

func (ps *PacketSerializer) AppendPosition(position *Position) {
	ps.AppendInt64(position.ToInt64())
}

func (ps *PacketSerializer) AppendSlot(slot *SlotData) {
	ps.AppendBool(slot.Present)
	if slot.Present {
		ps.AppendVarInt(slot.ItemID)
		ps.AppendByte(slot.ItemCount)
		ps.AppendNBT(&slot.NBT)
	}
}

func (ps *PacketSerializer) AppendBitSet(bitSet *BitSet) {
	ps.AppendVarInt(len(bitSet.v))
	for _, v := range bitSet.v {
		ps.AppendInt64(v)
	}
}

func (ps *PacketSerializer) wrapError(err error) {
	if err != nil {
		ps.err = err
	}
}

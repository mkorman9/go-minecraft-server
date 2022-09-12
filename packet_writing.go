package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"github.com/mkorman9/go-minecraft-server/nbt"
	"reflect"
)

type PacketWriter struct {
	compressionThreshold int
}

type PacketWriterContext struct {
	buffer               *bytes.Buffer
	err                  error
	compressionThreshold int
}

func NewPacketWriter() *PacketWriter {
	return &PacketWriter{compressionThreshold: -1}
}

func (pw *PacketWriter) New() *PacketWriterContext {
	return &PacketWriterContext{
		buffer:               bytes.NewBuffer(make([]byte, 0)),
		compressionThreshold: pw.compressionThreshold,
	}
}

func (pw *PacketWriter) EnableCompression(threshold int) {
	pw.compressionThreshold = threshold
}

func (pwc *PacketWriterContext) Error() error {
	return pwc.err
}

func (pwc *PacketWriterContext) Bytes() []byte {
	data := pwc.buffer.Bytes()
	dataSize := len(data)
	finalWriter := &PacketWriterContext{buffer: bytes.NewBuffer(make([]byte, 0)), compressionThreshold: -1}

	if pwc.compressionThreshold >= 0 {
		if dataSize >= pwc.compressionThreshold {
			var zlibBuffer bytes.Buffer

			zlibWriter := zlib.NewWriter(&zlibBuffer)
			_, _ = zlibWriter.Write(pwc.buffer.Bytes())
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

func (pwc *PacketWriterContext) AppendByte(value byte) {
	pwc.wrapError(
		pwc.buffer.WriteByte(value),
	)
}

func (pwc *PacketWriterContext) AppendBool(value bool) {
	var b byte
	if value {
		b = 1
	}

	pwc.AppendByte(b)
}

func (pwc *PacketWriterContext) AppendInt16(value int16) {
	pwc.wrapError(
		binary.Write(pwc.buffer, binary.BigEndian, value),
	)
}

func (pwc *PacketWriterContext) AppendInt32(value int32) {
	pwc.wrapError(
		binary.Write(pwc.buffer, binary.BigEndian, value),
	)
}

func (pwc *PacketWriterContext) AppendInt64(value int64) {
	pwc.wrapError(
		binary.Write(pwc.buffer, binary.BigEndian, value),
	)
}

func (pwc *PacketWriterContext) AppendFloat32(value float32) {
	pwc.wrapError(
		binary.Write(pwc.buffer, binary.BigEndian, value),
	)
}

func (pwc *PacketWriterContext) AppendFloat64(value float64) {
	pwc.wrapError(
		binary.Write(pwc.buffer, binary.BigEndian, value),
	)
}

func (pwc *PacketWriterContext) AppendUUID(value UUID) {
	pwc.wrapError(
		binary.Write(pwc.buffer, binary.BigEndian, value.Upper),
	)
	pwc.wrapError(
		binary.Write(pwc.buffer, binary.BigEndian, value.Lower),
	)
}

func (pwc *PacketWriterContext) AppendVarInt(value int) {
	for {
		if (value & ^SegmentBits) == 0 {
			pwc.AppendByte(byte(value))
			break
		}

		pwc.AppendByte(byte((value & SegmentBits) | ContinueBit))

		value >>= 7
	}
}

func (pwc *PacketWriterContext) AppendVarLong(value int64) {
	for {
		if (value & ^int64(SegmentBits)) == 0 {
			pwc.AppendByte(byte(value))
			break
		}

		pwc.AppendByte(byte((value & int64(SegmentBits)) | int64(ContinueBit)))

		value >>= 7
	}
}

func (pwc *PacketWriterContext) AppendByteArray(value []byte) {
	pwc.AppendVarInt(len(value))
	_, err := pwc.buffer.Write(value)
	pwc.wrapError(err)
}

func (pwc *PacketWriterContext) AppendString(value string) {
	pwc.AppendByteArray([]byte(value))
}

func (pwc *PacketWriterContext) AppendNBT(obj any) {
	if obj == nil || reflect.ValueOf(obj).IsNil() {
		pwc.AppendByte(nbt.TagEnd)
		return
	}

	data, err := nbt.Marshal(obj)
	pwc.wrapError(err)

	_, err = pwc.buffer.Write(data)
	pwc.wrapError(err)
}

func (pwc *PacketWriterContext) AppendPosition(position *Position) {
	pwc.AppendInt64(position.ToInt64())
}

func (pwc *PacketWriterContext) wrapError(err error) {
	if err != nil {
		pwc.err = err
	}
}

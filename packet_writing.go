package main

import (
	"bytes"
	"encoding/binary"
)

type PacketWriter struct {
	buffer *bytes.Buffer
}

func NewPacketWriter() *PacketWriter {
	return &PacketWriter{buffer: bytes.NewBuffer(make([]byte, 0, 1))}
}

func (pw *PacketWriter) Bytes() []byte {
	data := pw.buffer.Bytes()
	finalWriter := NewPacketWriter()

	finalWriter.AppendVarInt(len(data))
	finalWriter.buffer.Write(data)

	return finalWriter.buffer.Bytes()
}

func (pw *PacketWriter) AppendByte(value byte) {
	_ = pw.buffer.WriteByte(value)
}

func (pw *PacketWriter) AppendBool(value bool) {
	var b byte
	if value {
		b = 1
	}

	pw.AppendByte(b)
}

func (pw *PacketWriter) AppendInt16(value int16) {
	_ = binary.Write(pw.buffer, binary.BigEndian, value)
}

func (pw *PacketWriter) AppendInt32(value int) {
	_ = binary.Write(pw.buffer, binary.BigEndian, value)
}

func (pw *PacketWriter) AppendInt64(value int64) {
	_ = binary.Write(pw.buffer, binary.BigEndian, value)
}

func (pw *PacketWriter) AppendUUID(value UUID) {
	_ = binary.Write(pw.buffer, binary.BigEndian, value.Upper)
	_ = binary.Write(pw.buffer, binary.BigEndian, value.Lower)
}

func (pw *PacketWriter) AppendVarInt(value int) {
	for {
		if (value & ^SegmentBits) == 0 {
			pw.AppendByte(byte(value))
			break
		}

		pw.AppendByte(byte((value & SegmentBits) | ContinueBit))

		value >>= 7
	}
}

func (pw *PacketWriter) AppendVarLong(value int64) {
	for {
		if (value & ^int64(SegmentBits)) == 0 {
			pw.AppendByte(byte(value))
			break
		}

		pw.AppendByte(byte((value & int64(SegmentBits)) | int64(ContinueBit)))

		value >>= 7
	}
}

func (pw *PacketWriter) AppendString(value string) {
	pw.AppendVarInt(len(value))
	pw.buffer.Write([]byte(value))
}

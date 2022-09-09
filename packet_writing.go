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

func (pw *PacketWriter) AppendInt16(value int16) {
	_ = binary.Write(pw.buffer, binary.BigEndian, value)
}

func (pw *PacketWriter) AppendInt32(value int) {
	_ = binary.Write(pw.buffer, binary.BigEndian, value)
}

func (pw *PacketWriter) AppendInt64(value int64) {
	_ = binary.Write(pw.buffer, binary.BigEndian, value)
}

func (pw *PacketWriter) AppendVarInt(value int) {
	for {
		if (value & ^SEGMENT_BITS) == 0 {
			pw.AppendByte(byte(value))
			break
		}

		pw.AppendByte(byte((value & SEGMENT_BITS) | CONTINUE_BIT))

		value >>= 7
	}
}

func (pw *PacketWriter) AppendVarLong(value int64) {
	for {
		if (value & ^int64(SEGMENT_BITS)) == 0 {
			pw.AppendByte(byte(value))
			break
		}

		pw.AppendByte(byte((value & int64(SEGMENT_BITS)) | int64(CONTINUE_BIT)))

		value >>= 7
	}
}

func (pw *PacketWriter) AppendString(value string) {
	pw.AppendVarInt(len(value))
	pw.buffer.Write([]byte(value))
}

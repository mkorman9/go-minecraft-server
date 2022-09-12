package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/mkorman9/go-minecraft-server/nbt"
	"io"
	"math"
)

type PacketReaderContext struct {
	data   []byte
	cursor int
	err    error
}

func NewPacketReaderContext(data []byte) *PacketReaderContext {
	return &PacketReaderContext{
		data:   data,
		cursor: 0,
	}
}

func (prc *PacketReaderContext) Error() error {
	return prc.err
}

func (prc *PacketReaderContext) BytesLeft() int {
	return len(prc.data) - prc.cursor
}

func (prc *PacketReaderContext) FetchByte() byte {
	if prc.BytesLeft() <= 0 {
		prc.err = errors.New("out of bounds read")
		return 0
	}

	value := prc.data[prc.cursor]
	prc.cursor++
	return value
}

func (prc *PacketReaderContext) FetchBool() bool {
	return prc.FetchByte() > 0
}

func (prc *PacketReaderContext) FetchInt16() int16 {
	return int16(binary.BigEndian.Uint16([]byte{prc.FetchByte(), prc.FetchByte()}))
}

func (prc *PacketReaderContext) FetchInt32() int {
	return int(binary.BigEndian.Uint32([]byte{prc.FetchByte(), prc.FetchByte(), prc.FetchByte(), prc.FetchByte()}))
}

func (prc *PacketReaderContext) FetchInt64() int64 {
	return int64(binary.BigEndian.Uint64([]byte{
		prc.FetchByte(), prc.FetchByte(), prc.FetchByte(), prc.FetchByte(),
		prc.FetchByte(), prc.FetchByte(), prc.FetchByte(), prc.FetchByte(),
	}))
}

func (prc *PacketReaderContext) FetchUUID() UUID {
	return UUID{
		Upper: prc.FetchInt64(),
		Lower: prc.FetchInt64(),
	}
}

func (prc *PacketReaderContext) FetchVarInt() int {
	var value int
	var position int
	var currentByte byte

	for {
		currentByte = prc.FetchByte()
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

func (prc *PacketReaderContext) FetchVarLong() int64 {
	var value int64
	var position int64
	var currentByte byte

	for {
		currentByte = prc.FetchByte()
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

func (prc *PacketReaderContext) FetchByteArray() []byte {
	length := prc.FetchVarInt()
	if prc.cursor+length > len(prc.data) {
		prc.err = errors.New("out of bounds read")
		return nil
	}

	value := prc.data[prc.cursor : prc.cursor+length]
	prc.cursor += length
	return value
}

func (prc *PacketReaderContext) FetchString() string {
	return string(prc.FetchByteArray())
}

func (prc *PacketReaderContext) FetchNBT(v any) {
	buff := prc.data[prc.cursor:]
	reader := &io.LimitedReader{R: bytes.NewBuffer(buff), N: math.MaxInt64}

	_, err := nbt.NewDecoder(reader).Decode(v)
	if err != nil {
		prc.err = err
		return
	}

	prc.cursor += int(math.MaxInt64 - reader.N)
}

func (prc *PacketReaderContext) FetchPosition() *Position {
	value := prc.FetchInt64()
	return PositionFromInt64(value)
}

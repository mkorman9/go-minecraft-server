package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/mkorman9/go-minecraft-server/nbt"
	"io"
	"math"
)

type PacketReader struct {
	data   []byte
	cursor int
}

func (pr *PacketReader) BytesLeft() int {
	return len(pr.data) - pr.cursor
}

func (pr *PacketReader) FetchByte() byte {
	value := pr.data[pr.cursor]
	pr.cursor++
	return value
}

func (pr *PacketReader) FetchBool() bool {
	return pr.FetchByte() > 0
}

func (pr *PacketReader) FetchInt16() int16 {
	return int16(binary.BigEndian.Uint16([]byte{pr.FetchByte(), pr.FetchByte()}))
}

func (pr *PacketReader) FetchInt32() int {
	return int(binary.BigEndian.Uint32([]byte{pr.FetchByte(), pr.FetchByte(), pr.FetchByte(), pr.FetchByte()}))
}

func (pr *PacketReader) FetchInt64() int64 {
	return int64(binary.BigEndian.Uint64([]byte{
		pr.FetchByte(), pr.FetchByte(), pr.FetchByte(), pr.FetchByte(),
		pr.FetchByte(), pr.FetchByte(), pr.FetchByte(), pr.FetchByte(),
	}))
}

func (pr *PacketReader) FetchUUID() UUID {
	return UUID{
		Upper: pr.FetchInt64(),
		Lower: pr.FetchInt64(),
	}
}

func (pr *PacketReader) FetchVarInt() int {
	var value int
	var position int
	var currentByte byte

	for {
		currentByte = pr.FetchByte()
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

func (pr *PacketReader) FetchVarLong() int64 {
	var value int64
	var position int64
	var currentByte byte

	for {
		currentByte = pr.FetchByte()
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

func (pr *PacketReader) FetchString() string {
	length := pr.FetchVarInt()
	value := string(pr.data[pr.cursor : pr.cursor+length])
	pr.cursor += length
	return value
}

func (pr *PacketReader) FetchNBT(v any) {
	buff := pr.data[pr.cursor:]
	reader := &io.LimitedReader{R: bytes.NewBuffer(buff), N: math.MaxInt64}

	_, _ = nbt.NewDecoder(reader).Decode(v)
	pr.cursor += int(math.MaxInt64 - reader.N)
}

func (pr *PacketReader) FetchPosition() *Position {
	value := pr.FetchInt64()
	return PositionFromInt64(value)
}

func ReadPacketSize(reader io.Reader) (int, error) {
	var value int
	var position int
	var currentByte byte

	for {
		buff := make([]byte, 1)
		_, err := reader.Read(buff)
		if err != nil {
			return -1, err
		}

		currentByte = buff[0]
		value |= int(currentByte) & SegmentBits << position

		if (int(currentByte) & ContinueBit) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return -1, errors.New("invalid VarInt size")
		}
	}

	return value, nil
}

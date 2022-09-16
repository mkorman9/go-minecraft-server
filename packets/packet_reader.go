package packets

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
)

type PacketDelivery struct {
	PacketID int
	Header   *PacketHeader
	Reader   io.Reader
}

type PacketReader struct {
	reader               io.Reader
	compressionThreshold int
}

type PacketHeader struct {
	PacketSize           int
	UncompressedDataSize int
	UseCompression       bool
}

func NewPacketReader(reader io.Reader) *PacketReader {
	return &PacketReader{
		reader:               reader,
		compressionThreshold: -1,
	}
}

func (pr *PacketReader) SetCompression(threshold int) {
	pr.compressionThreshold = threshold
}

func (pr *PacketReader) SetEncryption(cipherStream *CipherStream) {
	pr.reader = cipherStream.WrapReader(pr.reader)
}

func (pr *PacketReader) Read() (*PacketDelivery, error) {
	header, err := pr.readHeader()
	if err != nil {
		return nil, err
	}

	packetData := make([]byte, header.PacketSize)
	_, err = pr.reader.Read(packetData)
	if err != nil {
		return nil, err
	}

	if header.UseCompression {
		zlibReader, err := zlib.NewReader(bytes.NewReader(packetData))
		if err != nil {
			return nil, err
		}

		zlibBuffer := make([]byte, header.UncompressedDataSize)
		_, err = zlibReader.Read(zlibBuffer)
		if err != nil {
			return nil, err
		}

		packetData = zlibBuffer
	}

	packetReader := bytes.NewBuffer(packetData)

	packetId, err := readVarInt(packetReader)
	if err != nil {
		return nil, err
	}

	return &PacketDelivery{
		PacketID: packetId,
		Header:   header,
		Reader:   packetReader,
	}, nil
}

func (pr *PacketReader) readHeader() (*PacketHeader, error) {
	switch pr.compressionThreshold {
	case -1:
		// no compression
		packetSize, err := readVarInt(pr.reader)
		if err != nil {
			return nil, err
		}

		if packetSize > MaxPacketSize {
			return nil, errors.New("invalid packet size")
		}

		return &PacketHeader{
			PacketSize:           packetSize,
			UncompressedDataSize: 0,
			UseCompression:       false,
		}, nil
	default:
		// compression
		compressedDataSize, err := readVarInt(pr.reader)
		if err != nil {
			return nil, err
		}

		uncompressedDataSize, err := readVarInt(pr.reader)
		if err != nil {
			return nil, err
		}

		if compressedDataSize > MaxPacketSize || uncompressedDataSize > MaxPacketSize {
			return nil, errors.New("invalid packet size")
		}

		compressedDataSize -= getVarIntSize(uncompressedDataSize)

		return &PacketHeader{
			PacketSize:           compressedDataSize,
			UncompressedDataSize: uncompressedDataSize,
			UseCompression:       uncompressedDataSize != 0,
		}, nil
	}
}

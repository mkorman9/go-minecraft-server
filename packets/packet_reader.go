package packets

import (
	"bytes"
	"compress/zlib"
	"errors"
	"github.com/mkorman9/go-minecraft-server/types"
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
	buffered             bool
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
		buffered:             false,
	}
}

func (pr *PacketReader) SetCompression(threshold int) {
	pr.compressionThreshold = threshold
}

func (pr *PacketReader) SetEncryption(cipherStream *CipherStream) {
	pr.reader = cipherStream.WrapReader(pr.reader)
}

func (pr *PacketReader) SetBuffered(buffered bool) {
	pr.buffered = buffered
}

func (pr *PacketReader) Read() (*PacketDelivery, error) {
	header, err := pr.readHeader()
	if err != nil {
		return nil, err
	}

	var packetReader io.Reader

	if pr.buffered {
		reader, err := pr.readBuffered(header)
		if err != nil {
			return nil, err
		}

		packetReader = reader
	} else {
		reader, err := pr.readerUnbuffered(header)
		if err != nil {
			return nil, err
		}

		packetReader = reader
	}

	packetId, err := types.ReadVarInt(packetReader)
	if err != nil {
		return nil, err
	}

	return &PacketDelivery{
		PacketID: packetId,
		Header:   header,
		Reader:   packetReader,
	}, nil
}

func (pr *PacketReader) readerUnbuffered(header *PacketHeader) (io.Reader, error) {
	var packetReader io.Reader

	if header.UseCompression {
		zlibReader, err := zlib.NewReader(pr.reader)
		if err != nil {
			return nil, err
		}

		packetReader = io.LimitReader(zlibReader, int64(header.UncompressedDataSize))
	} else {
		packetReader = io.LimitReader(pr.reader, int64(header.PacketSize))
	}

	return packetReader, nil
}

func (pr *PacketReader) readBuffered(header *PacketHeader) (io.Reader, error) {
	packetData := make([]byte, header.PacketSize)
	_, err := pr.reader.Read(packetData)
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
	return packetReader, nil
}

func (pr *PacketReader) readHeader() (*PacketHeader, error) {
	switch pr.compressionThreshold {
	case -1:
		// no compression
		packetSize, err := types.ReadVarInt(pr.reader)
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
		compressedDataSize, err := types.ReadVarInt(pr.reader)
		if err != nil {
			return nil, err
		}

		uncompressedDataSize, err := types.ReadVarInt(pr.reader)
		if err != nil {
			return nil, err
		}

		if compressedDataSize > MaxPacketSize || uncompressedDataSize > MaxPacketSize {
			return nil, errors.New("invalid packet size")
		}

		compressedDataSize -= types.GetVarIntSize(uncompressedDataSize)

		return &PacketHeader{
			PacketSize:           compressedDataSize,
			UncompressedDataSize: uncompressedDataSize,
			UseCompression:       uncompressedDataSize != 0,
		}, nil
	}
}

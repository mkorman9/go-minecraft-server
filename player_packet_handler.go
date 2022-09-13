package main

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type PlayerState = int

const (
	PlayerStateBeforeHandshake = iota
	PlayerStateLogin
	PlayerStateEncryption
	PlayerStatePlay
)

type PacketHandlingError struct {
	wrapped error
	reason  *ChatMessage
}

func NewPacketHandlingError(err error, reason *ChatMessage) *PacketHandlingError {
	return &PacketHandlingError{
		wrapped: err,
		reason:  reason,
	}
}

func (phe *PacketHandlingError) Error() string {
	return phe.wrapped.Error()
}

type PlayerPacketHandler struct {
	player       *Player
	world        *World
	connection   net.Conn
	state        PlayerState
	reader       io.Reader
	writer       io.Writer
	packetWriter *PacketWriter

	ip                          string
	verifyToken                 string
	sharedSecret                []byte
	serverHash                  string
	enabledCompressionThreshold int

	canceled      bool
	canceledMutex sync.Mutex
}

func NewPlayerPacketHandler(player *Player, world *World, connection net.Conn, ip string) *PlayerPacketHandler {
	return &PlayerPacketHandler{
		player:                      player,
		world:                       world,
		connection:                  connection,
		state:                       PlayerStateBeforeHandshake,
		reader:                      connection,
		writer:                      connection,
		packetWriter:                NewPacketWriter(),
		ip:                          ip,
		enabledCompressionThreshold: -1,
		canceled:                    false,
		canceledMutex:               sync.Mutex{},
	}
}

func (pph *PlayerPacketHandler) ReadLoop() {
	defer func() {
		pph.world.PlayerList().UnregisterPlayer(pph.player)
		pph.Cancel(nil)
	}()

	for {
		packetMetadata, err := pph.readPacketMetadata()
		if err != nil {
			if err == IgnorablePacketReadError {
				break
			}

			log.Printf("%v\n", err)
			break
		}

		packetData := make([]byte, packetMetadata.PacketSize)
		_, err = pph.reader.Read(packetData)
		if err != nil {
			log.Printf("%v\n", err)
			break
		}

		if packetMetadata.UseCompression {
			zlibReader, err := zlib.NewReader(bytes.NewReader(packetData))
			if err != nil {
				log.Printf("%v\n", err)
				break
			}

			zlibBuffer := make([]byte, packetMetadata.UncompressedDataSize)
			_, err = zlibReader.Read(zlibBuffer)
			if err != nil && err != io.EOF {
				log.Printf("%v\n", err)
				break
			}

			packetData = zlibBuffer
		}

		err = pph.HandlePacket(packetData)
		if err != nil {
			if handlingError, ok := err.(*PacketHandlingError); ok {
				pph.Cancel(handlingError.reason)
			}

			log.Printf("%v\n", err)
			break
		}
	}
}

func (pph *PlayerPacketHandler) setupEncryption() error {
	cipherStream, err := NewCipherStream(pph.sharedSecret)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	pph.reader = cipherStream.WrapReader(pph.reader)
	pph.writer = cipherStream.WrapWriter(pph.writer)

	return nil
}

func (pph *PlayerPacketHandler) setupCompression() error {
	compressionThreshold := pph.world.Settings().CompressionThreshold

	if compressionThreshold >= 0 {
		err := pph.sendSetCompressionRequest(compressionThreshold)
		if err != nil {
			return err
		}

		pph.packetWriter.EnableCompression(compressionThreshold)
		pph.enabledCompressionThreshold = compressionThreshold
	}

	return nil
}

func (pph *PlayerPacketHandler) readPacketMetadata() (*PacketMetadata, error) {
	switch pph.enabledCompressionThreshold {
	case -1:
		// no compression
		packetSize, err := pph.peekVarInt()
		if err != nil {
			if err == io.EOF {
				return nil, IgnorablePacketReadError
			}
			if netOpError, ok := err.(*net.OpError); ok {
				if netOpError.Err.Error() == "use of closed network connection" {
					return nil, IgnorablePacketReadError
				}
			}

			return nil, err
		}

		if packetSize > MaxPacketSize {
			return nil, errors.New("invalid packet size")
		}

		return &PacketMetadata{
			PacketSize:           packetSize,
			UncompressedDataSize: 0,
			UseCompression:       false,
		}, nil
	default:
		// compression
		compressedDataSize, err := pph.peekVarInt()
		if err != nil {
			if err == io.EOF {
				return nil, IgnorablePacketReadError
			}
			if netOpError, ok := err.(*net.OpError); ok {
				if netOpError.Err.Error() == "use of closed network connection" {
					return nil, IgnorablePacketReadError
				}
			}

			return nil, err
		}

		uncompressedDataSize, err := pph.peekVarInt()
		if err != nil {
			if err == io.EOF {
				return nil, IgnorablePacketReadError
			}
			if netOpError, ok := err.(*net.OpError); ok {
				if netOpError.Err.Error() == "use of closed network connection" {
					return nil, IgnorablePacketReadError
				}
			}

			return nil, err
		}

		if compressedDataSize > MaxPacketSize || uncompressedDataSize > MaxPacketSize {
			return nil, errors.New("invalid packet size")
		}

		compressedDataSize -= getVarIntSize(uncompressedDataSize)

		return &PacketMetadata{
			PacketSize:           compressedDataSize,
			UncompressedDataSize: uncompressedDataSize,
			UseCompression:       uncompressedDataSize != 0,
		}, nil
	}
}

func (pph *PlayerPacketHandler) peekVarInt() (int, error) {
	var value int
	var position int
	var currentByte byte

	for {
		buff := make([]byte, 1)
		_, err := pph.reader.Read(buff)
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

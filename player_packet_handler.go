package main

import (
	"github.com/mkorman9/go-minecraft-server/packets"
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
	packetReader *packets.PacketReader
	packetWriter *packets.PacketWriter

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
		packetReader:                packets.NewPacketReader(),
		packetWriter:                packets.NewPacketWriter(connection),
		ip:                          ip,
		enabledCompressionThreshold: -1,
		canceled:                    false,
		canceledMutex:               sync.Mutex{},
	}
}

func (pph *PlayerPacketHandler) ReadLoop() {
	defer func() {
		pph.world.RemovePlayer(pph.player)
		pph.Cancel(nil)
	}()

	for {
		packetDelivery, err := pph.packetReader.Read(pph.reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			if netOpError, ok := err.(*net.OpError); ok {
				if netOpError.Err.Error() == "use of closed network connection" {
					break
				}
			}

			log.Printf("%v\n", err)
			break
		}

		err = pph.HandlePacket(packetDelivery)
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

	pph.reader = cipherStream.WrapReader(pph.connection)
	pph.packetWriter.SetWriter(cipherStream.WrapWriter(pph.connection))

	return nil
}

func (pph *PlayerPacketHandler) setupCompression() error {
	compressionThreshold := pph.world.Settings().CompressionThreshold

	if compressionThreshold >= 0 {
		err := pph.sendSetCompressionRequest(compressionThreshold)
		if err != nil {
			return err
		}

		pph.enabledCompressionThreshold = compressionThreshold

		pph.packetReader.SetCompression(compressionThreshold)
		pph.packetWriter.SetCompression(compressionThreshold)
	}

	return nil
}

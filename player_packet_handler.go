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
		packetReader:                packets.NewPacketReader(connection),
		packetWriter:                packets.NewPacketWriter(connection),
		ip:                          ip,
		enabledCompressionThreshold: -1,
		canceled:                    false,
		canceledMutex:               sync.Mutex{},
	}
}

func (pph *PlayerPacketHandler) ReadLoop() {
	defer func() {
		pph.Cancel(nil)
	}()

	for {
		packetDelivery, err := pph.packetReader.Read()
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

		// discard the remaining (unread) part of data before reading next packet
		_, err = io.Copy(io.Discard, packetDelivery.Reader)
		if err != nil {
			log.Printf("%v\n", err)
			break
		}
	}
}

func (pph *PlayerPacketHandler) Cancel(reason *ChatMessage) {
	pph.canceledMutex.Lock()
	if pph.canceled {
		return
	}
	pph.canceled = true
	pph.canceledMutex.Unlock()

	switch pph.state {
	case PlayerStateBeforeHandshake:
		// nop
	case PlayerStateLogin:
		_ = pph.sendCancelLogin(reason)
	case PlayerStateEncryption:
		_ = pph.sendCancelLogin(reason)
	case PlayerStatePlay:
		_ = pph.sendDisconnect(reason)
		pph.player.OnDisconnect()
	}

	_ = pph.connection.Close()
}

func (pph *PlayerPacketHandler) setupEncryption() error {
	cipherStream, err := packets.NewCipherStream(pph.sharedSecret)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	pph.packetReader.SetEncryption(cipherStream)
	pph.packetWriter.SetEncryption(cipherStream)

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

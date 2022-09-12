package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
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

type HandshakeType = int

const (
	HandshakeTypeStatus = 1
	HandshakeTypeLogin  = 2
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

	ip              string
	playerPublicKey *rsa.PublicKey
	verifyToken     string
	sharedSecret    []byte
	serverHash      string

	canceled      bool
	canceledMutex sync.Mutex
}

func NewPlayerPacketHandler(player *Player, world *World, connection net.Conn, ip string) *PlayerPacketHandler {
	return &PlayerPacketHandler{
		player:        player,
		world:         world,
		connection:    connection,
		state:         PlayerStateBeforeHandshake,
		reader:        connection,
		writer:        connection,
		packetWriter:  NewPacketWriter(),
		ip:            ip,
		canceled:      false,
		canceledMutex: sync.Mutex{},
	}
}

func (pph *PlayerPacketHandler) ReadLoop() {
	defer func() {
		pph.world.PlayerList().UnregisterPlayer(pph.player)
		pph.Cancel(nil)
	}()

	for {
		packetSize, err := pph.readPacketSize()
		if err != nil {
			if err == io.EOF {
				return
			}
			if netOpError, ok := err.(*net.OpError); ok {
				if netOpError.Err.Error() == "use of closed network connection" {
					return
				}
			}

			log.Printf("%v\n", err)
			return
		}

		if packetSize > MaxPacketSize {
			log.Println("invalid packet size")
			return
		}

		packetData := make([]byte, packetSize)
		_, err = pph.reader.Read(packetData)
		if err != nil {
			log.Printf("%v\n", err)
			return
		}

		err = pph.HandlePacket(packetData)
		if err != nil {
			if handlingError, ok := err.(*PacketHandlingError); ok {
				pph.Cancel(handlingError.reason)
			}

			log.Printf("%v\n", err)
			return
		}
	}
}

func (pph *PlayerPacketHandler) HandlePacket(packet []byte) (err error) {
	packetReader := NewPacketReaderContext(packet)
	packetId := packetReader.FetchVarInt()

	switch pph.state {
	case PlayerStateBeforeHandshake:
		err = pph.OnBeforeHandshakePacket(packetId, packetReader)
	case PlayerStateLogin:
		err = pph.OnLoginPacket(packetId, packetReader)
	case PlayerStateEncryption:
		err = pph.OnEncryptionPacket(packetId, packetReader)
	case PlayerStatePlay:
		err = pph.OnPlayPacket(packetId, packetReader)
	}

	return
}

func (pph *PlayerPacketHandler) OnBeforeHandshakePacket(packetId int, packetReader *PacketReaderContext) error {
	switch packetId {
	case 0x00:
		if packetReader.BytesLeft() > 0 {
			return pph.OnHandshakeRequest(packetReader)
		} else {
			return pph.OnStatusRequest(packetReader)
		}
	case 0x01:
		return pph.OnPing(packetReader)
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in before handshake state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnLoginPacket(packetId int, packetReader *PacketReaderContext) error {
	switch packetId {
	case 0x00:
		return pph.OnLoginStartRequest(packetReader)
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in login state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnEncryptionPacket(packetId int, packetReader *PacketReaderContext) error {
	switch packetId {
	case 0x01:
		return pph.OnEncryptionResponse(packetReader)
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in encryption state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnPlayPacket(packetId int, packetReader *PacketReaderContext) error {
	switch packetId {
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in play state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnHandshakeRequest(packetReader *PacketReaderContext) error {
	log.Println("received HandshakeRequest")

	var request HandshakeRequest
	err := request.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	switch request.NextState {
	case HandshakeTypeStatus:
		return pph.sendHandshakeStatusResponse()
	case HandshakeTypeLogin:
		pph.state = PlayerStateLogin
	}

	return nil
}

func (pph *PlayerPacketHandler) OnStatusRequest(_ *PacketReaderContext) error {
	log.Println("received StatusRequest")

	// ignore

	return nil
}

func (pph *PlayerPacketHandler) OnPing(packetReader *PacketReaderContext) error {
	log.Println("received PingRequest")

	var request PingRequest
	err := request.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	return pph.sendPongResponse(request.Payload)
}

func (pph *PlayerPacketHandler) OnLoginStartRequest(packetReader *PacketReaderContext) error {
	log.Println("received LoginStartRequest")

	var request LoginStartRequest
	err := request.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.Name = request.Name
	pph.verifyToken, _ = getSecureRandomString(VerifyTokenLength)

	if request.PublicKey != "" {
		publicKey, err := loadPublicKey(request.PublicKey)
		if err != nil {
			log.Printf("%v\n", err)
			return NewPacketHandlingError(err, NewChatMessage("Malformed Public Key"))
		}

		pph.playerPublicKey = publicKey
	}

	if pph.world.Settings().OnlineMode {
		pph.state = PlayerStateEncryption
		return pph.sendEncryptionRequest()
	} else {
		err := pph.setupCompression()
		if err != nil {
			return err
		}

		err = pph.sendLoginSuccessResponse()
		if err != nil {
			return err
		}

		return pph.OnJoin()
	}
}

func (pph *PlayerPacketHandler) OnEncryptionResponse(packetReader *PacketReaderContext) error {
	log.Println("received EncryptionResponse")

	var response EncryptionResponse
	err := response.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	sharedSecret, err := pph.world.Server().DecryptMessage(response.SharedSecret)
	if err != nil {
		log.Printf("%v\n", err)
		return NewPacketHandlingError(err, NewChatMessage("Malformed Shared Secret"))
	}

	pph.sharedSecret = sharedSecret
	pph.serverHash = pph.world.Server().GenerateServerHash(sharedSecret)

	if response.VerifyToken != nil {
		verifyToken, err := pph.world.Server().DecryptMessage(response.VerifyToken)
		if err != nil {
			log.Printf("%v\n", err)
			return NewPacketHandlingError(err, NewChatMessage("Malformed Verify Token"))
		}

		if string(verifyToken) != pph.verifyToken {
			log.Printf("token mismatch\n")
			return NewPacketHandlingError(err, NewChatMessage("Token mismatch"))
		}
	} else {
		err = verifyRsaSignature(
			pph.playerPublicKey,
			pph.verifyToken,
			response.Salt,
			response.MessageSignature,
		)
		if err != nil {
			log.Printf("%v\n", err)
			return NewPacketHandlingError(err, NewChatMessage("Signature verification error"))
		}
	}

	err = pph.setupEncryption()
	if err != nil {
		return err
	}

	err = pph.setupCompression()
	if err != nil {
		return err
	}

	err = pph.sendLoginSuccessResponse()
	if err != nil {
		return err
	}

	return pph.OnJoin()
}

func (pph *PlayerPacketHandler) Cancel(reason *ChatMessage) {
	pph.canceledMutex.Lock()
	if pph.canceled {
		return
	}
	pph.canceled = true
	pph.canceledMutex.Unlock()

	if reason == nil {
		_ = pph.connection.Close()
		return
	}

	switch pph.state {
	case PlayerStateBeforeHandshake:
		// nop
	case PlayerStateLogin:
		_ = pph.sendCancelLogin(reason)
	case PlayerStateEncryption:
		_ = pph.sendCancelLogin(reason)
	case PlayerStatePlay:
		_ = pph.sendDisconnect(reason)
	}

	_ = pph.connection.Close()
}

func (pph *PlayerPacketHandler) OnJoin() error {
	pph.state = PlayerStatePlay
	pph.world.PlayerList().RegisterPlayer(pph.player)
	pph.player.OnJoin()

	return pph.sendPlayPacket()
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
	}

	return nil
}

func (pph *PlayerPacketHandler) sendHandshakeStatusResponse() error {
	serverStatus := pph.world.GetStatus()
	serverStatusJSON, err := serverStatus.Encode()
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	response := &HandshakeResponse{
		StatusJSON: serverStatusJSON,
	}

	return pph.writePacket(response)
}

func (pph *PlayerPacketHandler) sendPongResponse(payload int64) error {
	packet := &PongResponse{
		Payload: payload,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendEncryptionRequest() error {
	response := &EncryptionRequest{
		ServerID:    "",
		PublicKey:   pph.world.Server().PublicKey(),
		VerifyToken: pph.verifyToken,
	}

	return pph.writePacket(response)
}

func (pph *PlayerPacketHandler) sendSetCompressionRequest(compressionThreshold int) error {
	request := &SetCompressionRequest{
		Threshold: compressionThreshold,
	}

	return pph.writePacket(request)
}

func (pph *PlayerPacketHandler) sendCancelLogin(reason *ChatMessage) error {
	packet := &CancelLoginPacket{
		Reason: reason,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendLoginSuccessResponse() error {
	response := &LoginSuccessResponse{
		UUID:     pph.player.UUID,
		Username: pph.player.Name,
	}

	return pph.writePacket(response)
}

func (pph *PlayerPacketHandler) sendPlayPacket() error {
	packet := &PlayPacket{
		EntityID:            0,
		IsHardcore:          false,
		GameMode:            0,
		PreviousGameMode:    0xff,
		WorldNames:          []string{"minecraft:overworld", "minecraft:the_nether", "minecraft:the_nether"},
		DimensionCodec:      DefaultRegistryCodec(),
		WorldType:           "minecraft:overworld",
		WorldName:           "minecraft:overworld",
		HashedSeed:          1,
		MaxPlayers:          pph.world.Settings().MaxPlayers,
		ViewDistance:        10,
		SimulationDistance:  10,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
		IsDebug:             true,
		IsFlat:              false,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendDisconnect(reason *ChatMessage) error {
	response := &DisconnectPacket{
		Reason: reason,
	}

	return pph.writePacket(response)
}

func (pph *PlayerPacketHandler) writePacket(packet Packet) error {
	data, err := packet.Marshal(pph.packetWriter.New())
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	_, err = pph.writer.Write(data)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	return nil
}

func (pph *PlayerPacketHandler) readPacketSize() (int, error) {
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

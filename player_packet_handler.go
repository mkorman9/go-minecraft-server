package main

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
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
		_ = pph.connection.Close()
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
	case 0x00:
		return pph.OnTeleportConfirm(packetReader)
	case 0x07:
		return pph.OnSettings(packetReader)
	case 0x0c:
		return pph.OnCustomPayload(packetReader)
	case 0x13:
		return pph.OnPosition(packetReader)
	case 0x14:
		return pph.OnPositionLook(packetReader)
	case 0x15:
		return pph.OnLook(packetReader)
	case 0x2e:
		return pph.OnArmAnimation(packetReader)
	case 0x1b:
		return pph.OnAbilities(packetReader)
	case 0x2a:
		return pph.OnSetCreativeSlot(packetReader)
	case 0x03:
		return pph.OnChatCommand(packetReader)
	case 0x04:
		return pph.OnChatMessage(packetReader)
	default:
		log.Printf("unrecognized packet id: 0x%x in play state\n", packetId)
		return nil
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
	pph.player.Signature = request.Signature
	pph.verifyToken, _ = getSecureRandomString(VerifyTokenLength)

	if request.PublicKey != nil {
		publicKey, err := loadPublicKey(request.PublicKey)
		if err != nil {
			log.Printf("%v\n", err)
			return NewPacketHandlingError(err, NewChatMessage("Malformed Public Key"))
		}

		pph.player.PublicKey = publicKey
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
			pph.player.PublicKey,
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

func (pph *PlayerPacketHandler) OnTeleportConfirm(packetReader *PacketReaderContext) error {
	log.Println("received TeleportConfirm")

	var packet TeleportConfirmPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	// nop

	return nil
}

func (pph *PlayerPacketHandler) OnSettings(packetReader *PacketReaderContext) error {
	log.Println("received Settings")

	var packet SettingsPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnClientSettings(&PlayerClientSettings{
		Locale:              packet.Locale,
		ViewDistance:        packet.ViewDistance,
		ChatColors:          packet.ChatColors,
		SkinParts:           packet.SkinParts,
		MainHand:            packet.MainHand,
		EnableTextFiltering: packet.EnableTextFiltering,
		EnableServerListing: packet.EnableServerListing,
	})

	return nil
}

func (pph *PlayerPacketHandler) OnPosition(packetReader *PacketReaderContext) error {
	var packet PositionPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnPositionUpdate(packet.X, packet.Y, packet.Z)
	pph.player.OnGroundUpdate(packet.OnGround)

	return nil
}

func (pph *PlayerPacketHandler) OnPositionLook(packetReader *PacketReaderContext) error {
	var packet PositionLookPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnPositionUpdate(packet.X, packet.Y, packet.Z)
	pph.player.OnGroundUpdate(packet.OnGround)
	pph.player.OnLookUpdate(packet.Yaw, packet.Pitch)

	return nil
}

func (pph *PlayerPacketHandler) OnLook(packetReader *PacketReaderContext) error {
	var packet LookPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnGroundUpdate(packet.OnGround)
	pph.player.OnLookUpdate(packet.Yaw, packet.Pitch)

	return nil
}

func (pph *PlayerPacketHandler) OnCustomPayload(packetReader *PacketReaderContext) error {
	log.Println("received CustomPayload")

	var packet CustomPayloadPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnPluginChannel(packet.Channel, packet.Data)

	return nil
}

func (pph *PlayerPacketHandler) OnArmAnimation(packetReader *PacketReaderContext) error {
	log.Println("received ArmAnimation")

	var packet ArmAnimationPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnArmAnimation(packet.Hand)

	return nil
}

func (pph *PlayerPacketHandler) OnAbilities(packetReader *PacketReaderContext) error {
	log.Println("received Abilities")

	var packet AbilitiesPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	// TODO

	return nil
}

func (pph *PlayerPacketHandler) OnSetCreativeSlot(packetReader *PacketReaderContext) error {
	log.Println("received SetCreativeSlot")

	var packet SetCreativeSlotPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	// TODO

	return nil
}

func (pph *PlayerPacketHandler) OnChatCommand(packetReader *PacketReaderContext) error {
	log.Println("received ChatCommand")

	var packet ChatCommandPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnChatCommand(packet.Message, time.UnixMilli(packet.Timestamp))

	return nil
}

func (pph *PlayerPacketHandler) OnChatMessage(packetReader *PacketReaderContext) error {
	log.Println("received ChatMessage")

	var packet ChatMessagePacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnChatMessage(packet.Message, time.UnixMilli(packet.Timestamp))

	return nil
}

func (pph *PlayerPacketHandler) Cancel(reason *ChatMessage) {
	pph.canceledMutex.Lock()
	if pph.canceled {
		return
	}
	pph.canceled = true
	pph.canceledMutex.Unlock()

	if reason == nil {
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
}

func (pph *PlayerPacketHandler) OnJoin() error {
	err := pph.sendPlayPacket()
	if err != nil {
		return err
	}

	err = pph.sendSpawnPosition()
	if err != nil {
		return err
	}

	pph.state = PlayerStatePlay
	pph.player.OnJoin(GameModeSurvival)

	return nil
}

func (pph *PlayerPacketHandler) SendSystemChatMessage(message *ChatMessage) error {
	return pph.sendSystemChatMessage(message)
}

func (pph *PlayerPacketHandler) SynchronizePosition(x float64, y float64, z float64) error {
	return pph.sendPositionUpdate(x, y, z)
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
		GameMode:            GameModeSurvival,
		PreviousGameMode:    GameModeUnknown,
		WorldNames:          []string{"minecraft:overworld", "minecraft:the_nether", "minecraft:the_nether"},
		DimensionCodec:      *pph.world.Data().DimensionCodec,
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
	packet := &DisconnectPacket{
		Reason: reason,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendSystemChatMessage(message *ChatMessage) error {
	packet := &SystemChatPacket{
		Content: message,
		Type:    SystemChatMessageTypeChat,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendSpawnPosition() error {
	packet := &SpawnPositionPacket{
		Location: pph.world.Data().SpawnPosition,
		Angle:    0,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendPositionUpdate(x float64, y float64, z float64) error {
	packet := &UpdatePositionPacket{
		X:               x,
		Y:               y,
		Z:               z,
		Yaw:             pph.player.Yaw,
		Pitch:           pph.player.Pitch,
		Flags:           0,
		TeleportID:      0,
		DismountVehicle: false,
	}

	return pph.writePacket(packet)
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

package main

import (
	"crypto/rsa"
	"io"
	"log"
	"net"
	"strings"
)

type Player struct {
	world        *World
	connection   net.Conn
	reader       io.Reader
	writer       io.Writer
	packetWriter *PacketWriter

	name         string
	uuid         UUID
	ip           string
	state        PlayerState
	publicKey    *rsa.PublicKey
	verifyToken  string
	sharedSecret string
	serverHash   string
}

type PlayerState = int

const (
	PlayerStateBeforeHandshake = iota
	PlayerStateLogin           = iota
	PlayerStateEncryption      = iota
	PlayerStatePlay            = iota
)

func NewPlayer(world *World, connection net.Conn) *Player {
	ip, _, _ := strings.Cut(connection.RemoteAddr().String(), ":")

	return &Player{
		world:        world,
		connection:   connection,
		reader:       connection,
		writer:       connection,
		packetWriter: NewPacketWriter(),
		uuid:         getRandomUUID(),
		ip:           ip,
		state:        PlayerStateBeforeHandshake,
	}
}

func (p *Player) Disconnect() {
	_ = p.connection.Close()
}

func (p *Player) IsOnline() bool {
	return p.state == PlayerStatePlay
}

func (p *Player) HandlePacket(data []byte) {
	reader := NewPacketReaderContext(data)
	packetId := reader.FetchVarInt()

	switch p.state {
	case PlayerStateBeforeHandshake:
		switch packetId {
		case 0x00:
			if reader.BytesLeft() == 0 {
				p.OnStatusRequest()
			} else {
				p.OnHandshakeRequest(UnmarshalPacket(reader, &HandshakeRequest{}))
			}
		case 0x01:
			p.OnPing(UnmarshalPacket(reader, &PingRequest{}))
		}
	case PlayerStateLogin:
		switch packetId {
		case 0x00:
			p.OnLoginStartRequest(UnmarshalPacket(reader, &LoginStartRequest{}))
		}
	case PlayerStateEncryption:
		switch packetId {
		case 0x01:
			p.OnEncryptionResponse(UnmarshalPacket(reader, &EncryptionResponse{}))
		}
	case PlayerStatePlay:
		switch packetId {
		}
	}
}

func (p *Player) OnHandshakeRequest(request *HandshakeRequest, err error) {
	if err != nil {
		return // ignore request
	}

	log.Println("received HandshakeRequest")

	switch request.NextState {
	case HandshakeStateStatus:
		p.sendHandshakeResponse()
	case HandshakeStateLogin:
		p.state = PlayerStateLogin
	}
}

func (p *Player) OnStatusRequest() {
	log.Println("received StatusRequest")
}

func (p *Player) OnPing(request *PingRequest, err error) {
	if err != nil {
		return // ignore request
	}

	log.Println("received Ping")

	p.sendPongResponse(request.Payload)
}

func (p *Player) OnLoginStartRequest(request *LoginStartRequest, err error) {
	if err != nil {
		p.sendCancelLogin(NewChatMessage("Invalid Login Request"))
		p.Disconnect()
		return
	}

	log.Println("received LoginStartRequest")

	p.name = request.Name
	p.verifyToken, _ = getSecureRandomString(VerifyTokenLength)

	if request.PublicKey != "" {
		publicKey, err := loadPublicKey(request.PublicKey)
		if err != nil {
			log.Printf("%v\n", err)
			p.sendCancelLogin(NewChatMessage("Malformed Public Key"))
			p.Disconnect()
			return
		}

		p.publicKey = publicKey
	}

	if p.world.settings.OnlineMode {
		p.state = PlayerStateEncryption
		p.sendEncryptionRequest()
	} else {
		p.state = PlayerStatePlay
		p.setupCompression()
		p.sendLoginSuccessResponse()

		p.OnPlayStart()
	}
}

func (p *Player) OnEncryptionResponse(response *EncryptionResponse, err error) {
	if err != nil {
		p.sendCancelLogin(NewChatMessage("Invalid Login Request"))
		p.Disconnect()
		return
	}

	log.Println("received EncryptionResponse")

	sharedSecret, err := p.world.DecryptServerMessage(response.SharedSecret)
	if err != nil {
		log.Printf("%v\n", err)
		p.sendCancelLogin(NewChatMessage("Malformed Shared Secret"))
		p.Disconnect()
		return
	}

	p.sharedSecret = sharedSecret
	p.serverHash = p.world.GenerateServerHash(sharedSecret)

	if response.VerifyToken != "" {
		verifyToken, err := p.world.DecryptServerMessage(response.VerifyToken)
		if err != nil {
			log.Printf("%v\n", err)
			p.sendCancelLogin(NewChatMessage("Malformed Verify Token"))
			p.Disconnect()
			return
		}

		if verifyToken != p.verifyToken {
			log.Printf("token mismatch\n")
			p.sendCancelLogin(NewChatMessage("Token mismatch"))
			p.Disconnect()
			return
		}
	} else {
		err = verifyRsaSignature(
			p.publicKey,
			p.verifyToken,
			response.Salt,
			response.MessageSignature,
		)
		if err != nil {
			log.Printf("%v\n", err)
			p.sendCancelLogin(NewChatMessage("Signature verification error"))
			p.Disconnect()
			return
		}
	}

	// Verify user info here
	//fmt.Printf(
	//	"https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s\n",
	//	p.name,
	//	p.serverHash,
	//)

	p.state = PlayerStatePlay
	p.setupEncryption()
	p.setupCompression()
	p.sendLoginSuccessResponse()

	p.OnPlayStart()
}

func (p *Player) OnPlayStart() {
	p.sendPlayPacket()
}

func (p *Player) setupEncryption() {
	cipherStream, err := NewCipherStream(p.sharedSecret)
	if err != nil {
		log.Printf("%v\n", err)
		p.Disconnect()
		return
	}

	p.reader = cipherStream.WrapReader(p.reader)
	p.writer = cipherStream.WrapWriter(p.writer)
}

func (p *Player) setupCompression() {
	if p.world.settings.CompressionThreshold >= 0 {
		p.sendSetCompressionRequest()
		p.packetWriter.EnableCompression(p.world.settings.CompressionThreshold)
	}
}

func (p *Player) sendHandshakeResponse() {
	serverStatus := p.world.GetServerStatus()
	serverStatusJSON, err := serverStatus.Encode()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	response := &HandshakeResponse{
		StatusJSON: serverStatusJSON,
	}

	p.writePacket(response)
}

func (p *Player) sendPongResponse(payload int64) {
	packet := &PongResponse{
		Payload: payload,
	}

	p.writePacket(packet)
}

func (p *Player) sendCancelLogin(reason *ChatMessage) {
	packet := &CancelLoginPacket{
		Reason: reason,
	}

	p.writePacket(packet)
}

func (p *Player) sendEncryptionRequest() {
	serverKey := p.world.serverKey.publicASN1

	response := &EncryptionRequest{
		ServerID:    "",
		PublicKey:   serverKey,
		VerifyToken: p.verifyToken,
	}

	p.writePacket(response)
}

func (p *Player) sendSetCompressionRequest() {
	request := &SetCompressionRequest{
		Threshold: p.world.settings.CompressionThreshold,
	}

	p.writePacket(request)
}

func (p *Player) sendLoginSuccessResponse() {
	response := &LoginSuccessResponse{
		UUID:     p.uuid,
		Username: p.name,
	}

	p.writePacket(response)
}

func (p *Player) sendDisconnect(reason *ChatMessage) {
	response := &DisconnectPacket{
		Reason: reason,
	}

	p.writePacket(response)
	p.Disconnect()
}

func (p *Player) sendPlayPacket() {
	packet := &PlayPacket{
		EntityID:            0,
		IsHardcore:          false,
		GameMode:            0,
		PreviousGameMode:    0xff,
		WorldNames:          []string{"world"},
		RegistryCodec:       DefaultRegistryCodec(),
		WorldType:           "world",
		WorldName:           "world",
		HashedSeed:          1,
		MaxPlayers:          p.world.settings.MaxPlayers,
		ViewDistance:        10,
		SimulationDistance:  10,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
		IsDebug:             true,
		IsFlat:              false,
	}

	p.writePacket(packet)
}

func (p *Player) writePacket(packet Packet) {
	data, err := packet.Marshal(p.packetWriter.New())
	if err != nil {
		log.Printf("%v\n", err)
		p.Disconnect()
		return
	}

	_, err = p.writer.Write(data)
	if err != nil {
		log.Printf("%v\n", err)
		p.Disconnect()
		return
	}
}

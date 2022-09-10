package main

import (
	"crypto/rsa"
	"log"
	"net"
	"strings"
)

type Player struct {
	world      *World
	connection net.Conn
	ip         string

	name         string
	uuid         UUID
	state        PlayerState
	publicKey    *rsa.PublicKey
	verifyToken  string
	sharedSecret string
	serverHash   string
	cipherStream *CipherStream
}

type PlayerState = int

const (
	PlayerStateBeforeHandshake = iota
	PlayerStateAfterHandshake  = iota
	PlayerStateAfterLoginStart = iota
	PlayerStatePlay            = iota
)

func NewPlayer(world *World, connection net.Conn) *Player {
	ip, _, _ := strings.Cut(connection.RemoteAddr().String(), ":")

	return &Player{
		world:      world,
		connection: connection,
		ip:         ip,
		uuid:       getRandomUUID(),
		state:      PlayerStateBeforeHandshake,
	}
}

func (p *Player) Disconnect() {
	_ = p.connection.Close()
}

func (p *Player) IsLoggedIn() bool {
	return p.state >= PlayerStatePlay
}

func (p *Player) HandlePacket(data []byte) {
	reader := &PacketReader{data: data, cursor: 0}
	packetId := reader.FetchVarInt()

	switch p.state {
	case PlayerStateBeforeHandshake:
		switch packetId {
		case 0x00:
			if reader.BytesLeft() == 0 {
				p.OnStatusRequest()
			} else {
				p.OnHandshakeRequest(ReadHandshakeRequest(reader))
			}
		case 0x01:
			p.OnPing(ReadPingRequest(reader))
		}
	case PlayerStateAfterHandshake:
		switch packetId {
		case 0x00:
			p.OnLoginStartRequest(ReadLoginStartRequest(reader))
		}
	case PlayerStateAfterLoginStart:
		switch packetId {
		case 0x01:
			p.OnEncryptionResponse(ReadEncryptionResponse(reader))
		}
	case PlayerStatePlay:
		switch packetId {
		}
	}
}

func (p *Player) OnHandshakeRequest(request *HandshakeRequest) {
	log.Println("received HandshakeRequest")

	switch request.NextState {
	case HandshakeStateStatus:
		p.SendHandshakeResponse()
	case HandshakeStateLogin:
		p.state = PlayerStateAfterHandshake
	}
}

func (p *Player) OnStatusRequest() {
	log.Println("received StatusRequest")
}

func (p *Player) OnPing(request *PingRequest) {
	log.Println("received Ping")

	response := &PongResponse{
		Payload: request.Payload,
	}

	_, _ = p.connection.Write(response.Bytes())
}

func (p *Player) OnLoginStartRequest(request *LoginStartRequest) {
	log.Println("received LoginStartRequest")

	p.name = request.Name
	p.verifyToken, _ = getSecureRandomString(VerifyTokenLength)

	if request.PublicKey != "" {
		publicKey, err := loadPublicKey(request.PublicKey)
		if err != nil {
			log.Printf("%v\n", err)
			p.Disconnect()
			return
		}

		p.publicKey = publicKey
	}

	if p.world.settings.OnlineMode {
		p.state = PlayerStateAfterLoginStart
		p.SendEncryptionRequest()
	} else {
		p.state = PlayerStatePlay
		p.SendLoginSuccessResponse()
	}
}

func (p *Player) OnEncryptionResponse(response *EncryptionResponse) {
	log.Println("received EncryptionResponse")

	sharedSecret, err := p.world.DecryptServerMessage(response.SharedSecret)
	if err != nil {
		log.Printf("%v\n", err)
		p.Disconnect()
		return
	}

	p.sharedSecret = sharedSecret
	p.serverHash = p.world.GenerateServerHash(sharedSecret)

	if response.VerifyToken != "" {
		verifyToken, err := p.world.DecryptServerMessage(response.VerifyToken)
		if err != nil {
			log.Printf("%v\n", err)
			p.Disconnect()
			return
		}

		if verifyToken != p.verifyToken {
			log.Printf("token mismatch\n")
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

	cipherStream, err := NewCipherStream(p.sharedSecret)
	if err != nil {
		log.Printf("%v\n", err)
		p.Disconnect()
		return
	}

	p.cipherStream = cipherStream
	p.state = PlayerStatePlay

	p.SendLoginSuccessResponse()
}

func (p *Player) OnPlayStart() {
}

func (p *Player) SendHandshakeResponse() {
	serverStatus := p.world.GetServerStatus()
	serverStatusJSON, err := serverStatus.Encode()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	response := &HandshakeResponse{
		StatusJSON: serverStatusJSON,
	}

	_, _ = p.connection.Write(response.Bytes())
}

func (p *Player) SendEncryptionRequest() {
	serverKey := p.world.serverKey.publicASN1

	response := &EncryptionRequest{
		ServerID:    "",
		PublicKey:   serverKey,
		VerifyToken: p.verifyToken,
	}

	_, _ = p.connection.Write(response.Bytes())
}

func (p *Player) SendLoginSuccessResponse() {
	response := &LoginSuccessResponse{
		UUID:     p.uuid,
		Username: p.name,
	}

	_, _ = p.connection.Write(response.Bytes())

	p.OnPlayStart()
}

func (p *Player) SendDisconnect(reason *ChatMessage) {
	response := &DisconnectPacket{
		Reason: reason,
	}

	_, _ = p.connection.Write(response.Bytes())

	p.Disconnect()
}

package main

import (
	"fmt"
	"log"
	"net"
)

type Player struct {
	world      *World
	connection net.Conn

	name        string
	state       PlayerState
	loggedIn    bool
	publicKey   string
	signature   string
	verifyToken string
}

type PlayerState = int

const (
	PlayerStateBeforeHandshake        = iota
	PlayerStateAfterHandshake         = iota
	PlayerStateAfterEncryptionRequest = iota
)

func NewPlayer(world *World, connection net.Conn) *Player {
	return &Player{
		world:      world,
		connection: connection,
		state:      PlayerStateBeforeHandshake,
		loggedIn:   false,
	}
}

func (p *Player) Disconnect() {
	_ = p.connection.Close()
}

func (p *Player) LogIn() {
	p.loggedIn = true
}

func (p *Player) IsLoggedIn() bool {
	return p.loggedIn
}

func (p *Player) HandlePacket(data []byte) {
	reader := &PacketReader{data: data, cursor: 0}
	packetId := reader.FetchVarInt()

	fmt.Printf("Got packet %d\n", packetId)

	switch packetId {
	case 0x00:
		if reader.BytesLeft() > 0 {
			if p.state < PlayerStateAfterHandshake {
				p.OnHandshakeRequest(ReadHandshakeRequest(reader))
			} else {
				p.OnLoginStartRequest(ReadLoginStartRequest(reader))
			}
		} else {
			p.OnStatusRequest()
		}
	case 0x01:
		p.OnPing(ReadPingRequest(reader))
	}
}

func (p *Player) OnHandshakeRequest(request *HandshakeRequest) {
	switch request.NextState {
	case HandshakeStateStatus:
		p.SendHandshakeResponse()
	case HandshakeStateLogin:
		p.state = PlayerStateAfterHandshake
	}
}

func (p *Player) OnStatusRequest() {
}

func (p *Player) OnPing(request *PingRequest) {
	response := &PongResponse{
		Payload: request.Payload,
	}

	_, _ = p.connection.Write(response.Bytes())
}

func (p *Player) OnLoginStartRequest(request *LoginStartRequest) {
	p.name = request.Name
	p.publicKey = request.PublicKey
	p.signature = request.Signature
	p.verifyToken, _ = getSecureRandomString(VerifyTokenLength)

	p.SendEncryptionRequest()

	p.state = PlayerStateAfterEncryptionRequest
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

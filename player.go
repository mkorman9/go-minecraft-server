package main

import (
	"log"
	"net"
)

type Player struct {
	world      *World
	connection net.Conn

	loggedIn bool
}

func NewPlayer(world *World, connection net.Conn) *Player {
	return &Player{
		world:      world,
		connection: connection,
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

func (p *Player) OnHandshakeRequest(request *HandshakeRequest) {
	switch request.NextState {
	case HandshakeStateStatus:
		p.SendHandshakeResponse()
	case HandshakeStateLogin:
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

package main

import (
	"log"
	"net"
)

type World struct {
	settings       *Settings
	serverListener net.Listener
	serverKey      *ServerKey
	playerList     *PlayerList
}

func NewWorld(settings *Settings, serverListener net.Listener, serverKey *ServerKey) *World {
	return &World{
		settings:       settings,
		serverListener: serverListener,
		serverKey:      serverKey,
		playerList:     NewPlayerList(),
	}
}

func (w *World) Shutdown() {
	_ = w.serverListener.Close()
}

func (w *World) RegisterPlayer(player *Player) {
	log.Println("player connected")

	w.playerList.RegisterPlayer(player)
}

func (w *World) UnregisterPlayer(player *Player) {
	log.Println("player disconnected")

	w.playerList.UnregisterPlayer(player)
}

func (w *World) GetPlayersCount() int {
	return w.playerList.GetLoggedInSize()
}

func (w *World) GetServerStatus() *ServerStatus {
	return &ServerStatus{
		Version: ServerStatusVersion{
			Name:     ProtocolName,
			Protocol: ProtocolVersion,
		},
		Players: ServerStatusPlayers{
			Max:    w.settings.MaxPlayers,
			Online: w.GetPlayersCount(),
			Sample: nil,
		},
		Description: ServerStatusDescription{
			Text: w.settings.Description,
		},
		PreviewsChat:       true,
		EnforcesSecureChat: true,
	}
}

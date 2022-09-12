package main

import (
	"net"
)

type World struct {
	settings       *Settings
	data           *Data
	server         *Server
	playerList     *PlayerList
	serverListener net.Listener
}

func NewWorld(settings *Settings) (*World, error) {
	data, err := LoadData()
	if err != nil {
		return nil, err
	}

	server, err := NewServer(settings)
	if err != nil {
		return nil, err
	}

	return &World{
		data:       data,
		server:     server,
		settings:   settings,
		playerList: NewPlayerList(),
	}, nil
}

func (w *World) Settings() *Settings {
	return w.settings
}

func (w *World) Data() *Data {
	return w.data
}

func (w *World) Server() *Server {
	return w.server
}

func (w *World) PlayerList() *PlayerList {
	return w.playerList
}

func (w *World) GetStatus() *ServerStatus {
	return &ServerStatus{
		Version: ServerStatusVersion{
			Name:     ProtocolName,
			Protocol: ProtocolVersion,
		},
		Players: ServerStatusPlayers{
			Max:    w.settings.MaxPlayers,
			Online: w.PlayerList().Len(),
			Sample: nil,
		},
		Description: ChatMessage{
			Text: w.settings.Description,
		},
		PreviewsChat:       true,
		EnforcesSecureChat: true,
	}
}

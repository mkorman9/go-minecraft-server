package main

import (
	"math/rand"
	"net"
	"time"
)

type World struct {
	settings       *Settings
	data           *Data
	server         *Server
	playerList     *PlayerList
	backgroundJob  *BackgroundJob
	entityStore    *EntityStore
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

	world := &World{
		data:        data,
		server:      server,
		settings:    settings,
		playerList:  NewPlayerList(),
		entityStore: NewEntityStore(),
	}

	world.backgroundJob = NewBackgroundJob(world)
	world.backgroundJob.Start()

	return world, nil
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

func (w *World) JoinPlayer(player *Player) {
	w.PlayerList().RegisterPlayer(player)
}

func (w *World) RemovePlayer(player *Player) {
	w.PlayerList().UnregisterPlayer(player)
	w.entityStore.RemoveID(player.EntityID)
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

func (w *World) BroadcastKeepAlive() {
	w.PlayerList().All(func(p *Player) {
		keepAliveID := rand.Int63()
		p.SendKeepAlive(keepAliveID)
	})
}

func (w *World) KickUnresponsivePlayers() {
	timeout := w.Settings().PlayerTimeout * time.Second

	w.PlayerList().All(func(p *Player) {
		timeSinceLastHeartbeat := time.Now().Sub(p.lastHeartbeat)
		if timeSinceLastHeartbeat > timeout {
			p.Kick(NewChatMessage("Timed out"))
		}
	})
}

func (w *World) BroadcastPlayerJoined(player *Player) {
	w.PlayerList().All(func(p *Player) {
		p.SendAnotherPlayerJoined(player)
	})
}

func (w *World) BroadcastPlayerDisconnected(player *Player) {
	w.PlayerList().All(func(p *Player) {
		p.SendAnotherPlayerDisconnected(player)
	})
}

func (w *World) GenerateEntityID() int32 {
	return w.entityStore.GenerateID()
}

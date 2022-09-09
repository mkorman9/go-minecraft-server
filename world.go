package main

import "log"

type World struct {
	playerList *PlayerList
}

func NewWorld() *World {
	return &World{
		playerList: NewPlayerList(),
	}
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
			Name:     "1.19",
			Protocol: 759,
		},
		Players: ServerStatusPlayers{
			Max:    2137,
			Online: w.GetPlayersCount(),
			Sample: nil,
		},
		Description: ServerStatusDescription{
			Text: "Simple Go Server",
		},
		PreviewsChat:       true,
		EnforcesSecureChat: true,
	}
}

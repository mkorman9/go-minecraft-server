package main

import "sync"

type PlayerList struct {
	m    sync.RWMutex
	list []*Player
}

func NewPlayerList() *PlayerList {
	return &PlayerList{
		m: sync.RWMutex{},
	}
}

func (pl *PlayerList) RegisterPlayer(player *Player) {
	pl.m.Lock()
	defer pl.m.Unlock()

	pl.list = append(pl.list, player)
}

func (pl *PlayerList) UnregisterPlayer(player *Player) {
	pl.m.Lock()
	defer pl.m.Unlock()

	index := -1
	for i, p := range pl.list {
		if player == p {
			index = i
			break
		}
	}

	if index != -1 {
		pl.list = append(pl.list[:index], pl.list[index+1:]...)
	}
}

func (pl *PlayerList) GetOnlineSize() int {
	pl.m.RLock()
	defer pl.m.RUnlock()

	count := 0
	for _, player := range pl.list {
		if player.IsOnline() {
			count++
		}
	}

	return count
}

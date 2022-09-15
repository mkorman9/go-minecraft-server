package main

import (
	"strings"
	"sync"
)

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

func (pl *PlayerList) Len() int {
	pl.m.RLock()
	defer pl.m.RUnlock()

	return len(pl.list)
}

func (pl *PlayerList) All(handler func(*Player)) {
	pl.m.RLock()
	defer pl.m.RUnlock()

	for _, player := range pl.list {
		handler(player)
	}
}

func (pl *PlayerList) ByName(name string, handler func(*Player)) (found bool) {
	pl.m.RLock()
	defer pl.m.RUnlock()

	for _, player := range pl.list {
		if strings.EqualFold(player.Name, name) {
			found = true
			handler(player)
			break
		}
	}

	return
}

func (pl *PlayerList) Copy() []*Player {
	pl.m.RLock()
	defer pl.m.RUnlock()

	return pl.list[:]
}

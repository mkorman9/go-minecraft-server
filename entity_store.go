package main

import (
	"math/rand"
	"sync"
)

type EntityStore struct {
	store map[int32]struct{}
	m     sync.Mutex
}

func NewEntityStore() *EntityStore {
	return &EntityStore{
		store: make(map[int32]struct{}),
	}
}

func (es *EntityStore) GenerateID() int32 {
	es.m.Lock()
	defer es.m.Unlock()

	for {
		value := rand.Int31()
		if _, ok := es.store[value]; !ok {
			es.store[value] = struct{}{}
			return value
		}
	}
}

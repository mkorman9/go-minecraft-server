package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
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
	return w.playerList.GetOnlineSize()
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

func (w *World) DecryptServerMessage(message string) (string, error) {
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, w.serverKey.private, []byte(message))
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func (w *World) GenerateServerHash(sharedSecret string) string {
	hash := sha1.New()
	hash.Write([]byte(sharedSecret))
	hash.Write([]byte(w.serverKey.publicASN1))
	return hex.EncodeToString(hash.Sum(nil))
}

package main

import "time"

type Settings struct {
	ServerAddress         string        `json:"serverAddress"`
	Description           string        `json:"description"`
	MaxPlayers            int           `json:"maxPlayers"`
	OnlineMode            bool          `json:"onlineMode"`
	CompressionThreshold  int           `json:"compressionThreshold"`
	IsDebug               bool          `json:"isDebug"`
	ViewDistance          int           `json:"viewDistance"`
	SimulationDistance    int           `json:"simulationDistance"`
	KeepAliveSendInterval time.Duration `json:"keepAliveSendInterval"`
	PlayerTimeout         time.Duration `json:"playerTimeout"`
}

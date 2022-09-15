package main

type Settings struct {
	ServerAddress        string `json:"serverAddress"`
	Description          string `json:"description"`
	MaxPlayers           int    `json:"maxPlayers"`
	OnlineMode           bool   `json:"onlineMode"`
	CompressionThreshold int    `json:"compressionThreshold"`
	IsDebug              bool   `json:"isDebug"`
	ViewDistance         int    `json:"viewDistance"`
	SimulationDistance   int    `json:"simulationDistance"`
}

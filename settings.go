package main

type Settings struct {
	ServerAddress        string `json:"serverAddress"`
	Description          string `json:"description"`
	MaxPlayers           int    `json:"maxPlayers"`
	OnlineMode           bool   `json:"onlineMode"`
	CompressionThreshold int    `json:"compressionThreshold"`
}

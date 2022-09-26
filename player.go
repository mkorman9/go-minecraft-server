package main

import (
	"crypto/rsa"
	"github.com/mkorman9/go-minecraft-server/packets"
	"log"
	"time"
)

type Player struct {
	Name              string
	DisplayName       *ChatMessage
	UUID              packets.UUID
	EntityID          int32
	IP                string
	PublicKey         *rsa.PublicKey
	PublicKeyDER      []byte
	Signature         string
	Timestamp         int64
	ClientSettings    *PlayerClientSettings
	X                 float64
	Y                 float64
	Z                 float64
	Yaw               float32
	Pitch             float32
	OnGround          bool
	GameMode          GameMode
	Textures          string
	TexturesSignature string
	Ping              int

	packetHandler     *PlayerPacketHandler
	world             *World
	lastKeepAliveID   int64
	lastHeartbeat     time.Time
	lastHeartbeatSent time.Time
}

type PlayerClientSettings struct {
	Locale              string
	ViewDistance        byte
	ChatFlags           int
	ChatColors          bool
	SkinParts           byte
	MainHand            int
	EnableTextFiltering bool
	EnableServerListing bool
}

func NewPlayer(world *World, ip string) *Player {
	return &Player{
		Name:        "",
		DisplayName: NewChatMessage(""),
		UUID:        getRandomUUID(),
		EntityID:    -1,
		IP:          ip,
		GameMode:    GameModeUnknown,
		world:       world,
	}
}

func (p *Player) Kick(reason *ChatMessage) {
	p.packetHandler.Cancel(reason)
}

func (p *Player) AssignPacketHandler(packetHandler *PlayerPacketHandler) {
	p.packetHandler = packetHandler
}

func (p *Player) SendSystemChatMessage(message *ChatMessage) {
	err := p.packetHandler.SendSystemChatMessage(message)
	if err != nil {
		log.Printf("Failed to send system chat message: %v\n", err)
	}
}

func (p *Player) SetPosition(x, y, z float64) {
	p.X = x
	p.Y = y
	p.Z = z

	_ = p.packetHandler.SynchronizePosition(x, y, z)
}

func (p *Player) SendKeepAlive(keepAliveID int64) {
	p.lastKeepAliveID = keepAliveID
	p.lastHeartbeatSent = time.Now()
	_ = p.packetHandler.SendKeepAlive(keepAliveID)
}

func (p *Player) SendAnotherPlayerJoined(player *Player) {
	_ = p.packetHandler.sendPlayersAdded([]*Player{player})
}

func (p *Player) SendAnotherPlayerDisconnected(player *Player) {
	_ = p.packetHandler.sendPlayersRemoved([]*Player{player})
}

func (p *Player) OnJoin(gameMode GameMode) {
	p.world.BroadcastPlayerJoined(p)

	p.world.JoinPlayer(p)
	p.GameMode = gameMode
	p.lastHeartbeat = time.Now()
}

func (p *Player) OnDisconnect() {
	p.world.RemovePlayer(p)

	p.world.BroadcastPlayerDisconnected(p)
}

func (p *Player) OnClientSettings(clientSettings *PlayerClientSettings) {
	p.ClientSettings = clientSettings
}

func (p *Player) OnPositionUpdate(x float64, y float64, z float64) {
	p.X = x
	p.Y = y
	p.Z = z
}

func (p *Player) OnLookUpdate(yaw float32, pitch float32) {
	p.Yaw = yaw
	p.Pitch = pitch
}

func (p *Player) OnGroundUpdate(onGround bool) {
	p.OnGround = onGround
}

func (p *Player) OnPluginChannel(channel string, data []byte) {

}

func (p *Player) OnArmAnimation(hand int) {

}

func (p *Player) OnCloseWindow(windowId byte) {

}

func (p *Player) OnChatCommand(command string, timestamp time.Time) {

}

func (p *Player) OnChatMessage(message string, timestamp time.Time) {

}

func (p *Player) OnKeepAliveResponse(keepAliveID int64) {
	if keepAliveID == p.lastKeepAliveID {
		p.lastHeartbeat = time.Now()

		ping := time.Now().Sub(p.lastHeartbeatSent)
		p.Ping = int(ping / time.Millisecond)
	}
}

func (p *Player) OnAction(entityID int, actionID EntityAction, jumpBoost int) {

}

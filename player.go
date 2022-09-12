package main

type Player struct {
	Name           string
	UUID           UUID
	IP             string
	ClientSettings *PlayerClientSettings
	X              float64
	Y              float64
	Z              float64
	Yaw            float32
	Pitch          float32
	OnGround       bool

	packetHandler *PlayerPacketHandler
	world         *World
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
		Name:  "",
		UUID:  getRandomUUID(),
		IP:    ip,
		world: world,
	}
}

func (p *Player) AssignPacketHandler(packetHandler *PlayerPacketHandler) {
	p.packetHandler = packetHandler
}

func (p *Player) OnJoin() {
	p.world.PlayerList().RegisterPlayer(p)
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

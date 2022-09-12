package main

type Player struct {
	Name string
	UUID UUID
	IP   string

	packetHandler *PlayerPacketHandler
	world         *World
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

}

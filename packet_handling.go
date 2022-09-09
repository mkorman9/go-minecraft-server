package main

func HandlePacket(player *Player, data []byte) {
	reader := &PacketReader{data: data, cursor: 0}
	packetId := reader.FetchVarInt()

	//fmt.Printf("Got packet %d\n", packetId)

	switch packetId {
	case 0x00:
		if reader.BytesLeft() > 0 {
			player.OnHandshakeRequest(ReadHandshakeRequest(reader))
		} else {
			player.OnStatusRequest()
		}
	case 0x01:
		player.OnPing(ReadPingRequest(reader))
	}
}

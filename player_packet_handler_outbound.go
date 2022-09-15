package main

import (
	"github.com/mkorman9/go-minecraft-server/packets"
	"log"
)

func (pph *PlayerPacketHandler) Cancel(reason *ChatMessage) {
	pph.canceledMutex.Lock()
	if pph.canceled {
		return
	}
	pph.canceled = true
	pph.canceledMutex.Unlock()

	switch pph.state {
	case PlayerStateBeforeHandshake:
		// nop
	case PlayerStateLogin:
		_ = pph.sendCancelLogin(reason)
	case PlayerStateEncryption:
		_ = pph.sendCancelLogin(reason)
	case PlayerStatePlay:
		_ = pph.sendDisconnect(reason)
	}

	_ = pph.connection.Close()
}

func (pph *PlayerPacketHandler) SendSystemChatMessage(message *ChatMessage) error {
	return pph.sendSystemChatMessage(message)
}

func (pph *PlayerPacketHandler) SynchronizePosition(x float64, y float64, z float64) error {
	return pph.sendPositionUpdate(x, y, z)
}

func (pph *PlayerPacketHandler) SendKeepAlive(keepAliveID int64) error {
	return pph.sendKeepAlive(keepAliveID)
}

func (pph *PlayerPacketHandler) sendHandshakeStatusResponse() error {
	serverStatus := pph.world.GetStatus()
	serverStatusJSON, err := serverStatus.Encode()
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	handshakeResponse := HandshakeResponse.
		New().
		Set("statusJson", serverStatusJSON)

	return pph.packetWriter.Write(handshakeResponse)
}

func (pph *PlayerPacketHandler) sendPongResponse(payload int64) error {
	pongResponse := PongResponse.
		New().
		Set("payload", payload)

	return pph.packetWriter.Write(pongResponse)
}

func (pph *PlayerPacketHandler) sendEncryptionRequest() error {
	encryptionRequest := EncryptionRequest.
		New().
		Set("serverId", "").
		Set("publicKey", pph.world.Server().PublicKey()).
		Set("verifyToken", pph.verifyToken)

	return pph.packetWriter.Write(encryptionRequest)
}

func (pph *PlayerPacketHandler) sendSetCompressionRequest(compressionThreshold int) error {
	setCompressionRequest := SetCompressionRequest.
		New().
		Set("threshold", compressionThreshold)

	return pph.packetWriter.Write(setCompressionRequest)
}

func (pph *PlayerPacketHandler) sendCancelLogin(reason *ChatMessage) error {
	cancelLoginPacket := CancelLoginPacket.
		New().
		Set("reason", reason.Encode())

	return pph.packetWriter.Write(cancelLoginPacket)
}

func (pph *PlayerPacketHandler) sendLoginSuccessResponse() error {
	loginSuccessResponse := LoginSuccessResponse.
		New().
		Set("uuid", pph.player.UUID).
		Set("username", pph.player.Name)

	return pph.packetWriter.Write(loginSuccessResponse)
}

func (pph *PlayerPacketHandler) sendPlayPacket(entityID int32) error {
	playPacket := PlayPacket.
		New().
		Set("entityID", entityID).
		Set("isHardcore", pph.world.Data().IsHardcore).
		Set("gameMode", pph.world.Data().GameMode).
		Set("previousGameMode", GameModeUnknown).
		SetArray(
			"worldNames",
			packets.ConvertArrayValue(
				pph.world.Data().WorldNames,
				func(value string, packet *packets.PacketData) {
					packet.Set("value", value)
				},
			),
		).
		Set("dimensionCodec", pph.world.Data().DimensionCodec).
		Set("worldType", pph.world.Data().SpawnDimension).
		Set("worldName", pph.world.Data().SpawnDimension).
		Set("hashedSeed", pph.world.Data().HashedSeed).
		Set("maxPlayers", pph.world.Settings().MaxPlayers).
		Set("viewDistance", pph.world.Settings().ViewDistance).
		Set("simulationDistance", pph.world.Settings().SimulationDistance).
		Set("reducedDebugInfo", !pph.world.Settings().IsDebug).
		Set("enableRespawnScreen", pph.world.Data().EnableRespawnScreen).
		Set("isDebug", pph.world.Settings().IsDebug).
		Set("isFlat", pph.world.Data().IsFlat).
		Set("hasDeath", false)

	return pph.packetWriter.Write(playPacket)
}

func (pph *PlayerPacketHandler) sendDisconnect(reason *ChatMessage) error {
	disconnectPacket := DisconnectPacket.
		New().
		Set("reason", reason.Encode())

	return pph.packetWriter.Write(disconnectPacket)
}

func (pph *PlayerPacketHandler) sendSystemChatMessage(message *ChatMessage) error {
	systemChatPacket := SystemChatPacket.
		New().
		Set("content", message.Encode()).
		Set("type", SystemChatMessageTypeChat)

	return pph.packetWriter.Write(systemChatPacket)
}

func (pph *PlayerPacketHandler) sendSpawnPosition() error {
	spawnPositionPacket := SpawnPositionPacket.
		New().
		Set("location", pph.world.Data().SpawnPosition).
		Set("angle", float32(0))

	return pph.packetWriter.Write(spawnPositionPacket)
}

func (pph *PlayerPacketHandler) sendPositionUpdate(x float64, y float64, z float64) error {
	updatePositionPacket := UpdatePositionPacket.
		New().
		Set("x", x).
		Set("y", y).
		Set("z", z).
		Set("yaw", pph.player.Yaw).
		Set("pitch", pph.player.Pitch).
		Set("flags", 0).
		Set("teleportId", 0).
		Set("dismountVehicle", false)

	return pph.packetWriter.Write(updatePositionPacket)
}

func (pph *PlayerPacketHandler) sendKeepAlive(keepAliveID int64) error {
	keepAlivePacket := KeepAlivePacket.
		New().
		Set("keepAliveId", keepAliveID)

	return pph.packetWriter.Write(keepAlivePacket)
}

func (pph *PlayerPacketHandler) sendPlayerListUpdate(players []*Player) error {
	playerInfoPacket := PlayerInfoPacket.
		New().
		Set("actionId", 0).
		SetArray(
			"playersToAdd",
			packets.ConvertArrayValue(players, func(player *Player, packet *packets.PacketData) {
				packet.Set("uuid", player.UUID).
					Set("name", player.Name)

				var properties []SignedProperties
				if player.Textures != "" {
					properties = append(properties, SignedProperties{
						Name:      "textures",
						Value:     player.TexturesSignature,
						IsSigned:  true,
						Signature: player.Signature,
					})
				}

				packet.SetArray(
					"properties",
					packets.ConvertArrayValue(properties, func(property SignedProperties, packet2 *packets.PacketData) {
						packet2.Set("name", property.Name)
						packet2.Set("value", property.Value)
						packet2.Set("isSigned", property.IsSigned)
						packet2.Set("signature", property.Signature)
					}),
				)

				packet.Set("gameMode", int(player.GameMode)).
					Set("ping", player.Ping).
					Set("hasDisplayName", true).
					Set("displayName", player.DisplayName.Encode()).
					Set("hasSigData", true).
					Set("timestamp", player.Timestamp).
					Set("publicKey", player.PublicKeyDER).
					Set("signature", player.Signature)
			}),
		)

	return pph.packetWriter.Write(playerInfoPacket)
}

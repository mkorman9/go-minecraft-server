package main

import "log"

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

	response := &HandshakeResponse{
		StatusJSON: serverStatusJSON,
	}

	return pph.writePacket(response)
}

func (pph *PlayerPacketHandler) sendPongResponse(payload int64) error {
	packet := &PongResponse{
		Payload: payload,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendEncryptionRequest() error {
	response := &EncryptionRequest{
		ServerID:    "",
		PublicKey:   pph.world.Server().PublicKey(),
		VerifyToken: pph.verifyToken,
	}

	return pph.writePacket(response)
}

func (pph *PlayerPacketHandler) sendSetCompressionRequest(compressionThreshold int) error {
	request := &SetCompressionRequest{
		Threshold: compressionThreshold,
	}

	return pph.writePacket(request)
}

func (pph *PlayerPacketHandler) sendCancelLogin(reason *ChatMessage) error {
	packet := &CancelLoginPacket{
		Reason: reason,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendLoginSuccessResponse() error {
	response := &LoginSuccessResponse{
		UUID:     pph.player.UUID,
		Username: pph.player.Name,
	}

	return pph.writePacket(response)
}

func (pph *PlayerPacketHandler) sendPlayPacket(entityID int32) error {
	packet := &PlayPacket{
		EntityID:            entityID,
		IsHardcore:          pph.world.Data().IsHardcore,
		GameMode:            pph.world.Data().GameMode,
		PreviousGameMode:    GameModeUnknown,
		WorldNames:          pph.world.Data().WorldNames,
		DimensionCodec:      *pph.world.Data().DimensionCodec,
		WorldType:           pph.world.Data().SpawnDimension,
		WorldName:           pph.world.Data().SpawnDimension,
		HashedSeed:          pph.world.Data().HashedSeed,
		MaxPlayers:          pph.world.Settings().MaxPlayers,
		ViewDistance:        pph.world.Settings().ViewDistance,
		SimulationDistance:  pph.world.Settings().SimulationDistance,
		ReducedDebugInfo:    !pph.world.Settings().IsDebug,
		EnableRespawnScreen: pph.world.Data().EnableRespawnScreen,
		IsDebug:             pph.world.Settings().IsDebug,
		IsFlat:              pph.world.Data().IsFlat,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendDisconnect(reason *ChatMessage) error {
	packet := &DisconnectPacket{
		Reason: reason,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendSystemChatMessage(message *ChatMessage) error {
	packet := &SystemChatPacket{
		Content: message,
		Type:    SystemChatMessageTypeChat,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendSpawnPosition() error {
	packet := &SpawnPositionPacket{
		Location: pph.world.Data().SpawnPosition,
		Angle:    0,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendPositionUpdate(x float64, y float64, z float64) error {
	packet := &UpdatePositionPacket{
		X:               x,
		Y:               y,
		Z:               z,
		Yaw:             pph.player.Yaw,
		Pitch:           pph.player.Pitch,
		Flags:           0,
		TeleportID:      0,
		DismountVehicle: false,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) sendKeepAlive(keepAliveID int64) error {
	packet := &KeepAlivePacket{
		KeepAliveID: keepAliveID,
	}

	return pph.writePacket(packet)
}

func (pph *PlayerPacketHandler) writePacket(packet Packet) error {
	data, err := packet.Marshal(pph.packetWriter.New())
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	_, err = pph.writer.Write(data)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	return nil
}

package main

/*
	0x00: Handshake Response
*/

type HandshakeResponse struct {
	StatusJSON string
}

func (hr *HandshakeResponse) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x00)
	writer.AppendString(hr.StatusJSON)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (hr *HandshakeResponse) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x01: Pong
*/

type PongResponse struct {
	Payload int64
}

func (pr *PongResponse) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x01)
	writer.AppendInt64(pr.Payload)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (pr *PongResponse) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x01: Encryption Request
*/

type EncryptionRequest struct {
	ServerID    string
	PublicKey   []byte
	VerifyToken string
}

func (er *EncryptionRequest) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x01)
	writer.AppendString(er.ServerID)
	writer.AppendByteArray(er.PublicKey)
	writer.AppendString(er.VerifyToken)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (er *EncryptionRequest) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x00: Cancel Login
*/

type CancelLoginPacket struct {
	Reason *ChatMessage
}

func (clp *CancelLoginPacket) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x00)
	writer.AppendString(clp.Reason.Encode())

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (clp *CancelLoginPacket) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x02: Login Success
*/

type LoginSuccessResponse struct {
	UUID       UUID
	Username   string
	Properties []LoginSuccessResponseProperty
}

type LoginSuccessResponseProperty struct {
	Name      string
	Value     string
	IsSigned  bool
	Signature string
}

func (lsr *LoginSuccessResponse) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x02)
	writer.AppendUUID(lsr.UUID)
	writer.AppendString(lsr.Username)
	writer.AppendVarInt(len(lsr.Properties))

	for _, property := range lsr.Properties {
		writer.AppendString(property.Name)
		writer.AppendString(property.Value)
		writer.AppendBool(property.IsSigned)

		if property.IsSigned {
			writer.AppendString(property.Signature)
		}
	}

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (lsr *LoginSuccessResponse) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x03: Set Compression
*/

type SetCompressionRequest struct {
	Threshold int
}

func (scr *SetCompressionRequest) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x03)
	writer.AppendVarInt(scr.Threshold)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (scr *SetCompressionRequest) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x23: Play
*/

type PlayPacket struct {
	EntityID            int32
	IsHardcore          bool
	GameMode            GameMode
	PreviousGameMode    GameMode
	WorldNames          []string
	DimensionCodec      DimensionCodec
	WorldType           string
	WorldName           string
	HashedSeed          int64
	MaxPlayers          int
	ViewDistance        int
	SimulationDistance  int
	ReducedDebugInfo    bool
	EnableRespawnScreen bool
	IsDebug             bool
	IsFlat              bool
	DeathDimensionName  string
	DeathLocation       *Position
}

func (pp *PlayPacket) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x23)
	writer.AppendInt32(pp.EntityID)
	writer.AppendBool(pp.IsHardcore)
	writer.AppendByte(pp.GameMode)
	writer.AppendByte(pp.PreviousGameMode)

	writer.AppendVarInt(len(pp.WorldNames))
	for _, world := range pp.WorldNames {
		writer.AppendString(world)
	}

	writer.AppendNBT(&pp.DimensionCodec)
	writer.AppendString(pp.WorldType)
	writer.AppendString(pp.WorldName)
	writer.AppendInt64(pp.HashedSeed)
	writer.AppendVarInt(pp.MaxPlayers)
	writer.AppendVarInt(pp.ViewDistance)
	writer.AppendVarInt(pp.SimulationDistance)
	writer.AppendBool(pp.ReducedDebugInfo)
	writer.AppendBool(pp.EnableRespawnScreen)
	writer.AppendBool(pp.IsDebug)
	writer.AppendBool(pp.IsFlat)

	if pp.DeathDimensionName != "" && pp.DeathLocation != nil {
		writer.AppendBool(true)
		writer.AppendString(pp.DeathDimensionName)
		writer.AppendPosition(pp.DeathLocation)
	} else {
		writer.AppendBool(false)
	}

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (pp *PlayPacket) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x4a: Spawn Position
*/

type SpawnPositionPacket struct {
	Location *Position
	Angle    float32
}

func (spp *SpawnPositionPacket) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x4a)
	writer.AppendPosition(spp.Location)
	writer.AppendFloat32(spp.Angle)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (spp *SpawnPositionPacket) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x17: Disconnect
*/

type DisconnectPacket struct {
	Reason *ChatMessage
}

func (dp *DisconnectPacket) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x17)
	writer.AppendString(dp.Reason.Encode())

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (dp *DisconnectPacket) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x1e: Keep Alive
*/

type KeepAlivePacket struct {
	KeepAliveID int64
}

func (kap *KeepAlivePacket) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x1e)
	writer.AppendInt64(kap.KeepAliveID)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (kap *KeepAlivePacket) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x5f: System Chat
*/

type SystemChatPacket struct {
	Content *ChatMessage
	Type    SystemChatMessageType
}

func (scp *SystemChatPacket) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x5f)
	writer.AppendString(scp.Content.Encode())
	writer.AppendVarInt(scp.Type)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (scp *SystemChatPacket) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x36: Update Position
*/

type UpdatePositionPacket struct {
	X               float64
	Y               float64
	Z               float64
	Yaw             float32
	Pitch           float32
	Flags           byte
	TeleportID      int
	DismountVehicle bool
}

func (upp *UpdatePositionPacket) Marshal(writer *PacketSerializer) ([]byte, error) {
	writer.AppendByte(0x36)
	writer.AppendFloat64(upp.X)
	writer.AppendFloat64(upp.Y)
	writer.AppendFloat64(upp.Z)
	writer.AppendFloat32(upp.Yaw)
	writer.AppendFloat32(upp.Pitch)
	writer.AppendByte(upp.Flags)
	writer.AppendVarInt(upp.TeleportID)
	writer.AppendBool(upp.DismountVehicle)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (upp *UpdatePositionPacket) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

/*
	0x1f: Map Chunk
*/

type MapChunkPacket struct {
	X                   int32
	Z                   int32
	Heightmaps          Heightmap
	ChunkData           ChunkData
	BlockEntities       []BlockEntity
	TrustEdges          bool
	SkyLightMask        *BitSet
	BlockLightMask      *BitSet
	EmptySkyLightMask   *BitSet
	EmptyBlockLightMask *BitSet
}

func (mcp *MapChunkPacket) Marshal(writer *PacketSerializer) ([]byte, error) {
	skyLightMaskBits := mcp.SkyLightMask.BitsSet()
	blockLightMaskBits := mcp.BlockLightMask.BitsSet()

	writer.AppendByte(0x1f)
	writer.AppendInt32(mcp.X)
	writer.AppendInt32(mcp.Z)
	writer.AppendNBT(&mcp.Heightmaps)

	for _, chunkSection := range mcp.ChunkData.Data {
		writer.AppendInt16(chunkSection.BlockCount)

		for _, blockState := range chunkSection.BlockStates {
			writer.AppendByte(blockState.BitsPerEntry)

			switch {
			case blockState.BitsPerEntry == 0:
				writer.AppendVarInt(blockState.PaletteSingleValued.Value)
			case blockState.BitsPerEntry == 9:
			default:
				writer.AppendVarInt(len(blockState.PaletteIndirect.Palette))
				for _, p := range blockState.PaletteIndirect.Palette {
					writer.AppendVarInt(p)
				}
			}

			writer.AppendVarInt(len(blockState.Data))
			for _, d := range blockState.Data {
				writer.AppendInt64(d)
			}
		}

		for _, biome := range chunkSection.BlockStates {
			writer.AppendByte(biome.BitsPerEntry)

			switch {
			case biome.BitsPerEntry == 0:
				writer.AppendVarInt(biome.PaletteSingleValued.Value)
			case biome.BitsPerEntry == 4:
			default:
				writer.AppendVarInt(len(biome.PaletteIndirect.Palette))
				for _, p := range biome.PaletteIndirect.Palette {
					writer.AppendVarInt(p)
				}
			}

			writer.AppendVarInt(len(biome.Data))
			for _, d := range biome.Data {
				writer.AppendInt64(d)
			}
		}
	}

	writer.AppendVarInt(len(mcp.BlockEntities))
	for _, entity := range mcp.BlockEntities {
		writer.AppendByte(entity.PackedXZ)
		writer.AppendInt16(entity.Y)
		writer.AppendVarInt(entity.Type)
		writer.AppendNBT(&entity.Data)
	}

	writer.AppendBool(mcp.TrustEdges)
	writer.AppendBitSet(mcp.SkyLightMask)
	writer.AppendBitSet(mcp.BlockLightMask)
	writer.AppendBitSet(mcp.EmptySkyLightMask)
	writer.AppendBitSet(mcp.EmptyBlockLightMask)

	writer.AppendVarInt(skyLightMaskBits)
	for i := 0; i < skyLightMaskBits; i++ {
		writer.AppendVarInt(2048)
		writer.AppendByteArray(make([]byte, 2048))
	}

	writer.AppendVarInt(blockLightMaskBits)
	for i := 0; i < blockLightMaskBits; i++ {
		writer.AppendVarInt(2048)
		writer.AppendByteArray(make([]byte, 2048))
	}

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (mcp *MapChunkPacket) Unmarshal(reader *PacketDeserializer) error {
	return nil
}

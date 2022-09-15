package main

import "github.com/mkorman9/go-minecraft-server/packets"

/*
	0x00: Handshake Response
*/

var HandshakeResponse = packets.Packet(
	packets.ID(0x00),
	packets.String("statusJson"),
)

/*
	0x01: Pong
*/

var PongResponse = packets.Packet(
	packets.ID(0x01),
	packets.Int64("payload"),
)

/*
	0x01: Encryption Request
*/

var EncryptionRequest = packets.Packet(
	packets.ID(0x01),
	packets.String("serverId"),
	packets.ByteArray("publicKey"),
	packets.String("verifyToken"),
)

/*
	0x00: Cancel Login
*/

var CancelLoginPacket = packets.Packet(
	packets.ID(0x00),
	packets.String("reason"),
)

/*
	0x02: Login Success
*/

var LoginSuccessResponse = packets.Packet(
	packets.ID(0x02),
	packets.UUIDField("uuid"),
	packets.String("username"),
	packets.Array(
		"properties",
		packets.ArrayLengthPrefixed,
		packets.String("name"),
		packets.String("value"),
		packets.Bool("isSigned"),
		packets.String("signature", packets.OnlyIfTrue("isSigned")),
	),
)

/*
	0x03: Set Compression
*/

var SetCompressionRequest = packets.Packet(
	packets.ID(0x03),
	packets.VarInt("threshold"),
)

/*
	0x23: Play
*/

var PlayPacket = packets.Packet(
	packets.ID(0x23),
	packets.Int32("entityID"),
	packets.Bool("isHardcore"),
	packets.Byte("gameMode"),
	packets.Byte("previousGameMode"),
	packets.Array(
		"worldNames",
		packets.ArrayLengthPrefixed,
		packets.String("value"),
	),
	packets.NBT("dimensionCodec", &DimensionCodec{}),
	packets.String("worldType"),
	packets.String("worldName"),
	packets.Int64("hashedSeed"),
	packets.VarInt("maxPlayers"),
	packets.VarInt("viewDistance"),
	packets.VarInt("simulationDistance"),
	packets.Bool("reducedDebugInfo"),
	packets.Bool("enableRespawnScreen"),
	packets.Bool("isDebug"),
	packets.Bool("isFlat"),
	packets.Bool("hasDeath"),
	packets.String("deathDimension", packets.OnlyIfTrue("hasDeath")),
	packets.PositionField("deathLocation", packets.OnlyIfTrue("hasDeath")),
)

/*
	0x4a: Spawn Position
*/

var SpawnPositionPacket = packets.Packet(
	packets.ID(0x4a),
	packets.PositionField("location"),
	packets.Float32("angle"),
)

/*
	0x17: Disconnect
*/

var DisconnectPacket = packets.Packet(
	packets.ID(0x17),
	packets.String("reason"),
)

/*
	0x1e: Keep Alive
*/

var KeepAlivePacket = packets.Packet(
	packets.ID(0x1e),
	packets.Int64("keepAliveId"),
)

/*
	0x5f: System Chat
*/

var SystemChatPacket = packets.Packet(
	packets.ID(0x5f),
	packets.String("content"),
	packets.VarInt("type"),
)

/*
	0x36: Update Position
*/

var UpdatePositionPacket = packets.Packet(
	packets.ID(0x36),
	packets.Float64("x"),
	packets.Float64("y"),
	packets.Float64("z"),
	packets.Float32("yaw"),
	packets.Float32("pitch"),
	packets.Byte("flags"),
	packets.VarInt("teleportId"),
	packets.Bool("dismountVehicle"),
)

/*
	0x34: Player Info
*/

var PlayerInfoPacket = packets.Packet(
	packets.ID(0x34),
	packets.VarInt("actionId"),
	packets.ArrayWithOptions(
		"playersToAdd",
		packets.ArrayLengthPrefixed,
		[]packets.PacketOpt{
			packets.UUIDField("uuid"),
			packets.String("name"),
			packets.Array(
				"properties",
				packets.ArrayLengthPrefixed,
				packets.String("name"),
				packets.String("value"),
				packets.Bool("isSigned"),
				packets.String("signature", packets.OnlyIfTrue("isSigned")),
			),
			packets.VarInt("gameMode"),
			packets.VarInt("ping"),
			packets.Bool("hasDisplayName"),
			packets.String("displayName", packets.OnlyIfTrue("hasDisplayName")),
			packets.Bool("hasSigData"),
			packets.Int64("timestamp", packets.OnlyIfTrue("hasSigData")),
			packets.ByteArray("publicKey", packets.OnlyIfTrue("hasSigData")),
			packets.String("signature", packets.OnlyIfTrue("hasSigData")),
		},
		packets.OnlyIfEqual("actionId", 0),
	),
	packets.ArrayWithOptions(
		"playersToUpdateGameMode",
		packets.ArrayLengthPrefixed,
		[]packets.PacketOpt{
			packets.UUIDField("uuid"),
			packets.VarInt("gameMode"),
		},
		packets.OnlyIfEqual("actionId", 1),
	),
	packets.ArrayWithOptions(
		"playersToUpdateLatency",
		packets.ArrayLengthPrefixed,
		[]packets.PacketOpt{
			packets.UUIDField("uuid"),
			packets.VarInt("ping"),
		},
		packets.OnlyIfEqual("actionId", 2),
	),
	packets.ArrayWithOptions(
		"playersToUpdateDisplayName",
		packets.ArrayLengthPrefixed,
		[]packets.PacketOpt{
			packets.UUIDField("uuid"),
			packets.Bool("hasDisplayName"),
			packets.String("displayName", packets.OnlyIfTrue("hasDisplayName")),
		},
		packets.OnlyIfEqual("actionId", 3),
	),
	packets.ArrayWithOptions(
		"playersToRemove",
		packets.ArrayLengthPrefixed,
		[]packets.PacketOpt{
			packets.UUIDField("uuid"),
		},
		packets.OnlyIfEqual("actionId", 4),
	),
)

/*
	0x1f: Map Chunk
*/

//type MapChunkPacket struct {
//	X                   int32
//	Z                   int32
//	Heightmaps          Heightmap
//	ChunkData           ChunkData
//	BlockEntities       []BlockEntity
//	TrustEdges          bool
//	SkyLightMask        *BitSet
//	BlockLightMask      *BitSet
//	EmptySkyLightMask   *BitSet
//	EmptyBlockLightMask *BitSet
//}
//
//func (mcp *MapChunkPacket) Marshal(writer *PacketSerializer) ([]byte, error) {
//	skyLightMaskBits := mcp.SkyLightMask.BitsSet()
//	blockLightMaskBits := mcp.BlockLightMask.BitsSet()
//
//	writer.AppendByte(0x1f)
//	writer.AppendInt32(mcp.X)
//	writer.AppendInt32(mcp.Z)
//	writer.AppendNBT(&mcp.Heightmaps)
//
//	for _, chunkSection := range mcp.ChunkData.Data {
//		writer.AppendInt16(chunkSection.BlockCount)
//
//		for _, blockState := range chunkSection.BlockStates {
//			writer.AppendByte(blockState.BitsPerEntry)
//
//			switch {
//			case blockState.BitsPerEntry == 0:
//				writer.AppendVarInt(blockState.PaletteSingleValued.Value)
//			case blockState.BitsPerEntry == 9:
//			default:
//				writer.AppendVarInt(len(blockState.PaletteIndirect.Palette))
//				for _, p := range blockState.PaletteIndirect.Palette {
//					writer.AppendVarInt(p)
//				}
//			}
//
//			writer.AppendVarInt(len(blockState.Data))
//			for _, d := range blockState.Data {
//				writer.AppendInt64(d)
//			}
//		}
//
//		for _, biome := range chunkSection.BlockStates {
//			writer.AppendByte(biome.BitsPerEntry)
//
//			switch {
//			case biome.BitsPerEntry == 0:
//				writer.AppendVarInt(biome.PaletteSingleValued.Value)
//			case biome.BitsPerEntry == 4:
//			default:
//				writer.AppendVarInt(len(biome.PaletteIndirect.Palette))
//				for _, p := range biome.PaletteIndirect.Palette {
//					writer.AppendVarInt(p)
//				}
//			}
//
//			writer.AppendVarInt(len(biome.Data))
//			for _, d := range biome.Data {
//				writer.AppendInt64(d)
//			}
//		}
//	}
//
//	writer.AppendVarInt(len(mcp.BlockEntities))
//	for _, entity := range mcp.BlockEntities {
//		writer.AppendByte(entity.PackedXZ)
//		writer.AppendInt16(entity.Y)
//		writer.AppendVarInt(entity.Type)
//		writer.AppendNBT(&entity.Data)
//	}
//
//	writer.AppendBool(mcp.TrustEdges)
//	writer.AppendBitSet(mcp.SkyLightMask)
//	writer.AppendBitSet(mcp.BlockLightMask)
//	writer.AppendBitSet(mcp.EmptySkyLightMask)
//	writer.AppendBitSet(mcp.EmptyBlockLightMask)
//
//	writer.AppendVarInt(skyLightMaskBits)
//	for i := 0; i < skyLightMaskBits; i++ {
//		writer.AppendVarInt(2048)
//		writer.AppendByteArray(make([]byte, 2048))
//	}
//
//	writer.AppendVarInt(blockLightMaskBits)
//	for i := 0; i < blockLightMaskBits; i++ {
//		writer.AppendVarInt(2048)
//		writer.AppendByteArray(make([]byte, 2048))
//	}
//
//	if writer.Error() != nil {
//		return nil, writer.Error()
//	}
//
//	return writer.Bytes(), nil
//}
//
//func (mcp *MapChunkPacket) Unmarshal(reader *PacketDeserializer) error {
//	return nil
//}

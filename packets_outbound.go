package main

import (
	"github.com/mkorman9/go-minecraft-server/chunk"
	"github.com/mkorman9/go-minecraft-server/nbt"
	"github.com/mkorman9/go-minecraft-server/packets"
)

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
		packets.Fields(
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
		),
		packets.OnlyIfEqual("actionId", 0),
	),
	packets.ArrayWithOptions(
		"playersToUpdateGameMode",
		packets.ArrayLengthPrefixed,
		packets.Fields(
			packets.UUIDField("uuid"),
			packets.VarInt("gameMode"),
		),
		packets.OnlyIfEqual("actionId", 1),
	),
	packets.ArrayWithOptions(
		"playersToUpdateLatency",
		packets.ArrayLengthPrefixed,
		packets.Fields(
			packets.UUIDField("uuid"),
			packets.VarInt("ping"),
		),
		packets.OnlyIfEqual("actionId", 2),
	),
	packets.ArrayWithOptions(
		"playersToUpdateDisplayName",
		packets.ArrayLengthPrefixed,
		packets.Fields(
			packets.UUIDField("uuid"),
			packets.Bool("hasDisplayName"),
			packets.String("displayName", packets.OnlyIfTrue("hasDisplayName")),
		),
		packets.OnlyIfEqual("actionId", 3),
	),
	packets.ArrayWithOptions(
		"playersToRemove",
		packets.ArrayLengthPrefixed,
		packets.Fields(
			packets.UUIDField("uuid"),
		),
		packets.OnlyIfEqual("actionId", 4),
	),
)

/*
	0x1f: Map Chunk
*/

var MapChunkPacket = packets.Packet(
	packets.ID(0x1f),
	packets.Int32("x"),
	packets.Int32("z"),
	packets.NBT("heightmaps", &chunk.Heightmap{}),
	packets.ByteArray("data"),
	packets.Array(
		"blockEntities",
		packets.ArrayLengthPrefixed,
		packets.Byte("xz"),
		packets.Int16("y"),
		packets.VarInt("type"),
		packets.NBT("data", &nbt.RawMessage{}),
	),
	packets.Bool("trustEdges"),
	packets.BitSetField("skyLightMask"),
	packets.BitSetField("blockLightMask"),
	packets.BitSetField("emptySkyLightMask"),
	packets.BitSetField("emptyBlockLightMask"),
	packets.Array(
		"skyLights",
		packets.ArrayLengthPrefixed,
		packets.ByteArray("value"),
	),
	packets.Array(
		"blockLights",
		packets.ArrayLengthPrefixed,
		packets.ByteArray("value"),
	),
)

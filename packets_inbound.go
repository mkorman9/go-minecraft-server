package main

import "github.com/mkorman9/go-minecraft-server/packets"

/*
	0x00: Handshake
*/

var HandshakeRequest = packets.Packet(
	packets.ID(0x00),
	packets.VarInt("protocolVersion"),
	packets.String("serverAddress"),
	packets.Int16("serverPort"),
	packets.VarInt("nextState"),
)

/*
	0x01: Ping
*/

var PingRequest = packets.Packet(
	packets.ID(0x01),
	packets.Int64("payload"),
)

/*
	0x00: Login Start
*/

var LoginStartRequest = packets.Packet(
	packets.ID(0x00),
	packets.String("name"),
	packets.Bool("hasSigData"),
	packets.Int64("timestamp", packets.OnlyIfTrue("hasSigData")),
	packets.ByteArray("publicKey", packets.OnlyIfTrue("hasSigData")),
	packets.String("signature", packets.OnlyIfTrue("hasSigData")),
)

/*
	0x01: Encryption Response
*/

var EncryptionResponse = packets.Packet(
	packets.ID(0x01),
	packets.ByteArray("sharedSecret"),
	packets.Bool("hasVerifyToken"),
	packets.ByteArray("verifyToken", packets.OnlyIfTrue("hasVerifyToken")),
	packets.Int64("salt", packets.OnlyIfFalse("hasVerifyToken")),
	packets.ByteArray("messageSignature", packets.OnlyIfFalse("hasVerifyToken")),
)

/*
	0x07: Settings
*/

var SettingsPacket = packets.Packet(
	packets.ID(0x07),
	packets.String("locale"),
	packets.Byte("viewDistance"),
	packets.VarInt("chatFlags"),
	packets.Bool("chatColors"),
	packets.Byte("skinParts"),
	packets.VarInt("mainHand"),
	packets.Bool("enableTextFiltering"),
	packets.Bool("enableServerListing"),
)

/*
	0x0c: Custom Payload
*/

var CustomPayloadPacket = packets.Packet(
	packets.ID(0x0c),
	packets.String("channel"),
	packets.ByteArray("data"),
)

/*
	0x13: Position
*/

var PositionPacket = packets.Packet(
	packets.ID(0x13),
	packets.Float64("x"),
	packets.Float64("y"),
	packets.Float64("z"),
	packets.Bool("onGround"),
)

/*
	0x14: Position & Look
*/

var PositionLookPacket = packets.Packet(
	packets.ID(0x14),
	packets.Float64("x"),
	packets.Float64("y"),
	packets.Float64("z"),
	packets.Float32("yaw"),
	packets.Float32("pitch"),
	packets.Bool("onGround"),
)

/*
	0x15: Look
*/

var LookPacket = packets.Packet(
	packets.ID(0x15),
	packets.Float32("yaw"),
	packets.Float32("pitch"),
	packets.Bool("onGround"),
)

/*
	0x2e: Arm Animation (left click)
*/

var ArmAnimationPacket = packets.Packet(
	packets.ID(0x2e),
	packets.VarInt("hand"),
)

/*
	0x1b: Abilities
*/

var AbilitiesPacket = packets.Packet(
	packets.ID(0x1b),
	packets.Byte("flags"),
)

/*
	0x2a: SetCreativeSlot
*/

var SetCreativeSlotPacket = packets.Packet(
	packets.ID(0x2a),
	packets.Int16("slot"),
	packets.Slot("item"),
)

/*
	0x04: Chat Message
*/

var ChatMessagePacket = packets.Packet(
	packets.ID(0x04),
	packets.String("message"),
	packets.Int64("timestamp"),
	packets.Int64("salt"),
	packets.ByteArray("signature"),
	packets.Bool("signedPreview"),
)

/*
	0x03: Chat Command
*/

var ChatCommandPacket = packets.Packet(
	packets.ID(0x03),
	packets.String("message"),
	packets.Int64("timestamp"),
	packets.Int64("salt"),
	packets.Array(
		"arguments",
		packets.ArrayLengthPrefixed,
		packets.String("name"),
		packets.ByteArray("signature"),
	),
	packets.ByteArray("signature"),
)

/*
	0x00: Teleport confirm
*/

var TeleportConfirmPacket = packets.Packet(
	packets.ID(0x00),
	packets.VarInt("teleportId"),
)

/*
	0x11: Keep Alive Response
*/

var KeepAliveResponsePacket = packets.Packet(
	packets.ID(0x11),
	packets.Int64("keepAliveId"),
)

/*
	0x1d: Entity Action
*/

var EntityActionPacket = packets.Packet(
	packets.ID(0x1d),
	packets.VarInt("entityId"),
	packets.VarInt("actionId"),
	packets.VarInt("jumpBoost"),
)

/*
	0x0b: Close Window
*/

var CloseWindowPacket = packets.Packet(
	packets.ID(0x0b),
	packets.Byte("windowId"),
)

package main

/*
	0x00: Handshake
*/

type HandshakeRequest struct {
	ProtocolVersion int
	ServerAddress   string
	ServerPort      int16
	NextState       HandshakeType
}

func (hr *HandshakeRequest) Marshal(ctx *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (hr *HandshakeRequest) Unmarshal(reader *PacketReaderContext) error {
	hr.ProtocolVersion = reader.FetchVarInt()
	hr.ServerAddress = reader.FetchString()
	hr.ServerPort = reader.FetchInt16()
	hr.NextState = reader.FetchVarInt()

	return reader.Error()
}

/*
	0x01: Ping
*/

type PingRequest struct {
	Payload int64
}

func (pr *PingRequest) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (pr *PingRequest) Unmarshal(reader *PacketReaderContext) error {
	pr.Payload = reader.FetchInt64()
	return reader.Error()
}

/*
	0x00: Login Start
*/

type LoginStartRequest struct {
	Name      string
	Timestamp int64
	PublicKey []byte
	Signature []byte
}

func (lsr *LoginStartRequest) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (lsr *LoginStartRequest) Unmarshal(reader *PacketReaderContext) error {
	lsr.Name = reader.FetchString()
	hasSigData := reader.FetchBool()

	if hasSigData {
		lsr.Timestamp = reader.FetchInt64()
		lsr.PublicKey = reader.FetchByteArray()
		lsr.Signature = reader.FetchByteArray()
	}

	return reader.Error()
}

/*
	0x01: Encryption Response
*/

type EncryptionResponse struct {
	SharedSecret     []byte
	VerifyToken      []byte
	Salt             int64
	MessageSignature []byte
}

func (er *EncryptionResponse) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (er *EncryptionResponse) Unmarshal(reader *PacketReaderContext) error {
	er.SharedSecret = reader.FetchByteArray()
	hasVerifyToken := reader.FetchBool()

	if hasVerifyToken {
		er.VerifyToken = reader.FetchByteArray()
	} else {
		er.Salt = reader.FetchInt64()
		er.MessageSignature = reader.FetchByteArray()
	}

	return reader.Error()
}

/*
	0x07: Settings
*/

type SettingsPacket struct {
	Locale              string
	ViewDistance        byte
	ChatFlags           int
	ChatColors          bool
	SkinParts           byte
	MainHand            int
	EnableTextFiltering bool
	EnableServerListing bool
}

func (sp *SettingsPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (sp *SettingsPacket) Unmarshal(reader *PacketReaderContext) error {
	sp.Locale = reader.FetchString()
	sp.ViewDistance = reader.FetchByte()
	sp.ChatFlags = reader.FetchVarInt()
	sp.ChatColors = reader.FetchBool()
	sp.SkinParts = reader.FetchByte()
	sp.MainHand = reader.FetchVarInt()
	sp.EnableTextFiltering = reader.FetchBool()
	sp.EnableServerListing = reader.FetchBool()

	return reader.Error()
}

/*
	0x0c: Custom Payload
*/

type CustomPayloadPacket struct {
	Channel string
	Data    []byte
}

func (cpp *CustomPayloadPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (cpp *CustomPayloadPacket) Unmarshal(reader *PacketReaderContext) error {
	cpp.Channel = reader.FetchString()
	cpp.Data = reader.FetchByteArray()

	return reader.Error()
}

/*
	0x13: Position
*/

type PositionPacket struct {
	X        float64
	Y        float64
	Z        float64
	OnGround bool
}

func (pp *PositionPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (pp *PositionPacket) Unmarshal(reader *PacketReaderContext) error {
	pp.X = reader.FetchFloat64()
	pp.Y = reader.FetchFloat64()
	pp.Z = reader.FetchFloat64()
	pp.OnGround = reader.FetchBool()

	return reader.Error()
}

/*
	0x14: Position & Look
*/

type PositionLookPacket struct {
	X        float64
	Y        float64
	Z        float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

func (plp *PositionLookPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (plp *PositionLookPacket) Unmarshal(reader *PacketReaderContext) error {
	plp.X = reader.FetchFloat64()
	plp.Y = reader.FetchFloat64()
	plp.Z = reader.FetchFloat64()
	plp.Yaw = reader.FetchFloat32()
	plp.Pitch = reader.FetchFloat32()
	plp.OnGround = reader.FetchBool()

	return reader.Error()
}

/*
	0x15: Look
*/

type LookPacket struct {
	Yaw      float32
	Pitch    float32
	OnGround bool
}

func (lp *LookPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (lp *LookPacket) Unmarshal(reader *PacketReaderContext) error {
	lp.Yaw = reader.FetchFloat32()
	lp.Pitch = reader.FetchFloat32()
	lp.OnGround = reader.FetchBool()

	return reader.Error()
}

/*
	0x2e: Arm Animation (left click)
*/

type ArmAnimationPacket struct {
	Hand int
}

func (aap *ArmAnimationPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (aap *ArmAnimationPacket) Unmarshal(reader *PacketReaderContext) error {
	aap.Hand = reader.FetchVarInt()

	return reader.Error()
}

/*
	0x1b: Abilities
*/

type AbilitiesPacket struct {
	Flags byte
}

func (ap *AbilitiesPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (ap *AbilitiesPacket) Unmarshal(reader *PacketReaderContext) error {
	ap.Flags = reader.FetchByte()

	return reader.Error()
}

/*
	0x2a: SetCreativeSlot
*/

type SetCreativeSlotPacket struct {
	Slot int16
	Item SlotData
}

func (scsp *SetCreativeSlotPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (scsp *SetCreativeSlotPacket) Unmarshal(reader *PacketReaderContext) error {
	scsp.Slot = reader.FetchInt16()
	scsp.Item = *reader.FetchSlot()

	return reader.Error()
}

/*
	0x04: Chat Message
*/

type ChatMessagePacket struct {
	Message       string
	Timestamp     int64
	Salt          int64
	Signature     []byte
	SignedPreview bool
}

func (cmp *ChatMessagePacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (cmp *ChatMessagePacket) Unmarshal(reader *PacketReaderContext) error {
	cmp.Message = reader.FetchString()
	cmp.Timestamp = reader.FetchInt64()
	cmp.Salt = reader.FetchInt64()
	cmp.Signature = reader.FetchByteArray()
	cmp.SignedPreview = reader.FetchBool()

	return reader.Error()
}

/*
	0x03: Chat Command
*/

type ChatCommandPacket struct {
	Message       string
	Timestamp     int64
	Salt          int64
	Arguments     []ChatCommandPacketArgument
	SignedPreview bool
}

type ChatCommandPacketArgument struct {
	ArgumentName string
	Signature    []byte
}

func (ccp *ChatCommandPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (ccp *ChatCommandPacket) Unmarshal(reader *PacketReaderContext) error {
	ccp.Message = reader.FetchString()
	ccp.Timestamp = reader.FetchInt64()
	ccp.Salt = reader.FetchInt64()

	argumentsCount := reader.FetchVarInt()
	for i := 0; i < argumentsCount; i++ {
		argument := ChatCommandPacketArgument{
			ArgumentName: reader.FetchString(),
			Signature:    reader.FetchByteArray(),
		}

		ccp.Arguments = append(ccp.Arguments, argument)
	}

	ccp.SignedPreview = reader.FetchBool()

	return reader.Error()
}

/*
	0x00: Teleport confirm
*/

type TeleportConfirmPacket struct {
	TeleportID int
}

func (tcp *TeleportConfirmPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (tcp *TeleportConfirmPacket) Unmarshal(reader *PacketReaderContext) error {
	tcp.TeleportID = reader.FetchVarInt()

	return reader.Error()
}

/*
	0x11: Keep Alive Response
*/

type KeepAliveResponsePacket struct {
	KeepAliveID int64
}

func (karp *KeepAliveResponsePacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (karp *KeepAliveResponsePacket) Unmarshal(reader *PacketReaderContext) error {
	karp.KeepAliveID = reader.FetchInt64()

	return reader.Error()
}

/*
	0x1d: Entity Action
*/

type EntityActionPacket struct {
	EntityID  int
	ActionID  EntityAction
	JumpBoost int
}

func (eap *EntityActionPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (eap *EntityActionPacket) Unmarshal(reader *PacketReaderContext) error {
	eap.EntityID = reader.FetchVarInt()
	eap.ActionID = reader.FetchVarInt()
	eap.JumpBoost = reader.FetchVarInt()

	return reader.Error()
}

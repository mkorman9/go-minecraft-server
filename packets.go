package main

type Packet interface {
	Marshal(writer *PacketWriterContext) ([]byte, error)
	Unmarshal(reader *PacketReaderContext) error
}

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

type HandshakeResponse struct {
	StatusJSON string
}

func (hr *HandshakeResponse) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x00)
	writer.AppendString(hr.StatusJSON)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (hr *HandshakeResponse) Unmarshal(reader *PacketReaderContext) error {
	return nil
}

/*
	0x01: Ping/Pong
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

type PongResponse struct {
	Payload int64
}

func (pr *PongResponse) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x01)
	writer.AppendInt64(pr.Payload)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (pr *PongResponse) Unmarshal(reader *PacketReaderContext) error {
	return nil
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
	0x01: Encryption Request
*/

type EncryptionRequest struct {
	ServerID    string
	PublicKey   []byte
	VerifyToken string
}

func (er *EncryptionRequest) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x01)
	writer.AppendString(er.ServerID)
	writer.AppendByteArray(er.PublicKey)
	writer.AppendString(er.VerifyToken)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (er *EncryptionRequest) Unmarshal(reader *PacketReaderContext) error {
	return nil
}

/*
	0x00: Cancel Login
*/

type CancelLoginPacket struct {
	Reason *ChatMessage
}

func (clp *CancelLoginPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x00)
	writer.AppendString(clp.Reason.Encode())

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (clp *CancelLoginPacket) Unmarshal(reader *PacketReaderContext) error {
	return nil
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

func (lsr *LoginSuccessResponse) Marshal(writer *PacketWriterContext) ([]byte, error) {
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

func (lsr *LoginSuccessResponse) Unmarshal(reader *PacketReaderContext) error {
	return nil
}

/*
	0x03: Set Compression
*/

type SetCompressionRequest struct {
	Threshold int
}

func (scr *SetCompressionRequest) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x03)
	writer.AppendVarInt(scr.Threshold)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (scr *SetCompressionRequest) Unmarshal(reader *PacketReaderContext) error {
	return nil
}

/*
	0x17: Disconnect
*/

type DisconnectPacket struct {
	Reason *ChatMessage
}

func (dp *DisconnectPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x17)
	writer.AppendString(dp.Reason.Encode())

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (dp *DisconnectPacket) Unmarshal(reader *PacketReaderContext) error {
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

func (pp *PlayPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
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

func (pp *PlayPacket) Unmarshal(reader *PacketReaderContext) error {
	return nil
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
	0x5f: System Chat
*/

type SystemChatPacket struct {
	Content *ChatMessage
	Type    SystemChatMessageType
}

func (scp *SystemChatPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x5f)
	writer.AppendString(scp.Content.Encode())
	writer.AppendVarInt(scp.Type)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (scp *SystemChatPacket) Unmarshal(reader *PacketReaderContext) error {
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
	0x4a: Spawn Position
*/

type SpawnPositionPacket struct {
	Location *Position
	Angle    float32
}

func (spp *SpawnPositionPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x4a)
	writer.AppendPosition(spp.Location)
	writer.AppendFloat32(spp.Angle)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (spp *SpawnPositionPacket) Unmarshal(reader *PacketReaderContext) error {
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

func (upp *UpdatePositionPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
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

func (upp *UpdatePositionPacket) Unmarshal(reader *PacketReaderContext) error {
	return nil
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
	0x1e: Keep Alive
*/

type KeepAlivePacket struct {
	KeepAliveID int64
}

func (kap *KeepAlivePacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x1e)
	writer.AppendInt64(kap.KeepAliveID)

	if writer.Error() != nil {
		return nil, writer.Error()
	}

	return writer.Bytes(), nil
}

func (kap *KeepAlivePacket) Unmarshal(reader *PacketReaderContext) error {
	return nil
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

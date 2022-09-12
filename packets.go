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
	PublicKey string
	Signature string
}

func (lsr *LoginStartRequest) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (lsr *LoginStartRequest) Unmarshal(reader *PacketReaderContext) error {
	lsr.Name = reader.FetchString()
	hasSigData := reader.FetchBool()

	if hasSigData {
		lsr.Timestamp = reader.FetchInt64()
		lsr.PublicKey = reader.FetchString()
		lsr.Signature = reader.FetchString()
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
	0x19: Disconnect
*/

type DisconnectPacket struct {
	Reason *ChatMessage
}

func (dp *DisconnectPacket) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x19)
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
	GameMode            byte
	PreviousGameMode    byte
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

func (sp *SettingsPacket) Marshal(ctx *PacketWriterContext) ([]byte, error) {
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

func (cpp *CustomPayloadPacket) Marshal(ctx *PacketWriterContext) ([]byte, error) {
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

func (pp *PositionPacket) Marshal(ctx *PacketWriterContext) ([]byte, error) {
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

func (plp *PositionLookPacket) Marshal(ctx *PacketWriterContext) ([]byte, error) {
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
	0x2e: Arm Animation (left click)
*/

type ArmAnimationPacket struct {
	Hand int
}

func (aap *ArmAnimationPacket) Marshal(ctx *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (aap *ArmAnimationPacket) Unmarshal(reader *PacketReaderContext) error {
	aap.Hand = reader.FetchVarInt()

	return reader.Error()
}

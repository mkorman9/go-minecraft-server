package main

type Packet interface {
	Marshal(writer *PacketWriterContext) ([]byte, error)
	Unmarshal(reader *PacketReaderContext) error
}

func UnmarshalPacket[P Packet](reader *PacketReaderContext, p P) (P, error) {
	err := p.Unmarshal(reader)
	return p, err
}

/*
	0x00: Handshake
*/

type HandshakeState = int

const (
	HandshakeStateStatus = 1
	HandshakeStateLogin  = 2
)

type HandshakeRequest struct {
	ProtocolVersion int
	ServerAddress   string
	ServerPort      int16
	NextState       HandshakeState
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
	PublicKey   string
	VerifyToken string
}

func (er *EncryptionRequest) Marshal(writer *PacketWriterContext) ([]byte, error) {
	writer.AppendByte(0x01)
	writer.AppendString(er.ServerID)
	writer.AppendString(er.PublicKey)
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
	SharedSecret     string
	VerifyToken      string
	Salt             int64
	MessageSignature string
}

func (er *EncryptionResponse) Marshal(writer *PacketWriterContext) ([]byte, error) {
	return nil, nil
}

func (er *EncryptionResponse) Unmarshal(reader *PacketReaderContext) error {
	er.SharedSecret = reader.FetchString()
	hasVerifyToken := reader.FetchBool()

	if hasVerifyToken {
		er.VerifyToken = reader.FetchString()
	} else {
		er.Salt = reader.FetchInt64()
		er.MessageSignature = reader.FetchString()
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
	Dimension           Dimension
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
	writer.AppendNBT(&pp.Dimension)
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

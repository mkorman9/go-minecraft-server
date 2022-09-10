package main

type UUID struct {
	Upper int64
	Lower int64
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

func ReadHandshakeRequest(reader *PacketReader) *HandshakeRequest {
	request := &HandshakeRequest{}
	request.ProtocolVersion = reader.FetchVarInt()
	request.ServerAddress = reader.FetchString()
	request.ServerPort = reader.FetchInt16()
	request.NextState = reader.FetchVarInt()
	return request
}

type HandshakeResponse struct {
	StatusJSON string
}

func (hr *HandshakeResponse) Bytes() []byte {
	writer := NewPacketWriter()
	writer.AppendByte(0x00)
	writer.AppendString(hr.StatusJSON)
	return writer.Bytes()
}

/*
	0x01: Ping/Pong
*/

type PingRequest struct {
	Payload int64
}

func ReadPingRequest(reader *PacketReader) *PingRequest {
	request := &PingRequest{}
	request.Payload = reader.FetchInt64()
	return request
}

type PongResponse struct {
	Payload int64
}

func (pr *PongResponse) Bytes() []byte {
	writer := NewPacketWriter()
	writer.AppendByte(0x01)
	writer.AppendInt64(pr.Payload)
	return writer.Bytes()
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

func ReadLoginStartRequest(reader *PacketReader) *LoginStartRequest {
	request := &LoginStartRequest{}
	request.Name = reader.FetchString()
	hasSigData := reader.FetchBool()

	if hasSigData {
		request.Timestamp = reader.FetchInt64()
		request.PublicKey = reader.FetchString()
		request.Signature = reader.FetchString()
	}

	return request
}

/*
	0x01: Encryption Request
*/

type EncryptionRequest struct {
	ServerID    string
	PublicKey   string
	VerifyToken string
}

func (er *EncryptionRequest) Bytes() []byte {
	writer := NewPacketWriter()
	writer.AppendByte(0x01)
	writer.AppendString(er.ServerID)
	writer.AppendString(er.PublicKey)
	writer.AppendString(er.VerifyToken)
	return writer.Bytes()
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

func ReadEncryptionResponse(reader *PacketReader) *EncryptionResponse {
	response := &EncryptionResponse{}
	response.SharedSecret = reader.FetchString()
	hasVerifyToken := reader.FetchBool()

	if hasVerifyToken {
		response.VerifyToken = reader.FetchString()
	} else {
		response.Salt = reader.FetchInt64()
		response.MessageSignature = reader.FetchString()
	}

	return response
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

func (lsr *LoginSuccessResponse) Bytes() []byte {
	writer := NewPacketWriter()
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

	return writer.Bytes()
}

/*
	0x19: Disconnect
*/

type DisconnectPacket struct {
	Reason *ChatMessage
}

func (dp *DisconnectPacket) Bytes() []byte {
	writer := NewPacketWriter()
	writer.AppendByte(0x19)
	writer.AppendString(dp.Reason.Encode())

	return writer.Bytes()
}

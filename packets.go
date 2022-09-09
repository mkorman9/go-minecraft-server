package main

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

package packets

import "bytes"

func getVarIntSize(value int) int {
	var buff bytes.Buffer

	_, _ = Packet(
		VarInt("a"),
	).
		New().
		Set("a", value).
		WriteTo(&buff)

	return buff.Len()
}

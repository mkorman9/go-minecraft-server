package packets

type ArrayValue = []PacketData

type ConvertedArrayValue = func(*PacketDefinition) ArrayValue

func ConvertArrayValue[T any](data []T, convertFunc func(*T, *PacketData)) ConvertedArrayValue {
	return func(definition *PacketDefinition) ArrayValue {
		packets := make([]PacketData, len(data))

		for i, element := range data {
			packetData := definition.New()
			convertFunc(&element, packetData)
			packets[i] = *packetData
		}

		return packets
	}
}

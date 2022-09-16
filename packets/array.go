package packets

type ArrayLengthOption = func(*PacketData) int

const arrayLengthPrefix = -1

var ArrayLengthPrefixed = func(_ *PacketData) int {
	return arrayLengthPrefix
}

func ArrayLengthStatic(length int) ArrayLengthOption {
	return func(_ *PacketData) int {
		return length
	}
}

func ArrayLengthDerive(fieldName string) ArrayLengthOption {
	return func(data *PacketData) int {
		return data.VarInt(fieldName)
	}
}

type ArrayValue = []PacketData

type ConvertedArrayValue = func(*PacketDefinition) ArrayValue

func ConvertArrayValue[T any](data []T, convertFunc func(T, *PacketData)) ConvertedArrayValue {
	return func(definition *PacketDefinition) ArrayValue {
		packets := make([]PacketData, len(data))

		for i, element := range data {
			packetData := definition.New()
			convertFunc(element, packetData)
			packets[i] = *packetData
		}

		return packets
	}
}

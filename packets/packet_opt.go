package packets

type ArrayLengthOption = int

const (
	ArrayLengthStatic = iota
	ArrayLengthAppend
)

type PacketOpt = func(*PacketDefinition)

func ID(id int) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.PacketID = id
	}
}

func Array(name string, lengthOption ArrayLengthOption, opts ...PacketOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeArray)
		packet.specifyArrayOptions(name, Packet(opts...), lengthOption)
	}
}

func Byte(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeByte)
	}
}

func Bool(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeBool)
	}
}

func Int16(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeInt16)
	}
}

func Int32(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeInt32)
	}
}

func Int64(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeInt64)
	}
}

func VarInt(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeVarInt)
	}
}

func Float32(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeFloat32)
	}
}

func Float64(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeFloat64)
	}
}

func UUIDField(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeUUID)
	}
}

func VarLong(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeVarLong)
	}
}

func ByteArray(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeByteArray)
	}
}

func String(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeString)
	}
}

func NBT(name string, blueprint any) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeNBT)
		packet.SetBlueprint(name, blueprint)
	}
}

func PositionField(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypePosition)
	}
}

func Slot(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeSlot)
	}
}

func BitSetField(name string) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeBitSet)
	}
}

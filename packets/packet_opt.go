package packets

type ArrayLengthOption = int

const (
	ArrayLengthStatic = iota
	ArrayLengthPrefixed
)

type PacketOpt = func(*PacketDefinition)

type PacketFieldOpt = func(data *PacketData) bool

func ID(id int) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.PacketID = id
	}
}

func Array(name string, lengthOption ArrayLengthOption, fields ...PacketOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeArray)
		packet.specifyArrayOptions(name, Packet(fields...), lengthOption)
	}
}
func ArrayWithOptions(name string, lengthOption ArrayLengthOption, fields []PacketOpt, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeArray)
		packet.specifyArrayOptions(name, Packet(fields...), lengthOption)
		packet.setFieldOpts(name, opts)
	}
}

func Byte(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeByte)
		packet.setFieldOpts(name, opts)
	}
}

func Bool(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeBool)
		packet.setFieldOpts(name, opts)
	}
}

func Int16(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeInt16)
		packet.setFieldOpts(name, opts)
	}
}

func Int32(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeInt32)
		packet.setFieldOpts(name, opts)
	}
}

func Int64(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeInt64)
		packet.setFieldOpts(name, opts)
	}
}

func VarInt(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeVarInt)
		packet.setFieldOpts(name, opts)
	}
}

func Float32(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeFloat32)
		packet.setFieldOpts(name, opts)
	}
}

func Float64(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeFloat64)
		packet.setFieldOpts(name, opts)
	}
}

func UUIDField(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeUUID)
		packet.setFieldOpts(name, opts)
	}
}

func VarLong(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeVarLong)
		packet.setFieldOpts(name, opts)
	}
}

func ByteArray(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeByteArray)
		packet.setFieldOpts(name, opts)
	}
}

func String(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeString)
		packet.setFieldOpts(name, opts)
	}
}

func NBT(name string, blueprint any, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeNBT)
		packet.SetBlueprint(name, blueprint)
		packet.setFieldOpts(name, opts)
	}
}

func PositionField(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypePosition)
		packet.setFieldOpts(name, opts)
	}
}

func Slot(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeSlot)
		packet.setFieldOpts(name, opts)
	}
}

func BitSetField(name string, opts ...PacketFieldOpt) PacketOpt {
	return func(packet *PacketDefinition) {
		packet.AddField(name, TypeBitSet)
		packet.setFieldOpts(name, opts)
	}
}

func OnlyIfTrue(fieldName string) PacketFieldOpt {
	return func(packet *PacketData) bool {
		return packet.Bool(fieldName)
	}
}

func OnlyIfFalse(fieldName string) PacketFieldOpt {
	return func(packet *PacketData) bool {
		return !packet.Bool(fieldName)
	}
}

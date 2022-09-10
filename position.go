package main

type Position struct {
	X int
	Y int
	Z int
}

func NewPosition(x, y, z int) *Position {
	return &Position{
		X: x,
		Y: y,
		Z: z,
	}
}

func PositionFromInt64(value int64) *Position {
	return &Position{
		X: int(value >> 38),
		Y: int(value << 52 >> 52),
		Z: int(value << 26 >> 38),
	}
}

func (p *Position) ToInt64() int64 {
	return ((int64(p.X) & 0x3FFFFFF) << 38) | ((int64(p.Z) & 0x3FFFFFF) << 12) | (int64(p.Y) & 0xFFF)
}

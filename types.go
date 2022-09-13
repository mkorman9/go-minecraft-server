package main

import "errors"

var IgnorablePacketReadError = errors.New("ignored")

type PacketMetadata struct {
	PacketSize           int
	UncompressedDataSize int
	UseCompression       bool
}

type UUID struct {
	Upper int64
	Lower int64
}

type BitSet struct {
	v []int64
}

func NewBitSet() *BitSet {
	return &BitSet{
		v: []int64{0},
	}
}

func (b *BitSet) Value(n int) bool {
	if n >= len(b.v)*64 {
		return false
	}

	return (b.v[n/64] & (1 << (n % 64))) != 0
}

func (b *BitSet) Set1(n int) {
	for len(b.v) <= n/64 {
		b.v = append(b.v, 0)
	}

	b.v[n/64] |= 1 << (n % 64)
}

func (b *BitSet) Set0(n int) {
	for len(b.v) <= n/64 {
		b.v = append(b.v, 0)
	}

	b.v[n/64] &= ^(1 << (n % 64))
}

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

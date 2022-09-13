package main

import "github.com/mkorman9/go-minecraft-server/nbt"

type ChunkData struct {
	Data []ChunkSection
}

type ChunkSection struct {
	BlockCount  int16
	BlockStates [4096]ChunkPalettedContainer
	Biomes      [64]ChunkPalettedContainer
}

type ChunkPalettedContainer struct {
	BitsPerEntry        byte
	PaletteSingleValued *ChunkPaletteSingleValued
	PaletteIndirect     *ChunkPaletteIndirect
	Data                []int64
}

type ChunkPaletteSingleValued struct {
	Value int
}

type ChunkPaletteIndirect struct {
	Palette []int
}

type BlockEntity struct {
	PackedXZ byte
	Y        int16
	Type     int
	Data     nbt.RawMessage
}

type Heightmap struct {
	MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
	WorldSurface   []int64 `nbt:"WORLD_SURFACE"`
}

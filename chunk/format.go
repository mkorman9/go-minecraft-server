package chunk

import (
	"github.com/mkorman9/go-minecraft-server/nbt"
	"github.com/mkorman9/go-minecraft-server/types"
	"io"
)

type Chunk struct {
	Sections []Section
}

type Section struct {
	BlockCount  int16
	BlockStates []PalettedContainer
	Biomes      []PalettedContainer
}

type PalettedContainer struct {
	BitsPerEntry        byte
	PaletteSingleValued *PaletteSingleValued
	PaletteIndirect     *PaletteIndirect
	Data                []int64
}

type PaletteSingleValued struct {
	Value int
}

type PaletteIndirect struct {
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

func GenerateExampleHeightmap() *Heightmap {
	return &Heightmap{
		MotionBlocking: make([]int64, 256),
		WorldSurface:   make([]int64, 256),
	}
}

func GenerateExampleChunk() *Chunk {
	sectionsCount := ChunkHeight / SectionHeight
	sections := make([]Section, sectionsCount)

	for i := 0; i < sectionsCount; i++ {
		if i <= sectionsCount/2 {
			sections[i] = getSolidSection()
		} else {
			sections[i] = getAirSection()
		}
	}

	return &Chunk{
		Sections: sections,
	}
}

func (c *Chunk) WriteTo(writer io.Writer) (int64, error) {
	sectionsCount := ChunkHeight / SectionHeight
	for section := 0; section < sectionsCount; section++ {
		err := types.WriteInt16(writer, c.Sections[section].BlockCount)
		if err != nil {
			return 0, err
		}

		for i := 0; i < BlocksPerSection; i++ {
			bitsPerEntry := c.Sections[section].BlockStates[i].BitsPerEntry

			err := types.WriteByte(writer, bitsPerEntry)
			if err != nil {
				return 0, err
			}

			switch {
			case bitsPerEntry == 0:
				err := types.WriteVarInt(writer, c.Sections[section].BlockStates[i].PaletteSingleValued.Value)
				if err != nil {
					return 0, err
				}
			default:
				// TODO
			}

			err = types.WriteVarInt(writer, len(c.Sections[section].BlockStates[i].Data))
			if err != nil {
				return 0, err
			}

			for dataElement := 0; dataElement < len(c.Sections[section].BlockStates[i].Data); dataElement++ {
				err = types.WriteInt64(writer, c.Sections[section].BlockStates[i].Data[dataElement])
				if err != nil {
					return 0, err
				}
			}
		}

		for i := 0; i < BiomesPerSection; i++ {
			bitsPerEntry := c.Sections[section].Biomes[i].BitsPerEntry

			err := types.WriteByte(writer, bitsPerEntry)
			if err != nil {
				return 0, err
			}

			switch {
			case bitsPerEntry == 0:
				err := types.WriteVarInt(writer, c.Sections[section].Biomes[i].PaletteSingleValued.Value)
				if err != nil {
					return 0, err
				}
			default:
				// TODO
			}

			err = types.WriteVarInt(writer, len(c.Sections[section].Biomes[i].Data))
			if err != nil {
				return 0, err
			}

			for dataElement := 0; dataElement < len(c.Sections[section].Biomes[i].Data); dataElement++ {
				err = types.WriteInt64(writer, c.Sections[section].Biomes[i].Data[dataElement])
				if err != nil {
					return 0, err
				}
			}
		}
	}

	return 0, nil
}

func getSolidSection() Section {
	section := Section{
		BlockCount:  BlocksPerSection,
		BlockStates: make([]PalettedContainer, BlocksPerSection),
		Biomes:      make([]PalettedContainer, BiomesPerSection),
	}

	for i := 0; i < BlocksPerSection; i++ {
		section.BlockStates[i] = PalettedContainer{
			BitsPerEntry: 0,
			PaletteSingleValued: &PaletteSingleValued{
				Value: 1,
			},
		}
	}

	for i := 0; i < BiomesPerSection; i++ {
		section.Biomes[i] = PalettedContainer{
			BitsPerEntry: 0,
			PaletteSingleValued: &PaletteSingleValued{
				Value: 1,
			},
		}
	}

	return section
}

func getAirSection() Section {
	section := Section{
		BlockCount:  0,
		BlockStates: make([]PalettedContainer, BlocksPerSection),
		Biomes:      make([]PalettedContainer, BiomesPerSection),
	}

	for i := 0; i < BlocksPerSection; i++ {
		section.BlockStates[i] = PalettedContainer{
			BitsPerEntry: 0,
			PaletteSingleValued: &PaletteSingleValued{
				Value: 0,
			},
		}
	}

	for i := 0; i < BiomesPerSection; i++ {
		section.Biomes[i] = PalettedContainer{
			BitsPerEntry: 0,
			PaletteSingleValued: &PaletteSingleValued{
				Value: 1,
			},
		}
	}

	return section
}

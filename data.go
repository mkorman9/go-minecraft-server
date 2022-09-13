package main

import (
	"encoding/json"
	"os"
)

type Data struct {
	DimensionCodec      *DimensionCodec
	SpawnPosition       *Position
	SpawnDimension      string
	WorldNames          []string
	IsHardcore          bool
	GameMode            GameMode
	HashedSeed          int64
	EnableRespawnScreen bool
	IsFlat              bool
}

func LoadData() (*Data, error) {
	data := Data{
		SpawnPosition:       NewPosition(0, 64, 0),
		SpawnDimension:      "minecraft:overworld",
		WorldNames:          []string{"minecraft:overworld", "minecraft:the_nether", "minecraft:the_end"},
		IsHardcore:          false,
		GameMode:            GameModeSurvival,
		HashedSeed:          0,
		EnableRespawnScreen: true,
		IsFlat:              true,
	}

	dimmensionCodecData, err := os.ReadFile("./data/1_19/dimension_codec.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dimmensionCodecData, &data.DimensionCodec)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

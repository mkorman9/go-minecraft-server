package main

import (
	"encoding/json"
	"os"
)

type Data struct {
	DimensionCodec *DimensionCodec
	SpawnPosition  *Position
}

func LoadData() (*Data, error) {
	data := Data{
		SpawnPosition: NewPosition(0, 64, 0),
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

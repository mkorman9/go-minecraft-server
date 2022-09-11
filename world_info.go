package main

type Registry[E any] struct {
	Type  string             `nbt:"type"`
	Value []RegistryValue[E] `nbt:"value"`
}

type RegistryValue[E any] struct {
	Name    string `nbt:"name"`
	ID      int32  `nbt:"id"`
	Element E      `nbt:"element"`
}

type RegistryCodec struct {
	Type  string             `nbt:"type"`
	Name  string             `nbt:"name"`
	Value RegistryCodecValue `nbt:"value"`
}

type RegistryCodecValue struct {
	DimensionType Registry[Dimension]     `nbt:"minecraft:dimension_type"`
	WorldGenBiome Registry[WorldGenBiome] `nbt:"minecraft:worldgen/biome"`
	//ChatType      Registry[nbt.RawMessage] `nbt:"minecraft:chat_type"`
}

type Dimension struct {
	HasSkylight        bool    `nbt:"has_skylight"`
	HasCeiling         bool    `nbt:"has_ceiling"`
	Ultrawarm          bool    `nbt:"ultrawarm"`
	Natural            bool    `nbt:"natural"`
	CoordinateScale    float64 `nbt:"coordinate_scale"`
	BedWorks           bool    `nbt:"bed_works"`
	RespawnAnchorWorks byte    `nbt:"respawn_anchor_works"`
	MinY               int32   `nbt:"min_y"`
	Height             int32   `nbt:"height"`
	LogicalHeight      int32   `nbt:"logical_height"`
	InfiniteBurn       string  `nbt:"infiniburn"`
	Effects            string  `nbt:"effects"`
	AmbientLight       float64 `nbt:"ambient_light"`
	PiglinSafe         byte    `nbt:"piglin_safe"`
	HasRaids           byte    `nbt:"has_raids"`
	//MonsterSpawnLightLevel      nbt.RawMessage `nbt:"monster_spawn_light_level"` // Tag_Int or {type:"minecraft:uniform", value:{min_inclusive: Tag_Int, max_inclusive: Tag_Int}}
	//MonsterSpawnBlockLightLimit int            `nbt:"monster_spawn_block_light_limit"`
}

type WorldGenBiome struct {
	Precipitation       string               `nbt:"precipitation"`
	Depth               float64              `nbt:"depth,omitempty"`
	Temperature         float64              `nbt:"temperature"`
	Scale               float64              `nbt:"scale,omitempty"`
	Downfall            float64              `nbt:"downfall"`
	Category            string               `nbt:"category,omitempty"`
	TemperatureModifier string               `nbt:"temperature_modifier,omitempty"`
	Effects             WorldGenBiomeEffects `nbt:"effects"`
}

type WorldGenBiomeEffects struct {
	SkyColor           int32  `nbt:"sky_color"`
	WaterFogColor      int32  `nbt:"water_fog_color"`
	FogColor           int32  `nbt:"fog_color"`
	WaterColor         int32  `nbt:"water_color"`
	FoliageColor       int32  `nbt:"foliage_color,omitempty"`
	GrassColor         int32  `nbt:"grass_color,omitempty"`
	GrassColorModifier string `nbt:"grass_color_modifier,omitempty"`
	// ...
}

func DefaultRegistryCodec() RegistryCodec {
	return RegistryCodec{
		Type: "compound",
		Name: "",
		Value: RegistryCodecValue{
			DimensionType: Registry[Dimension]{
				Type: "minecraft:dimension_type",
				Value: []RegistryValue[Dimension]{
					{
						Name: "minecraft:overworld",
						ID:   0,
						Element: Dimension{
							HasSkylight:        true,
							HasCeiling:         false,
							Ultrawarm:          false,
							Natural:            true,
							CoordinateScale:    1,
							BedWorks:           true,
							RespawnAnchorWorks: 0,
							MinY:               -64,
							Height:             384,
							LogicalHeight:      384,
							InfiniteBurn:       "minecraft:infiniburn_overworld",
							Effects:            "minecraft:overworld",
							AmbientLight:       0,
							PiglinSafe:         0,
							HasRaids:           1,
							//MonsterSpawnLightLevel:      nbt.RawMessage{Type: nbt.TagInt, Data: []byte{0, 0, 0, 8}},
							//MonsterSpawnBlockLightLimit: 8,
						},
					},
				},
			},
			WorldGenBiome: Registry[WorldGenBiome]{
				Type: "minecraft:worldgen/biome",
				Value: []RegistryValue[WorldGenBiome]{
					{
						Name: "minecraft:badlands",
						ID:   0,
						Element: WorldGenBiome{
							Precipitation: "none",
							Temperature:   2,
							Downfall:      0,
							Category:      "mesa",
							Effects: WorldGenBiomeEffects{
								SkyColor:      7254527,
								WaterFogColor: 329011,
								FogColor:      12638463,
								WaterColor:    4159204,
								FoliageColor:  10387789,
								GrassColor:    9470285,
							},
						},
					},
				},
			},
		},
	}
}

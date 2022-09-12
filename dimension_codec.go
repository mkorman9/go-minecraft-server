package main

import "github.com/mkorman9/go-minecraft-server/nbt"

type Registry[E any] struct {
	Type  string             `json:"type" nbt:"type"`
	Value []RegistryValue[E] `json:"value" nbt:"value"`
}

type RegistryValue[E any] struct {
	Name    string `json:"name" nbt:"name"`
	ID      int32  `json:"id" nbt:"id"`
	Element E      `json:"element" nbt:"element"`
}

type DimensionCodec struct {
	DimensionType Registry[Dimension]     `json:"dimension_type" nbt:"minecraft:dimension_type"`
	WorldGenBiome Registry[WorldGenBiome] `json:"world_gen_biome" nbt:"minecraft:worldgen/biome"`
	ChatType      Registry[ChatType]      `json:"chat_type" nbt:"minecraft:chat_type"`
}

type Dimension struct {
	PiglinSafe                  byte           `json:"piglin_safe" nbt:"piglin_safe"`
	Natural                     bool           `json:"natural" nbt:"natural"`
	AmbientLight                float64        `json:"ambient_light" nbt:"ambient_light"`
	FixedTime                   int64          `json:"fixed_time,omitempty" nbt:"fixed_time,omitempty"`
	InfiniteBurn                string         `json:"infiniburn" nbt:"infiniburn"`
	RespawnAnchorWorks          byte           `json:"respawn_anchor_works" nbt:"respawn_anchor_works"`
	HasSkylight                 bool           `json:"has_skylight" nbt:"has_skylight"`
	BedWorks                    bool           `json:"bed_works" nbt:"bed_works"`
	Effects                     string         `json:"effects" nbt:"effects"`
	HasRaids                    byte           `json:"has_raids" nbt:"has_raids"`
	LogicalHeight               int32          `json:"logical_height" nbt:"logical_height"`
	CoordinateScale             float64        `json:"coordinate_scale" nbt:"coordinate_scale"`
	MinY                        int32          `json:"min_y" nbt:"min_y"`
	HasCeiling                  bool           `json:"has_ceiling" nbt:"has_ceiling"`
	Ultrawarm                   bool           `json:"ultrawarm" nbt:"ultrawarm"`
	Height                      int32          `json:"height" nbt:"height"`
	MonsterSpawnLightLevel      nbt.RawMessage `json:"monster_spawn_light_level" nbt:"monster_spawn_light_level"`
	MonsterSpawnBlockLightLimit int32          `json:"monster_spawn_block_light_limit" nbt:"monster_spawn_block_light_limit"`
}

type WorldGenBiome struct {
	Precipitation       string               `json:"precipitation" nbt:"precipitation"`
	Effects             WorldGenBiomeEffects `json:"effects" nbt:"effects"`
	Temperature         float64              `json:"temperature" nbt:"temperature"`
	Downfall            float64              `json:"downfall" nbt:"downfall"`
	Category            string               `json:"category,omitempty" nbt:"category,omitempty"`
	Depth               float64              `json:"depth,omitempty" nbt:"depth,omitempty"`
	Scale               float64              `json:"scale,omitempty" nbt:"scale,omitempty"`
	TemperatureModifier string               `json:"temperature_modifier,omitempty" nbt:"temperature_modifier,omitempty"`
}

type WorldGenBiomeEffects struct {
	SkyColor           int32  `json:"sky_color" nbt:"sky_color"`
	WaterFogColor      int32  `json:"water_fog_color" nbt:"water_fog_color"`
	FogColor           int32  `json:"fog_color" nbt:"fog_color"`
	WaterColor         int32  `json:"water_color" nbt:"water_color"`
	FoliageColor       int32  `json:"omitempty" nbt:"foliage_color,omitempty"`
	GrassColor         int32  `json:"grass_color,omitempty" nbt:"grass_color,omitempty"`
	GrassColorModifier string `json:"grass_color_modifier,omitempty" nbt:"grass_color_modifier,omitempty"`
	// ...
}

type ChatType struct {
	Chat      ChatProperties `json:"chat" nbt:"chat"`
	Narration ChatNarration  `json:"narration" nbt:"narration"`
}

type ChatProperties struct {
	Decoration ChatDecoration `json:"decoration" nbt:"decoration"`
}

type ChatNarration struct {
	Decoration ChatDecoration `json:"decoration" nbt:"decoration"`
	Priority   string         `json:"priority" nbt:"priority"`
}

type ChatDecoration struct {
	Parameters     []string  `json:"parameters" nbt:"parameters"`
	TranslationKey string    `json:"translation_key" nbt:"translation_key"`
	Style          ChatStyle `json:"style" nbt:"style"`
}

type ChatStyle struct {
}

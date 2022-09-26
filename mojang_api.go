package main

import (
	"encoding/json"
	"fmt"
	"github.com/mkorman9/go-minecraft-server/types"
	"net/http"
	"time"
)

type MojangPlayerVerification struct {
	Verified          bool
	UUID              *types.UUID
	Textures          string
	TexturesSignature string
}

type mojangVerifyPlayerResponse struct {
	ID         string                       `json:"id"`
	Name       string                       `json:"name"`
	Properties []mojangVerifyPlayerProperty `json:"properties"`
}

type mojangVerifyPlayerProperty struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature,omitempty"`
}

func MojangVerifyPlayer(username, hash string) (*MojangPlayerVerification, error) {
	time.Sleep(5 * time.Second)

	url := fmt.Sprintf(
		"https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s",
		username,
		hash,
	)

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return &MojangPlayerVerification{Verified: false}, nil
	}

	var mojangResponse mojangVerifyPlayerResponse
	err = json.NewDecoder(response.Body).Decode(&mojangResponse)
	if err != nil {
		return nil, err
	}

	if username != mojangResponse.Name {
		return &MojangPlayerVerification{Verified: false}, nil
	}

	uuid, err := mojangIdToUUID(mojangResponse.ID)
	if err != nil {
		return nil, err
	}

	var textures string
	var texturesSignature string
	for _, property := range mojangResponse.Properties {
		if property.Name == "textures" {
			textures = property.Value
			texturesSignature = property.Signature
		}
	}

	return &MojangPlayerVerification{
		Verified:          true,
		UUID:              uuid,
		Textures:          textures,
		TexturesSignature: texturesSignature,
	}, nil
}

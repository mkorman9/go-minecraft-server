package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MojangPlayerVerification struct {
	Verified bool
	ID       string
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

	// introduce retry on 204
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

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

	return &MojangPlayerVerification{Verified: true, ID: mojangResponse.ID}, nil
}
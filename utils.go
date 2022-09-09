package main

import (
	"crypto/rand"
	"encoding/hex"
)

func getSecureRandomString(lengthBytes int) (string, error) {
	var bytes = make([]byte, lengthBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

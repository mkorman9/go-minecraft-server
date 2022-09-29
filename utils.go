package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"github.com/mkorman9/go-minecraft-server/types"
	"net"
	"strings"
)

func getSecureRandomString(lengthBytes int) (string, error) {
	var buff = make([]byte, lengthBytes)
	if _, err := rand.Read(buff); err != nil {
		return "", err
	}

	return hex.EncodeToString(buff), nil
}

func loadPublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

func verifyRsaSignature(publicKey *rsa.PublicKey, msg string, salt int64, signature []byte) error {
	saltEncoded := make([]byte, 8)
	binary.BigEndian.PutUint64(saltEncoded, uint64(salt))

	msgToHash := []byte(msg)
	msgToHash = append(msgToHash, saltEncoded...)

	hash := sha256.Sum256(msgToHash)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
}

func parseRemoteAddress(connection net.Conn) string {
	address := connection.RemoteAddr().String()
	portIndex := strings.LastIndex(address, ":")
	if portIndex >= 0 {
		address = address[:portIndex]
	}

	return address
}

func mojangIdToUUID(mojangId string) (*types.UUID, error) {
	upperPart := mojangId[:16]
	lowerPart := mojangId[16:]

	upperPartBytes, err := hex.DecodeString(upperPart)
	if err != nil {
		return nil, err
	}

	lowerPartBytes, err := hex.DecodeString(lowerPart)
	if err != nil {
		return nil, err
	}

	return &types.UUID{
		Upper: int64(binary.BigEndian.Uint64(upperPartBytes)),
		Lower: int64(binary.BigEndian.Uint64(lowerPartBytes)),
	}, nil
}

func packBlockXZ(x, z byte) byte {
	return ((x & 15) << 4) | (z & 15)
}

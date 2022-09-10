package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"github.com/gofrs/uuid"
)

func getSecureRandomString(lengthBytes int) (string, error) {
	var bytes = make([]byte, lengthBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func getRandomUUID() UUID {
	v, _ := uuid.NewV4()
	upper := int64(binary.BigEndian.Uint64([]byte{v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7]}))
	lower := int64(binary.BigEndian.Uint64([]byte{v[8], v[9], v[10], v[11], v[12], v[13], v[14], v[15]}))

	return UUID{
		Upper: upper,
		Lower: lower,
	}
}

func loadPublicKey(publicKey string) (*rsa.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

func verifyRsaSignature(publicKey *rsa.PublicKey, msg string, salt int64, signature string) error {
	saltEncoded := make([]byte, 8)
	binary.BigEndian.PutUint64(saltEncoded, uint64(salt))

	msgToHash := []byte(msg)
	msgToHash = append(msgToHash, saltEncoded...)

	hash := sha256.Sum256(msgToHash)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], []byte(signature))
}

package main

import (
	"crypto/aes"
	"crypto/cipher"
)

type CipherStream struct {
	block     cipher.Block
	encrypter cipher.Stream
	decrypter cipher.Stream
}

func NewCipherStream(key string) (*CipherStream, error) {
	//keyBytes := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	encrypter := cipher.NewCFBEncrypter(block, []byte(key))
	decrypter := cipher.NewCFBDecrypter(block, []byte(key))

	return &CipherStream{
		block:     block,
		encrypter: encrypter,
		decrypter: decrypter,
	}, nil
}

func (cs *CipherStream) Encrypt(plainText []byte) []byte {
	cipherText := make([]byte, len(plainText))
	cs.encrypter.XORKeyStream(cipherText, plainText)
	return cipherText
}

func (cs *CipherStream) Decrypt(cipherText []byte) []byte {
	plainText := make([]byte, len(cipherText))
	cs.encrypter.XORKeyStream(plainText, cipherText)
	return plainText
}

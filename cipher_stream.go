package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
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

func (cs *CipherStream) WrapReader(reader io.Reader) io.Reader {
	return &cipher.StreamReader{
		S: cs.decrypter,
		R: reader,
	}
}

func (cs *CipherStream) WrapWriter(writer io.Writer) io.Writer {
	return &cipher.StreamWriter{
		S: cs.encrypter,
		W: writer,
	}
}

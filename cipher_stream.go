package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type CipherStream struct {
	block     cipher.Block
	encrypter *CFB8
	decrypter *CFB8
}

type CFB8 struct {
	block       cipher.Block
	key         []byte
	tmp         []byte
	encryptMode bool
}

func NewCipherStream(key []byte) (*CipherStream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)
	encrypter := &CFB8{block: block, key: keyCopy, tmp: make([]byte, block.BlockSize()), encryptMode: true}

	keyCopy = make([]byte, len(key))
	copy(keyCopy, key)
	decrypter := &CFB8{block: block, key: keyCopy, tmp: make([]byte, block.BlockSize()), encryptMode: false}

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

func (c *CFB8) XORKeyStream(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		val := src[i]
		copy(c.tmp, c.key)
		c.block.Encrypt(c.key, c.key)

		val = val ^ c.key[0]

		copy(c.key, c.tmp[1:])
		if c.encryptMode {
			c.key[15] = val
		} else {
			c.key[15] = src[i]
		}

		dst[i] = val
	}
}

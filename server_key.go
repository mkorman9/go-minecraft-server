package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type ServerKey struct {
	private *rsa.PrivateKey
	public  crypto.PublicKey

	publicASN1 string
}

func GenerateServerKey() (*ServerKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, ServerKeyLength)
	if err != nil {
		return nil, err
	}

	public := key.Public()

	publicASN1, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return nil, err
	}

	return &ServerKey{
		private:    key,
		public:     &public,
		publicASN1: string(publicASN1),
	}, nil
}

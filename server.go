package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"net"
	"strings"
)

type Server struct {
	key      *serverKey
	listener net.Listener
}

func NewServer(settings *Settings) (*Server, error) {
	key, err := generateServerKey()
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", settings.ServerAddress)
	if err != nil {
		return nil, err
	}

	return &Server{
		key:      key,
		listener: listener,
	}, nil
}
func (s *Server) PublicKey() []byte {
	return s.key.publicDER
}

func (s *Server) Shutdown() {
	_ = s.listener.Close()
}

func (s *Server) AcceptLoop(handleConnection func(conn net.Conn, ip string)) error {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			if netOpError, ok := err.(*net.OpError); ok {
				if netOpError.Err.Error() == "use of closed network connection" {
					break
				}
			}

			return err
		}

		ip, _, _ := strings.Cut(connection.RemoteAddr().String(), ":")

		go handleConnection(connection, ip)
	}

	return nil
}

func (s *Server) DecryptMessage(message []byte) ([]byte, error) {
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, s.key.private, message)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func (s *Server) GenerateServerHash(sharedSecret []byte) string {
	hash := sha1.New()
	hash.Write(sharedSecret)
	hash.Write(s.key.publicDER)
	return hex.EncodeToString(hash.Sum(nil))
}

type serverKey struct {
	private   *rsa.PrivateKey
	public    crypto.PublicKey
	publicDER []byte
}

func generateServerKey() (*serverKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, ServerKeyLength)
	if err != nil {
		return nil, err
	}

	public := key.Public()

	publicASN1, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return nil, err
	}

	return &serverKey{
		private:   key,
		public:    &public,
		publicDER: publicASN1,
	}, nil
}

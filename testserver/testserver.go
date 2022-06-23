package testserver

import (
	"crypto/tls"
	"github.com/emersion/go-imap/backend/memory"
	"github.com/emersion/go-imap/server"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

const Address = "127.0.0.1:0"

func New(t *testing.T) (s *server.Server, l net.Listener) {
	l = createListener(t)
	return createServer(t, l), l
}

func NewTls(t *testing.T) (s *server.Server, l net.Listener) {
	l = createTlsListener(t)
	return createServer(t, l), l
}

func createServer(t *testing.T, listener net.Listener) *server.Server {
	s := server.New(memory.New())
	s.AllowInsecureAuth = true

	go func(l net.Listener) {
		assert.NoError(t, s.Serve(l))
	}(listener)

	return s
}

func createTlsListener(t *testing.T) net.Listener {
	cert, err := tls.LoadX509KeyPair("../testdata/cert.pem", "../testdata/key.pem")
	assert.NoError(t, err)

	listener, err := tls.Listen("tcp", Address, &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	})
	assert.NoError(t, err)
	return listener
}

func createListener(t *testing.T) net.Listener {
	listener, err := net.Listen("tcp", Address)
	assert.NoError(t, err)
	return listener
}

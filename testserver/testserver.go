package testserver

import (
	"github.com/emersion/go-imap/backend/memory"
	"github.com/emersion/go-imap/server"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

const Address = "127.0.0.1:0"

func New(t *testing.T) (s *server.Server, l net.Listener) {
	l = createListener(t)
	s = server.New(memory.New())
	s.AllowInsecureAuth = true

	go func(l net.Listener) {
		assert.NoError(t, s.Serve(l))
	}(l)

	return s, l
}

func createListener(t *testing.T) net.Listener {
	listener, err := net.Listen("tcp", Address)
	assert.NoError(t, err)
	return listener
}

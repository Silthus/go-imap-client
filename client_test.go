package main

import (
	"fmt"
	"github.com/emersion/go-imap/backend/memory"
	"github.com/emersion/go-imap/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"net"
	"testing"
)

const testServerAddress = "127.0.0.1:0"

type ClientTestSuite struct {
	suite.Suite
	client Client

	imapServer *server.Server
	listener   net.Listener
}

func (s *ClientTestSuite) SetupTest() {
	s.listener = createListener(s.T())
	s.imapServer = createServer(s.T(), s.listener)

	s.client = create(s.listener.Addr().String())
}

func createListener(t *testing.T) net.Listener {
	listener, err := net.Listen("tcp", testServerAddress)
	assert.NoError(t, err)
	return listener
}

func createServer(t *testing.T, listener net.Listener) *server.Server {
	imapServer := server.New(memory.New())
	imapServer.AllowInsecureAuth = true

	go func(l net.Listener) {
		assert.NoError(t, imapServer.Serve(l))
	}(listener)

	return imapServer
}

func (s *ClientTestSuite) TearDownTest() {
	s.NoError(s.imapServer.Close())
}

func (s *ClientTestSuite) TestConnect() {
	s.Run("connect with valid username and password", func() {
		err := s.client.Connect("username", "password")
		s.NoError(err)
		s.True(s.client.IsConnected())
	})
	//s.Run("connect with wrong password fails", func() {
	//	client, err := Connect(testServerAddress, "username", "wrong")
	//	s.Error(err)
	//	s.NotNil(client)
	//	s.False(client.IsConnected())
	//})
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func ExampleConnect() {
	client, err := Connect("imap.my-server.local", "username", "password")
	if err != nil {
		log.Fatalf("error while connecting to mail server %s: %s", testServerAddress, err)
	}
	fmt.Printf("IMAP client is connected to %s: %t", client.ServerAddress(), client.IsConnected())
}

package main

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

const testServerAddress = "localhost"

type ClientTestSuite struct {
	suite.Suite
}

func (s *ClientTestSuite) SetupTest() {
}

func (s *ClientTestSuite) TestNew() {
	s.Run("new client contains server address", func() {
		client := New(testServerAddress)
		s.Equal(testServerAddress, client.ServerAddress())
	})
}

func (s *ClientTestSuite) TestConnect() {
	s.Run("connect with valid username and password", func() {
		client, err := Connect(testServerAddress, "username", "password")
		s.NoError(err)
		s.True(client.IsConnected())
	})
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func ExampleConnect() {
	client, err := Connect(testServerAddress, "username", "password")
	if err != nil {
		log.Fatalf("error while connecting to mail server %s: %s", testServerAddress, err)
	}
	fmt.Printf("IMAP client is connected to %s: %t", client.ServerAddress(), client.IsConnected())
	// Output: IMAP client is connected to localhost: true
}

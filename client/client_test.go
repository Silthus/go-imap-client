package client

import (
	"fmt"
	"github.com/Silthus/go-imap-client/testserver"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/server"
	"github.com/stretchr/testify/suite"
	"log"
	"net"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
	client Client

	imapServer *server.Server
	listener   net.Listener
}

func (s *ClientTestSuite) SetupSuite() {
	s.imapServer, s.listener = testserver.New(s.T())
}

func (s *ClientTestSuite) SetupTest() {
	s.client = create(s.listener.Addr().String())
}

func (s *ClientTestSuite) BeforeTest(_, testName string) {
	if testName == "TestConnect" {
		return
	}
	s.NoError(s.connect())
}

func (s *ClientTestSuite) AfterTest(string, string) {
}

func (s *ClientTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *ClientTestSuite) TearDownSuite() {
	s.NoError(s.imapServer.Close())
}

func (s *ClientTestSuite) connect() error {
	return s.client.Connect("username", "password")
}

func (s *ClientTestSuite) TestConnect() {
	s.Run("connect with valid username and password", func() {
		s.NoError(s.connect())
		s.True(s.client.IsConnected())
		s.Equal("username", s.client.AuthenticatedUser())
		s.client.Close()
	})
	s.Run("connect with wrong password fails", func() {
		client, err := Connect(testserver.Address, "username", "wrong")
		s.Error(err)
		s.NotNil(client)
		s.False(client.IsConnected())
	})
	s.Run("connect when connected throws already connected error", func() {
		s.NoError(s.connect())
		s.ErrorIs(s.connect(), AlreadyConnectedError)
		s.client.Close()
	})
}

func (s *ClientTestSuite) TestClose() {
	s.Run("IsConnected() is false after Close()", func() {
		s.client.Close()
		s.False(s.client.IsConnected())
	})
	s.Run("AuthenticatedUser() is zero value after Close()", func() {
		s.client.Close()
		s.Zero(s.client.AuthenticatedUser())
	})
}

func (s *ClientTestSuite) TestSearchMailbox() {
	tests := []struct {
		name       string
		mailbox    string
		searchTerm string
		assertion  func(messages []*imap.Message, err error)
	}{
		{"search unknown mailbox throws error", "Unknown", "search term", s.assertOpenMailboxError},
		{"search mailbox with unknown search term returns empty slice", InboxName, "unknown", s.assertEmptySearch},
		// available data in memory backend server: https://github.com/emersion/go-imap/blob/master/backend/memory/backend.go
		{"search term matches subject", InboxName, "just for you", s.assertMessageCount(1)},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			test.assertion(s.client.SearchMailbox(test.mailbox, test.searchTerm))
		})
	}
}

func (s *ClientTestSuite) assertOpenMailboxError(messages []*imap.Message, err error) {
	s.T().Helper()
	s.ErrorAs(err, &OpenMailboxError{})
	s.Empty(messages)
}

func (s *ClientTestSuite) assertEmptySearch(messages []*imap.Message, err error) {
	s.T().Helper()
	s.NoError(err)
	s.Len(messages, 0)
}

func (s *ClientTestSuite) assertMessageCount(count int) func(messages []*imap.Message, err error) {
	return func(messages []*imap.Message, err error) {
		s.NoError(err)
		s.Len(messages, count)
	}
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func ExampleConnect() {
	client, err := Connect("imap.my-server.local", "username", "password")
	if err != nil {
		log.Fatalf("error while connecting to mail server %s: %s", testserver.Address, err)
	}
	fmt.Printf("IMAP client is connected to %s: %t", client.ServerAddress(), client.IsConnected())
}

func ExampleClient_SearchMailbox() {
	// connect to the imap server and login
	client, err := Connect("imap.my-server.local", "username", "password")
	if err != nil {
		log.Fatalf("error while connecting to mail server %s: %s", testserver.Address, err)
	}
	// search the inbox for messages with "my message" in their subject
	messages, err := client.SearchMailbox(InboxName, "my message")
	if err != nil {
		log.Fatalf("error while searching mailbox %q: %s", InboxName, err)
	}
	for _, msg := range messages {
		log.Printf("Found message: %s", msg.Envelope.Subject)
	}
}

package main

import (
	"github.com/emersion/go-imap/client"
)

// Client wraps the github.com/emersion/go-imap library into a more easily usable client
// to search and read mails in a mailbox.
//
// Use Connect to create a create client and connect to the mailbox server.
type Client interface {
	// ServerAddress the Client targets.
	ServerAddress() string
	// IsConnected is true if the Client is currently connected to the server.
	//
	// Close must be called to clean up after usage.
	IsConnected() bool
	// Connect uses the given username and password to connect to the server
	// address of this client.
	//
	// Close must be called to clean up the open connection after usage.
	Connect(username, password string) error
	// Close the Client connection to the IMAP server.
	Close()
}

// new creates an IMAP client targeting the given mail server.
//
// Connect can be used to directly connect to the server with the default configuration.
func create(server string) Client {
	return &clientImpl{serverAddress: server}
}

// Connect creates a new Client and directly connects to the given mail server using the provided username and password.
func Connect(server string, username string, password string) (client Client, err error) {
	client = create(server)
	return client, client.Connect(username, password)
}

type clientImpl struct {
	serverAddress string

	imapClient *client.Client
}

func (c *clientImpl) ServerAddress() string {
	return c.serverAddress
}

func (c *clientImpl) IsConnected() bool {
	return true
}

func (c *clientImpl) Connect(username, password string) (err error) {
	c.imapClient, err = client.Dial(c.serverAddress)
	if err != nil {
		return err
	}

	err = c.imapClient.Login(username, password)
	return err
}

func (c *clientImpl) Close() {
	//TODO implement me
	panic("implement me")
}

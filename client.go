package main

// Client wraps the github.com/emersion/go-imap library into a more easily usable client
// to search and read mails in a mailbox.
//
// Use Connect or New to create a new client and connect to the mailbox server.
type Client interface {
	// ServerAddress the Client targets.
	ServerAddress() string
	// IsConnected is true if the Client is currently connected to the server.
	//
	// Close must be called to clean up after usage.
	IsConnected() bool
	// Connect with the given username and password against the ServerAddress of this client.
	//
	// Close must be called to clean up the open connection after usage.
	Connect(username, password string) error
	// Close the Client connection to the IMAP server.
	Close()
}

// New creates an IMAP client targeting the given mail server.
//
// Connect can be used to directly connect to the server with the default configuration.
func New(server string) Client {
	return &clientImpl{serverAddress: server}
}

// Connect creates a new Client and directly connects to the given mail server using the provided username and password.
//
// Use New to customize the configuration of the Client.
func Connect(server string, username string, password string) (client Client, err error) {
	return New(server), nil
}

type clientImpl struct {
	serverAddress string
}

func (c *clientImpl) ServerAddress() string {
	return c.serverAddress
}

func (c *clientImpl) IsConnected() bool {
	return true
}

func (c *clientImpl) Connect(username, password string) error {
	//TODO implement me
	panic("implement me")
}

func (c *clientImpl) Close() {
	//TODO implement me
	panic("implement me")
}

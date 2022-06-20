package main

import (
	"errors"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"strings"
)

const InboxName = imap.InboxName

// Client wraps the github.com/emersion/go-imap library into a more easily usable client
// to search and read mails in a mailbox.
//
// Use Connect to create a create client and connect to the mailbox server.
type Client interface {
	// ServerAddress the Client targets.
	ServerAddress() string
	// AuthenticatedUser that is currently logged in to the server.
	// nil if authentication failed or not yet connected.
	AuthenticatedUser() string
	// IsConnected is true if the Client is currently connected to the server.
	//
	// Close must be called to clean up after usage.
	IsConnected() bool
	// Connect uses the given username and password to connect to the server
	// address of this client.
	//
	// Close must be called to clean up the open connection after usage.
	Connect(username, password string) error
	// SearchMailbox searches the given mailbox for messages that contain the searchTerm in their subject.
	SearchMailbox(mailbox string, searchTerm string) (messages []*imap.Message, err error)
	// Close the Client connection to the IMAP server.
	Close()
}

type OpenMailboxError struct {
	error
	User    string
	Server  string
	Mailbox string
	Cause   error
}

func (e OpenMailboxError) Error() string {
	return fmt.Sprintf("cannot open mailbox %q of %q on %q: %s", e.Mailbox, e.User, e.Server, e.Cause)
}

var AlreadyConnectedError = errors.New("already connected to the mail server")

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

	imapClient        *client.Client
	authenticatedUser string
}

func (c *clientImpl) ServerAddress() string {
	return c.serverAddress
}

func (c *clientImpl) AuthenticatedUser() string {
	return c.authenticatedUser
}

func (c *clientImpl) IsConnected() bool {
	return c.imapClient != nil
}

func (c *clientImpl) Connect(username, password string) (err error) {
	if c.IsConnected() {
		return AlreadyConnectedError
	}
	c.imapClient, err = client.Dial(c.serverAddress)
	if err != nil {
		return fmt.Errorf("unable to connect to the mail server: %w", err)
	}

	err = c.imapClient.Login(username, password)
	if err != nil {
		return fmt.Errorf("unable to login into the mail server: %w", err)
	}

	c.authenticatedUser = username
	return nil
}

func (c *clientImpl) SearchMailbox(mailboxName, searchTerm string) (messages []*imap.Message, err error) {
	mailbox, err := c.imapClient.Select(mailboxName, true)
	if err != nil {
		return nil, c.throwOpenMailboxError(mailboxName, err)
	}

	messageChannel, err := c.fetchMessages(mailbox)
	if err != nil {
		return nil, err
	}

	return c.filterMessages(messageChannel, searchTerm), nil
}

func (c *clientImpl) filterMessages(messageChannel chan *imap.Message, searchTerm string) (messages []*imap.Message) {
	for msg := range messageChannel {
		if c.subjectContains(msg, searchTerm) {
			messages = append(messages, msg)
		}
	}
	return messages
}

func (c *clientImpl) subjectContains(msg *imap.Message, searchTerm string) bool {
	return strings.Contains(msg.Envelope.Subject, searchTerm)
}

func (c *clientImpl) fetchMessages(mailbox *imap.MailboxStatus) (messages chan *imap.Message, err error) {
	seqset := new(imap.SeqSet)
	seqset.AddRange(mailbox.Messages, mailbox.Messages)
	messages = make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	if err = <-done; err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *clientImpl) Close() {
	if c.IsConnected() {
		c.imapClient.Close()
	}
	c.authenticatedUser = ""
	c.imapClient = nil
}

func (c *clientImpl) throwOpenMailboxError(mailbox string, cause error) OpenMailboxError {
	return OpenMailboxError{Server: c.ServerAddress(), User: c.AuthenticatedUser(), Mailbox: mailbox, Cause: cause}
}

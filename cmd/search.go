/*
Copyright © 2022 Michael Reichenbach <me@silthus.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/emersion/go-imap"
	imapClient "github.com/emersion/go-imap/client"
	"github.com/spf13/cobra"
	"net"
	"strings"
	"time"
)

var (
	mailbox        string
	noResultsError bool
)

func newSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <search term>",
		Short: "Searches a mailbox for mails matching a subject",
		Long: `Allows searching a mailbox for mails where their subject matches the given search term. For example:

imap --server "my-server:993" --username username --password password --tls search my mail`,
		Args: cobra.MinimumNArgs(1),
		RunE: searchMailbox,
	}

	cmd.Flags().StringVarP(&mailbox, "mailbox", "m", imap.InboxName, "name of the mailbox")
	cmd.Flags().BoolVarP(&noResultsError, "no-results-error", "e", false, "set to exit with an error code if no mails are found")

	return cmd
}

func searchMailbox(cmd *cobra.Command, args []string) error {
	client, err := connectAndLogin()
	if err != nil {
		return err
	}
	mailbox, err := client.Select(imap.InboxName, true)
	if err != nil {
		return err
	}

	subject := args[0]
	messagesChan, err := fetchAndFilterMessages(client, mailbox, subject)
	if err != nil {
		return err
	}

	return printMessagesOrRaiseError(cmd, subject, collectMessages(messagesChan))
}

func printMessagesOrRaiseError(cmd *cobra.Command, subject string, messages []*imap.Message) error {
	errorMessage := fmt.Sprintf("Found no messages matching the search term: %q", subject)

	if noResultsError && len(messages) < 1 {
		return errors.New(errorMessage)
	} else if len(messages) < 1 {
		cmd.Println(errorMessage)
	} else {
		printResults(cmd, messages)
	}

	return nil
}

func collectMessages(messagesChan <-chan *imap.Message) []*imap.Message {
	messages := make([]*imap.Message, 0)
	for msg := range messagesChan {
		messages = append(messages, msg)
	}
	return messages
}

func connectAndLogin() (client *imapClient.Client, err error) {
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	if useTls {
		client, err = imapClient.DialWithDialerTLS(dialer, server, &tls.Config{
			InsecureSkipVerify: skipVerifyTls,
		})
	} else {
		client, err = imapClient.DialWithDialer(dialer, server)
	}
	if err != nil {
		return nil, err
	}

	err = client.Login(username, password)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func fetchAndFilterMessages(client *imapClient.Client, mailbox *imap.MailboxStatus, subject string) (<-chan *imap.Message, error) {
	messages, done := fetchMessages(client, mailbox)
	messages = filterMessages(messages, subject)
	if err := <-done; err != nil {
		return messages, err
	}
	return messages, nil
}

func printResults(cmd *cobra.Command, messages []*imap.Message) {
	for _, msg := range messages {
		cmd.Println(msg.Envelope.Subject)
	}
}

func fetchMessages(client *imapClient.Client, mailbox *imap.MailboxStatus) (<-chan *imap.Message, <-chan error) {
	seqset := new(imap.SeqSet)
	seqset.AddRange(mailbox.Messages, mailbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- client.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	return messages, done
}

func filterMessages(messages <-chan *imap.Message, subject string) <-chan *imap.Message {
	out := make(chan *imap.Message)
	go func() {
		for msg := range messages {
			if strings.Contains(msg.Envelope.Subject, subject) {
				out <- msg
			}
		}
		close(out)
	}()
	return out
}

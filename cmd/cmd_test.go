package cmd

import (
	"bytes"
	"fmt"
	"github.com/Silthus/go-imap-client/testserver"
	imapSrv "github.com/emersion/go-imap/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

type CmdTestSuite struct {
	suite.Suite
	cmd               *cobra.Command
	b                 *bytes.Buffer
	testServerAddress string
	testServer        *imapSrv.Server
}

func (s *CmdTestSuite) SetupTest() {
	s.setupCommand()
	s.setupTestServer()
}

func (s *CmdTestSuite) setupCommand() {
	s.cmd = newRootCmd()
	s.b = bytes.NewBufferString("")
	s.cmd.Println()
	s.cmd.SetOut(s.b)
}

func (s *CmdTestSuite) setupTestServer() {
	srv, listener := testserver.New(s.T())
	s.testServerAddress = listener.Addr().String()
	s.testServer = srv
}

func (s *CmdTestSuite) executeErr(args ...string) (string, error) {
	s.T().Helper()
	s.cmd.SetArgs(args)
	execErr := s.cmd.Execute()
	out, err := ioutil.ReadAll(s.b)
	s.NoError(err)
	return string(out), execErr
}

func (s *CmdTestSuite) execute(args ...string) string {
	s.T().Helper()
	out, _ := s.executeErr(args...)
	return out
}

func (s *CmdTestSuite) assertExecuteError(expectedError string, args ...string) {
	s.T().Helper()
	_, err := s.executeErr(args...)
	s.EqualError(err, expectedError)
}

func (s *CmdTestSuite) assertOptionalFlag(flag string, value, expected interface{}) {
	s.T().Helper()
	var additionalFlags []string
	if value != "" {
		additionalFlags = append(additionalFlags, "--"+fmt.Sprintf("%v", flag)+"="+fmt.Sprintf("%v", value))
	}
	s.T().Helper()
	args := append([]string{"search", "--server=" + s.testServerAddress, "--username=username", "--password=password", "any"}, additionalFlags...)
	_, err := s.executeErr(args...)
	s.NoError(err)
	f := s.cmd.Flags().Lookup(flag)
	s.Equal(expected, f.Value.String())
}

func TestCmdTestSuite(t *testing.T) {
	suite.Run(t, new(CmdTestSuite))
}

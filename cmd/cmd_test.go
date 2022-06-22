package cmd

import (
	"bytes"
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
	s.cmd.SetArgs(args)
	execErr := s.cmd.Execute()
	out, err := ioutil.ReadAll(s.b)
	s.NoError(err)
	return string(out), execErr
}

func (s *CmdTestSuite) execute(args ...string) string {
	out, _ := s.executeErr(args...)
	return out
}

func (s *CmdTestSuite) assertExecuteError(expectedError string, args ...string) {
	_, err := s.executeErr(args...)
	s.EqualError(err, expectedError)
}

func (s *CmdTestSuite) assertRequiredFlag(flag string, args ...string) {
	_, err := s.executeErr(args...)
	s.ErrorContains(err, "required flag(s)")
	s.ErrorContains(err, flag)
}

func TestCmdTestSuite(t *testing.T) {
	suite.Run(t, new(CmdTestSuite))
}

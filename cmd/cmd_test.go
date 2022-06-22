package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

type CmdTestSuite struct {
	suite.Suite
	cmd *cobra.Command
	b   *bytes.Buffer
}

func (s *CmdTestSuite) SetupTest() {
	s.cmd = newRootCmd()
	s.b = bytes.NewBufferString("")
	s.cmd.Println()
	s.cmd.SetOut(s.b)
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

func TestCmdTestSuite(t *testing.T) {
	suite.Run(t, new(CmdTestSuite))
}

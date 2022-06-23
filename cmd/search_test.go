package cmd

import (
	"fmt"
	"github.com/Silthus/go-imap-client/testserver"
	"time"
)

func (s *CmdTestSuite) TestSearchCmd() {
	s.Run("search with invalid credentials returns error", func() {
		s.assertExecuteError("Bad username or password", "search", "--server="+s.testServerAddress, "--username=username", "--password=wrong", "any")
	})

	tests := []struct {
		name           string
		searchTerm     string
		expectedResult string
	}{
		{
			"search with no results has exit code 0",
			"some mail",
			"Found no messages matching the search term: \"some mail\"",
		},
		{
			"search with matches prints message subject",
			"just for you",
			"A little message",
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			result := s.executeSearch(test.searchTerm)
			s.Contains(result, test.expectedResult)
		})
	}
}

func (s *CmdTestSuite) TestSearchWithTls() {
	_, listener := testserver.NewTls(s.T())
	s.Contains(s.executeSearch("just for you", "--tls", "--skip-verify", "--server="+listener.Addr().String()), "A little message")
}

func (s *CmdTestSuite) TestOptionalFlags() {
	tests := []struct {
		flag     string
		value    interface{}
		expected interface{}
	}{
		{"mailbox", "", "INBOX"},
		{"mailbox", "Other", "Other"},
		{"tls", "", "false"},
		{"skip-verify", "", "false"},
		{"timeout", "", fmt.Sprint(5 * time.Second)},
	}
	for _, test := range tests {
		s.Run(test.flag+":"+fmt.Sprintf("%v", test.value)+"->"+fmt.Sprintf("%v", test.expected), func() {
			s.assertOptionalFlag(test.flag, test.value, test.expected)
		})
	}
}

func (s *CmdTestSuite) TestOptionalBoolFlags() {
	s.Run("--tls", func() {
		s.execute("search", "--tls", "any")
		s.True(s.cmd.Flags().GetBool("tls"))
	})
	s.Run("--skip-verify", func() {
		s.execute("search", "--skip-verify", "any")
		s.True(s.cmd.Flags().GetBool("skip-verify"))
	})
}

func (s *CmdTestSuite) TestRequiredFlags() {
	tests := []struct {
		flag string
	}{
		{"server"},
		{"username"},
		{"password"},
	}
	for _, test := range tests {
		s.Run("flag:"+test.flag, func() {
			s.assertRequiredFlag(test.flag, "search", "any")
		})
	}
}

func (s *CmdTestSuite) executeSearch(searchTerm string, flags ...string) string {
	s.T().Helper()
	args := []string{"search", "--server=" + s.testServerAddress, "--username=username", "--password=password", "--timeout=10ms", searchTerm}
	args = append(args, flags...)
	out, err := s.executeErr(args...)
	s.NoError(err)
	return out
}

package cmd

func (s *CmdTestSuite) TestRootCmd() {
	s.Contains(s.execute(), "go-imap-client")
}

func (s *CmdTestSuite) TestLoadConfig() {
	out, err := s.executeErr("--config=../testdata/test-config.yaml", "--server="+s.testServerAddress, "search", "any")
	s.NoError(err)
	s.Regexp("Using config file: \"(.*)test-config.yaml\"", out)

	username, err := s.cmd.Flags().GetString("username")
	s.NoError(err)
	s.Equal("config-username", username)
}

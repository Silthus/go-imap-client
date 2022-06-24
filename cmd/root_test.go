package cmd

func (s *CmdTestSuite) TestRootCmd() {
	_, err := s.executeErr()
	s.NoError(err)
}

func (s *CmdTestSuite) TestLoadConfig() {
	out, _ := s.executeErr("--config=../testdata/test-config.yaml", "--server="+s.testServerAddress, "search", "any")
	s.Regexp("Using config file: \"(.*)test-config.yaml\"", out)

	username, err := s.cmd.Flags().GetString("username")
	s.NoError(err)
	s.Equal("config-test", username)
}

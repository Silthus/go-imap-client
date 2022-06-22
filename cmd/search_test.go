package cmd

func (s *CmdTestSuite) TestSearchCmd() {
	s.Run("search with invalid credentials returns error", func() {
		s.assertExecuteError("access denied: Bad username or password", "search", "--server="+s.testServerAddress, "--user=username", "--password=wrong", "any")
	})
	s.Run("search with no results has exit code 0", func() {
		out := s.execute("search", "--server="+s.testServerAddress, "--user=username", "--password=password", "some mail")
		s.Contains(out, "Found no messages matching the given search term.")
	})
	s.Run("search with matches prints message subject", func() {
		out := s.execute("search", "--server="+s.testServerAddress, "--user=username", "--password=password", "just for you")
		s.Contains(out, "A little message")
	})
}

func (s *CmdTestSuite) TestRequiredFlags() {
	tests := []struct {
		flag string
	}{
		{"server"},
		{"user"},
		{"password"},
	}
	for _, test := range tests {
		s.Run("flag: "+test.flag, func() {
			s.assertRequiredFlag(test.flag, "search", "any")
		})
	}
}

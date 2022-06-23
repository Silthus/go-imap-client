package cmd

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

func (s *CmdTestSuite) executeSearch(searchTerm string) string {
	s.T().Helper()
	out, err := s.executeErr("search", "--server="+s.testServerAddress, "--username=username", "--password=password", searchTerm)
	s.NoError(err)
	return out
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
		s.Run("flag: "+test.flag, func() {
			s.assertRequiredFlag(test.flag, "search", "any")
		})
	}
}

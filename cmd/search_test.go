package cmd

func (s *CmdTestSuite) TestSearchCmd() {
	s.Run("executes search command", func() {
		s.Contains(s.execute("search"), "search")
	})
}

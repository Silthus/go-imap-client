package cmd

func (s *CmdTestSuite) TestRootCmd() {
	s.Run("executes root command", func() {
		s.Contains(s.execute(), "Cobra")
	})
	s.Run("execute without server fails", func() {

	})
}

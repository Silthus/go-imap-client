package cmd

func (s *CmdTestSuite) TestRootCmd() {
	s.Contains(s.execute(), "Cobra")
}

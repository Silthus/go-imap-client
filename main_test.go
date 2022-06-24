package main_test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestMainExec(t *testing.T) {
	cmd := exec.Command("go-imap-client", "--help")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "go-imap-client")
}

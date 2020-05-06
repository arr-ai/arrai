package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExecCommand(t *testing.T) {
	t.Parallel()

	assert.True(t, isExecCommand("eval", cmds))
	assert.True(t, isExecCommand("e", cmds))
	assert.False(t, isExecCommand("./eval", cmds))
	assert.False(t, isExecCommand("eval.arrai", cmds))
	assert.False(t, isExecCommand("file", cmds))
}

func TestFetchCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		"command",
		fetchCommand([]string{"-flag", "--flag", "command"}),
	)

	assert.Equal(t,
		"command",
		fetchCommand([]string{"-flag", "command", "--flag"}),
	)

	assert.Equal(t,
		"",
		fetchCommand([]string{"-flag", "-flag-again", "--flag"}),
	)
}

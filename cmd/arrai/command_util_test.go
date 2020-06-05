package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
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

func TestInsertRunCommand(t *testing.T) {
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "stuff",
			Aliases: []string{"s", "st"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
		},
	}

	assert.Equal(t,
		[]string{"arrai", "-d", "run", "file.arrai"},
		insertRunCommand(flags, []string{"arrai", "-d", "file.arrai"}),
	)
	assert.Equal(t,
		[]string{"arrai", "--debug", "run", "file.arrai"},
		insertRunCommand(flags, []string{"arrai", "--debug", "file.arrai"}),
	)
	assert.Equal(t,
		[]string{"arrai", "-d", "--stuff", "run", "file.arrai"},
		insertRunCommand(flags, []string{"arrai", "-d", "--stuff", "file.arrai"}),
	)
	// separating global flags and subcommand flags
	assert.Equal(t,
		[]string{"arrai", "--debug", "run", "--subcommand-flag", "file.arrai"},
		insertRunCommand(flags, []string{"arrai", "--debug", "--subcommand-flag", "file.arrai"}),
	)
	assert.Equal(t,
		[]string{"arrai", "run", "file.arrai"},
		insertRunCommand(flags, []string{"arrai", "file.arrai"}),
	)

	assert.Panics(t, func() { insertRunCommand(flags, []string{"arrai"}) })
}

package main

import (
	"github.com/arr-ai/arrai/internal/shell"
	"github.com/arr-ai/arrai/rel"
	"github.com/urfave/cli/v2"
)

var shellCommand = &cli.Command{
	Name:    "shell",
	Aliases: []string{"i"},
	Usage:   "start the arrai interactive shell",
	Action:  iShell,
}

func iShell(_ *cli.Context) error {
	return shell.Shell([]rel.ContextErr{})
}

func createDebuggerShell(err error) error {
	if err != nil {
		if ctxErr, isContextError := err.(rel.ContextErr); isContextError {
			return shell.Shell(ctxErr.GetImportantFrames())
		}
		return err
	}
	return nil
}

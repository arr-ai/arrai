package main

import (
	"context"

	"github.com/arr-ai/arrai/internal/shell"
	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/buildinfo"
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
	// when "ai" is used to open the shell, build info is not set. This resets the build info before it is injected
	// into the context.
	buildinfo.SetBuildInfo(Version, BuildDate, GitFullCommit, GitTags, BuildOS, BuildArch, GoVersion)
	return shell.Shell(arraictx.InitRunCtx(context.Background()), []rel.ContextErr{})
}

func createDebuggerShell(err error) error {
	if err != nil {
		if ctxErr, isContextError := err.(rel.ContextErr); isContextError {
			return shell.Shell(arraictx.InitRunCtx(context.Background()), ctxErr.GetImportantFrames())
		}
		return err
	}
	return nil
}

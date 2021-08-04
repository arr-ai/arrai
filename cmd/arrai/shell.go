package main

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/buildinfo"
	"github.com/arr-ai/arrai/pkg/shell"
	"github.com/arr-ai/arrai/rel"
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

// createDebuggerShell creates an interactive shell to explore the context at which err occurred.
func createDebuggerShell(err rel.ContextErr) error {
	return shell.Shell(arraictx.InitRunCtx(context.Background()), err.GetImportantFrames())
}

package main

import (
	"github.com/arr-ai/arrai/pkg/mod"
	"github.com/arr-ai/arrai/tools"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var modCommand = &cli.Command{
	Name:      "mod",
	Usage:     "<command> <repo>",
	UsageText: "mod your dependencies",
	Action:    modRun,
	SkipFlagParsing: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "cmd",
			Usage:   "get or update",
		},
		&cli.StringFlag{
			Name:    "repo",
			Usage:   "repo to get or update (leave blank to update all repos)",
		},
	},
}

func modRun(ctx *cli.Context) error {
	tools.SetArgs(ctx)
	cmd := ctx.Args().Get(0)
	repo := ctx.Args().Get(1)
	return mod.New(afero.NewOsFs()).Command(cmd, repo)
}

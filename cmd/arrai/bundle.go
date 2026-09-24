package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/bundle"
)

var stripSourceFlag = &cli.BoolFlag{
	Name: "strip-source",
	Usage: "Omit the plain arr.ai source tree from the bundle, keeping only the compiled plan. " +
		"Produces a smaller .arraiz, but it can only be run by arrai v0.340.0 or later.",
}

var bundleCommand = &cli.Command{
	Name:    "bundle",
	Aliases: []string{"b"},
	Usage:   "bundles an arr.ai script and its imports into an .arraiz file that can be `arrai run`.",
	Action:  bundleCmd,
	Flags: []cli.Flag{
		outFlag,
		stripSourceFlag,
	},
}

func bundleCmd(c *cli.Context) error {
	file := c.Args().Get(0)
	return bundle.BundledScriptsTo(
		arraictx.InitCliCtx(context.Background(), c),
		file, os.Stdout, c.Value("out").(string), c.Bool("strip-source"),
	)
}

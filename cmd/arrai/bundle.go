package main

import "github.com/urfave/cli/v2"

var bundleCommand = &cli.Command{
	Name:    "bundle",
	Aliases: []string{"b"},
	Usage:   "bundle arrai script and its dependencies into a runnable file",
	Action:  bundle,
	Flags: []cli.Flag{
		outFlag,
	},
}

func bundle(c *cli.Context) error {
	return nil
}

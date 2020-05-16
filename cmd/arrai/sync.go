package main

// TODO: Allow watching a single file to represent a non-tuple database.

import (
	"github.com/urfave/cli/v2"
)

var syncCommand = &cli.Command{
	Name:    "sync",
	Aliases: []string{"s"},
	Usage:   "sync local files to a server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "template",
			Value: "%v",
			Usage: "Template for command to send (%v denotes the content)",
		},
	},
	Action: sync,
}

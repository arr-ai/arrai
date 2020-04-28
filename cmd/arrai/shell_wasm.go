//+build wasm

package main

import (
	"github.com/urfave/cli/v2"
)

var shellCommand = &cli.Command{
	Name:    "shell",
	Aliases: []string{"i"},
	Usage:   "start the arrai interactive shell",
	Action:  shell,
}

func shell(c *cli.Context) error {
	panic("not implemented")
}

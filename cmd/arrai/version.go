package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var versionCommand = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "evaluate the version information of arrai",
	Action:  version,
}

func version(c *cli.Context) error {
	fmt.Println("version 0.1")
	return nil
}

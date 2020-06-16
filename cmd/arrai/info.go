package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var infoCommand = &cli.Command{
	Name:   "info",
	Usage:  "display arrai release information",
	Action: info,
}

func info(c *cli.Context) error {
	fmt.Printf("Version    : %s\n", Version)
	fmt.Printf("Git commit : %s\n", GitFullCommit)
	fmt.Printf("Date       : %s\n", BuildDate)
	fmt.Printf("OS/arch    : %s\n", BuildOS)
	fmt.Printf("Go version : %s\n", GoVersion)

	return nil
}

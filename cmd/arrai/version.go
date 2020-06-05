package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var versionCommand = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "display arrai version information",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "display full arrai version information",
		}},
	Action: version,
}

func version(c *cli.Context) error {
	if c.Bool("verbose") {
		fmt.Printf("Version    : %s\n", Version)
		fmt.Printf("Git commit : %s\n", GitFullCommit)
		fmt.Printf("Date       : %s\n", BuildDate)
		fmt.Printf("OS/arch    : %s\n", BuildOS)
		fmt.Printf("Go version : %s\n", GoVersion)
	} else {
		fmt.Printf("Version    : %s\n", Version)
		fmt.Printf("OS/arch    : %s\n", BuildOS)
	}

	return nil
}

package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var infoCommand = &cli.Command{
	Name:   "info",
	Usage:  "display arrai release information",
	Action: info,
}

func info(c *cli.Context) error {
	fmt.Printf("Version    : %s\n", strings.TrimSpace(Version))
	fmt.Printf("Git commit : %s\n", strings.TrimSpace(GitFullCommit))
	fmt.Printf("Date       : %s\n", strings.TrimSpace(BuildDate))
	fmt.Printf("OS/arch    : %s\n", strings.TrimSpace(BuildOS))
	fmt.Printf("Go version : %s\n", strings.TrimSpace(GoVersion))

	return nil
}

package main

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli/v2"
)

var versionCommand = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "display arrai version information",
	Action:  version,
}

func version(c *cli.Context) error {
	fmt.Printf("Build:\n")
	fmt.Printf("  Version      : %s\n", Version)
	fmt.Printf("  Git Commit   : %s\n", GitCommit)
	fmt.Printf("  Date         : %s\n", BuildDate)
	fmt.Printf("  Go Version   : %s\n", GoVersion)
	fmt.Printf("  OS           : %s\n", BuildOS)
	fmt.Printf("Runtime:\n")
	fmt.Printf("  GOOS/GOARCH  : %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  Go Version   : %s\n", runtime.Version())

	return nil
}

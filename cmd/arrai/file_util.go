package main

import (
	"strings"

	"github.com/urfave/cli/v2"
)

func isExecCommand(arg string, cmds []*cli.Command) bool {
	for _, cmd := range cmds {
		if arg == cmd.Name {
			return true
		}
		for _, name := range cmd.Aliases {
			if arg == name {
				return true
			}
		}
	}
	return false
}

func fetchCommand(args []string) string {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			return arg
		}
	}
	return ""
}

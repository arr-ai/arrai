package main

import (
	"strings"

	"github.com/urfave/cli/v2"
)

func isExecCommand(arg string, cmds []*cli.Command) bool {
	for _, cmd := range cmds {
		names := append(make([]string, 0, 1+len(cmd.Aliases)), cmd.Name)
		names = append(names, cmd.Aliases...)
		for _, name := range names {
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

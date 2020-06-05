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

func insertRunCommand(globalFlags []cli.Flag, args []string) []string {
	globalFlagsEndIndex := 1
	globalFlagsMap := make(map[string]struct{})

	for _, f := range globalFlags {
		for _, n := range f.Names() {
			globalFlagsMap[n] = struct{}{}
		}
	}

	for i, a := range args[1:] {
		if strings.HasPrefix(a, "-") {
			name := strings.TrimLeftFunc(a, func(r rune) bool { return r == '-' })
			if _, isGlobalFlag := globalFlagsMap[name]; isGlobalFlag {
				globalFlagsEndIndex = i + 2
			} else {
				break
			}
		}
	}

	tmpArgs := append(make([]string, 0, globalFlagsEndIndex), args[:globalFlagsEndIndex]...)
	tmpArgs = append(tmpArgs, "run")
	args = append(tmpArgs, args[globalFlagsEndIndex:]...)
	return args
}

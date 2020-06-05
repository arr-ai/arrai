package main

import (
	"fmt"
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

// insertRunCommand adds a run command on file evaluation commands that does not
// contain the `run` command by adding the `run` command manually.
//
// e.g. `arrai path/to/file.arrai`
func insertRunCommand(globalFlags []cli.Flag, args []string) []string {
	if len(args) < 2 {
		panic(fmt.Errorf("need at least 2 arguments, program and arrai file, received %v", args))
	}

	globalFlagsEndIndex := 1
	globalFlagsMap := make(map[string]bool)

	for _, f := range globalFlags {
		for _, n := range f.Names() {
			globalFlagsMap[n] = true
		}
	}

	for i, a := range args[1:] {
		if strings.HasPrefix(a, "-") {
			name := strings.TrimLeft(a, "-")
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

//+build !wasm

package main

import (
	"regexp"
	"strings"

	"github.com/go-errors/errors"
)

var cmdRe = regexp.MustCompile(`^/[a-zA-Z_]\w*[ \t]*`)

func isCommand(line string) bool {
	return cmdRe.MatchString(line)
}

func tryRunCommand(line string, shellData *shellInstance) error {
	var name string
	if name = cmdRe.FindString(line); name == "" {
		return errors.Errorf("%s is not a command", line)
	}
	name = strings.TrimSpace(name[1:])
	line = strings.TrimSpace(cmdRe.ReplaceAllString(line, ""))
	cmd, isCmd := shellData.cmds[name]
	if !isCmd {
		return errors.Errorf("command %s not found", name)
	}
	return cmd.process(line, shellData)
}

type shellCmd interface {
	name() string
	process(line string, shellData *shellInstance) error
}

type setCmd struct{}

func (sc *setCmd) name() string {
	return "set"
}

func (sc *setCmd) process(line string, shellData *shellInstance) error {
	identRe := regexp.MustCompile(`^(?P<ident>[a-zA-Z_]\w*)[ \t]+=`)
	identMatches := identRe.FindStringSubmatch(line)
	if len(identMatches) != 2 {
		return errors.Errorf(`/set command error, usage: /set <name> = <expr>`)
	}
	name := identMatches[1]
	expr := identRe.ReplaceAllString(line, "")
	val, err := shellEval(expr, shellData.scope)
	if err != nil {
		return err
	}
	shellData.scope = shellData.scope.With(name, val)
	return nil
}

type unsetCmd struct{}

func (uc *unsetCmd) name() string {
	return "unset"
}

func (uc *unsetCmd) process(line string, shellData *shellInstance) error {
	ident := regexp.MustCompile(`^[a-zA-Z_]\w*`).FindString(line)
	if ident == "" {
		return errors.Errorf(`/unset command error, usage: /unset <name>`)
	}
	shellData.scope = shellData.scope.Without(ident)
	return nil
}

type exitError struct{}

func (exitError) Error() string {
	return "exiting interactive shell"
}

type exitCommand struct{}

func (*exitCommand) name() string {
	return "exit"
}

func (ec *exitCommand) process(_ string, _ *shellInstance) error {
	return exitError{}
}

func initCommands() map[string]shellCmd {
	cmds := []shellCmd{&setCmd{}, &unsetCmd{}, &exitCommand{}}
	cmdMap := make(map[string]shellCmd)
	for _, cmd := range cmds {
		cmdMap[cmd.name()] = cmd
	}
	return cmdMap
}

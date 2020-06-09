//+build !wasm

package shell

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/anz-bank/pkg/log"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/wbnf/parser"
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
	names() []string
	process(line string, shellData *shellInstance) error
}

type setCmd struct{}

func (sc *setCmd) names() []string {
	return []string{"set"}
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

func (uc *unsetCmd) names() []string {
	return []string{"unset"}
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

func (*exitCommand) names() []string {
	return []string{"exit"}
}

func (ec *exitCommand) process(_ string, _ *shellInstance) error {
	return exitError{}
}

type upFrameCmd struct{}

func (*upFrameCmd) names() []string {
	return []string{"up", "u"}
}

func (*upFrameCmd) process(_ string, sh *shellInstance) error {
	return changeFrame(sh.currentFrameIndex-1, sh)
}

type downFrameCmd struct{}

func (d *downFrameCmd) names() []string {
	return []string{"down", "d"}
}

func (*downFrameCmd) process(_ string, sh *shellInstance) error {
	return changeFrame(sh.currentFrameIndex+1, sh)
}

func changeFrame(i int, sh *shellInstance) error {
	if i < 0 || i >= len(sh.frames) {
		return fmt.Errorf("frame index out of range, frame length: %d", len(sh.frames))
	}
	sh.currentFrameIndex = i
	log.Infof(context.Background(), "Stack: %d\n%s\n", i, sh.frames[i].GetSource().Context(parser.DefaultLimit))
	sh.scope = syntax.StdScope().Update(sh.frames[i].GetScope())
	return nil
}

func initCommands() map[string]shellCmd {
	cmds := []shellCmd{&setCmd{}, &unsetCmd{}, &exitCommand{}, &upFrameCmd{}, &downFrameCmd{}}
	cmdMap := make(map[string]shellCmd)
	for _, cmd := range cmds {
		for _, n := range cmd.names() {
			cmdMap[n] = cmd
		}
	}
	return cmdMap
}

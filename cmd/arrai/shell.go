//+build !wasm

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/anz-bank/pkg/log"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/chzyer/readline"
	"github.com/urfave/cli/v2"
)

var shellCommand = &cli.Command{
	Name:    "shell",
	Aliases: []string{"i"},
	Usage:   "start the arrai interactive shell",
	Action:  shell,
}

func tryEval(line string, scope rel.Scope) (_ rel.Value, err error) {
	defer func() {
		if i := recover(); i != nil {
			if i, is := i.(error); is {
				err = i
			} else {
				err = fmt.Errorf("")
			}
		}
	}()
	return syntax.EvalWithScope("", line, scope)
}

func shell(c *cli.Context) error {
	ctx := log.WithConfigs(log.SetVerboseMode(true)).Onto(context.Background())
	sh := newShellInstance(newLineCollector(), syntax.StdScope())
	l, err := readline.NewEx(&readline.Config{
		Prompt:       "@> ",
		HistoryFile:  os.ExpandEnv("${HOME}/.arrai_history"),
		AutoComplete: sh,
		EOFPrompt:    "exit",
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		line, err := l.Readline()
		if err != nil {
			switch err {
			case io.EOF:
				return nil
			case readline.ErrInterrupt:
				continue
			}
			panic(err)
		}
		if err = sh.parseCmd(line, l); err != nil {
			log.Error(ctx, err)
		}
	}
}

type shellInstance struct {
	collector *lineCollector
	scope     rel.Scope
	cmds      map[string]shellCmd
}

func newShellInstance(c *lineCollector, initialScope rel.Scope) *shellInstance {
	return &shellInstance{c, initialScope, initCommands()}
}

func (s *shellInstance) parseCmd(line string, l *readline.Instance) error {
	if line = strings.TrimSpace(line); line != "" {
		s.collector.appendLine(line)
	}
	if len(s.collector.lines) != 0 && s.collector.isBalanced() {
		l.SetPrompt("@> ")
		lines := strings.Join(s.collector.lines, "\n")
		s.collector.reset()
		if isCommand(lines) {
			return tryRunCommand(lines, s)
		} else if _, err := shellEval(lines, s.scope); err != nil {
			return err
		}
	}
	if len(s.collector.lines) != 0 {
		l.SetPrompt(" > ")
	}
	return nil
}

func (s *shellInstance) Do(line []rune, pos int) (newLine [][]rune, length int) {
	l := string(line[:pos])
	switch {
	case strings.HasSuffix(l, "///"):
		return [][]rune{line}, len(line)
	case strings.HasPrefix(l, "//"):
		var names []string
		var lastName string
		if l == "//" {
			newLine = append(newLine, []rune("{"))
			lastName, names = "", []string{}
		} else {
			names = strings.Split(l[2:], ".")
			lastName, names = names[len(names)-1], names[:len(names)-1]
		}

		t := s.scope.MustGet(".").(rel.Tuple)
		for _, name := range names {
			if value, has := t.Get(name); has {
				if u, is := value.(rel.Tuple); is {
					t = u
					continue
				}
			}
			return
		}

		length = len(lastName)
		for _, name := range t.Names().OrderedNames() {
			if strings.HasPrefix(name, lastName) {
				newLine = append(newLine, []rune(name[length:]))
			}
		}
	}
	return newLine, length
}

func shellEval(lines string, scope rel.Scope) (rel.Value, error) {
	value, err := tryEval(lines, scope)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stdout, "%s\n", rel.Repr(value))
	return value, nil
}

type lineCollector struct {
	lines        []string
	stack        []*closer
	opener       map[string]*closer
	maxOpenerLen int
}

type closer struct {
	char               string
	recursive          bool
	contextBasedOpener map[string]*closer
}

func (l *lineCollector) pop() {
	if len(l.stack) == 0 {
		return
	}
	l.stack = l.stack[:len(l.stack)-1]
}

func (l *lineCollector) push(close *closer) {
	l.stack = append(l.stack, close)
}

func (l *lineCollector) peek() *closer {
	if len(l.stack) == 0 {
		return nil
	}
	return l.stack[len(l.stack)-1]
}

func newLineCollector() *lineCollector {
	templateContext := map[string]*closer{"${": {"}", true, nil}}
	return &lineCollector{
		[]string{},
		[]*closer{},
		map[string]*closer{
			"{":   {"}", true, nil},
			"(":   {")", true, nil},
			"[":   {"]", true, nil},
			"$\"": {"\"", true, templateContext},
			"$'":  {"'", true, templateContext},
			"$`":  {"`", true, templateContext},
			"\"":  {"\"", false, nil},
			"'":   {"'", false, nil},
			"`":   {"`", false, nil},
		},
		2, //TODO: automate finding max length
	}
}

func (l *lineCollector) appendLine(line string) {
	increment := 1
	for i := 0; i < len(line); i += increment {
		if line[i] == '\\' {
			increment = 2
		} else if nextCloser := l.peek(); nextCloser != nil && strings.HasPrefix(line[i:], nextCloser.char) {
			l.pop()
			increment = len(nextCloser.char)
		} else {
			if nextCloser != nil && !nextCloser.recursive {
				continue
			}
			openers := l.opener
			if nextCloser != nil && nextCloser.contextBasedOpener != nil {
				openers = nextCloser.contextBasedOpener
			}
			possibleOpener := make([]string, 0, l.maxOpenerLen)
			for j := 0; j < l.maxOpenerLen && i+j < len(line); j++ {
				possibleOpener = append(possibleOpener, line[i:i+j+1])
			}
			for _, p := range possibleOpener {
				if close, isOpener := openers[p]; isOpener {
					l.push(close)
					increment = len(p)
					break
				} else {
					increment = 1
				}
			}
		}
	}
	l.lines = append(l.lines, line)
}

func (l *lineCollector) isBalanced() bool {
	if nextClosure := l.peek(); nextClosure != nil && !nextClosure.recursive {
		return false
	}

	if len(l.lines) == 0 {
		return true
	}

	lastLine := l.lines[len(l.lines)-1]
	if strings.HasSuffix(lastLine, ";") || strings.HasSuffix(lastLine, ":") {
		return false
	}

	// check for function argument
	if regexp.MustCompile(`\\[^ \t\n]+$`).Match([]byte(lastLine)) {
		return false
	}

	return len(l.stack) == 0
}

func (l *lineCollector) reset() {
	l.lines = []string{}
	l.stack = []*closer{}
}

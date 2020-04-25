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

func tryEval(line string) (_ rel.Value, err error) {
	defer func() {
		if i := recover(); i != nil {
			if i, is := i.(error); is {
				err = i
			} else {
				err = fmt.Errorf("")
			}
		}
	}()
	return syntax.EvaluateExpr("", line)
}

type autoCompleter struct {
	std rel.Tuple
}

var matchStdPrefixRE = regexp.MustCompile(`//((?:\.\w+)*)(\.?)$`)

func newAutoCompleter() *autoCompleter {
	return &autoCompleter{
		std: syntax.StdScope().MustGet(".").(rel.Tuple),
	}
}

func (c *autoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	s := string(line[:pos])
	switch {
	case strings.HasSuffix(s, "///"):
		return [][]rune{line}, len(line)
	case strings.HasSuffix(s, "//"):
		return [][]rune{[]rune("."), []rune("{")}, 0
	}
	if m := matchStdPrefixRE.FindStringSubmatch(s); m != nil {
		t := c.std
		lastName := ""
		if m[1] != "" {
			names := strings.Split(m[1][1:], ".")
			allNamesButLast := names
			if m[2] != "." {
				allNamesButLast = names[:len(names)-1]
				lastName = names[len(names)-1]
			}
			for _, name := range allNamesButLast {
				if value, has := t.Get(name); has {
					if u, is := value.(rel.Tuple); is {
						t = u
					} else {
						return
					}
				} else {
					return
				}
			}
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

func shell(c *cli.Context) error {
	ctx := log.WithConfigs(log.SetVerboseMode(true)).Onto(context.Background())
	l, err := readline.NewEx(&readline.Config{
		Prompt:       "@> ",
		HistoryFile:  os.ExpandEnv("${HOME}/.arrai_history"),
		AutoComplete: newAutoCompleter(),
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
		if line != "" {
			value, err := tryEval(line)
			if err != nil {
				log.Error(ctx, err)
			} else {
				fmt.Fprintf(os.Stdout, "%s\n", rel.Repr(value))
			}
		}
	}
}

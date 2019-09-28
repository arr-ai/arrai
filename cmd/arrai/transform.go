package main

import (
	"os"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/urfave/cli"
)

var transformCommand = cli.Command{
	Name:    "transform",
	Aliases: []string{"x"},
	Usage:   "transform a stream of input data with an expression",
	Action:  transform,
}

func transform(c *cli.Context) error {
	source := c.Args().Get(0)

	expr, err := syntax.Parse(syntax.NewStringLexer(source))
	if err != nil {
		return err
	}

	stream := syntax.NewLexer(os.Stdin)
	for {
		if stream.Peek() == syntax.EOF {
			break
		}
		value, err := syntax.ParseAtom(stream)
		if err != nil {
			return err
		}

		global := rel.EmptyScope.With(".", value)
		xvalue, err := expr.Eval(global, global)
		if err != nil {
			return err
		}

		s := xvalue.String()
		os.Stdout.WriteString(s)
		if s[len(s)-1] != '\n' {
			os.Stdout.Write([]byte{'\n'})
		}
	}

	return nil
}

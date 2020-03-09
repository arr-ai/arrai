package main

import (
	"fmt"
	"io"
	"os"

	"github.com/arr-ai/arrai/syntax"
	"github.com/urfave/cli/v2"
)

var evalCommand = &cli.Command{
	Name:    "eval",
	Aliases: []string{"e"},
	Usage:   "evaluate an arrai expression",
	Action:  eval,
}

func evalImpl(source string, w io.Writer) error {
	return evalExpr(".", source, w)
}

func evalExpr(path, source string, w io.Writer) error {
	value, err := syntax.EvaluateExpr(path, source)
	if err != nil {
		return err
	}

	s := value.String()
	fmt.Fprintf(w, "%s", s)
	if s[len(s)-1] != '\n' {
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	return nil
}

func eval(c *cli.Context) error {
	source := c.Args().Get(0)
	return evalImpl(source, os.Stdout)
}

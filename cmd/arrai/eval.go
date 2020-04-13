package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arr-ai/arrai/rel"
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
	return evalExpr(syntax.NoPath, source, w)
}

func evalExpr(path, source string, w io.Writer) error {
	value, err := syntax.EvaluateExpr(path, source)
	if err != nil {
		return err
	}

	var s string
	switch v := value.(type) {
	case rel.String:
		s = v.String()
	case rel.Set:
		if !v.IsTrue() {
			s = ""
		} else {
			s = rel.Repr(v)
		}
	default:
		s = rel.Repr(v)
	}
	fmt.Fprintf(w, "%s", s)
	if s != "" && !strings.HasSuffix(s, "\n") {
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

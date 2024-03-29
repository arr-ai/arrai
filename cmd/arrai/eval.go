package main

import (
	"context"
	"io"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/arr-ai/arrai/pkg/arrai"
	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/syntax"
)

var outFlag = &cli.StringFlag{
	Name:    "out",
	Aliases: []string{"o"},
	Usage:   "Control output behaviour",
}

var evalCommand = &cli.Command{
	Name:    "eval",
	Aliases: []string{"e"},
	Usage:   "evaluate an arrai expression",
	Action:  eval,
	Flags: []cli.Flag{
		outFlag,
	},
}

func eval(c *cli.Context) error {
	source := c.Args().Get(0)
	ctx := arraictx.InitCliCtx(context.Background(), c)

	return evalImpl(ctx, source, os.Stdout, c.Value("out").(string))
}

func evalImpl(ctx context.Context, source string, w io.Writer, out string) error {
	return evalExpr(ctx, ".", source, w, out)
}

func evalExpr(ctx context.Context, path, source string, w io.Writer, out string) error {
	value, err := syntax.EvaluateExpr(ctx, path, source)
	if err != nil {
		return err
	}

	return arrai.OutputValue(ctx, value, w, out)
}

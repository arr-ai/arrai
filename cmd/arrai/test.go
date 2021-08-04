package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/test"
)

var testCommand = &cli.Command{
	Name:    "test",
	Aliases: []string{"t"},
	Usage:   "run arrai tests",
	Action:  doTest,
}

func doTest(c *cli.Context) error {
	path := c.Args().Get(0)
	ctx := arraictx.InitCliCtx(context.Background(), c)

	return test.RunTests(ctx, os.Stdout, path)
}

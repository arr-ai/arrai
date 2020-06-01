package main

import (
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

var runCommand = &cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "evaluate an arrai file",
	Action:  run,
}

func run(c *cli.Context) error {
	file := c.Args().Get(0)
	return evalFile(file, os.Stdout)
}

func evalFile(path string, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	buf := make([]byte, fi.Size())
	if _, err := f.Read(buf); err != nil {
		return err
	}
	return createDebuggerShell(evalExpr(path, string(buf), w))
}

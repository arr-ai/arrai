package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arr-ai/arrai/tools"

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
	if !tools.FileExists(path) {
		if !strings.Contains(path, string([]rune{os.PathSeparator})) {
			return fmt.Errorf(`"%s": not a command and not found as a file in the current directory`, path)
		}
		return fmt.Errorf(`"%s": file not found`, path)
	}

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
	return evalExpr(path, string(buf), w)
}

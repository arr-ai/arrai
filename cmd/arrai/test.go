package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arr-ai/arrai/internal/test"
	"github.com/arr-ai/arrai/tools"

	"github.com/urfave/cli/v2"
)

var testCommand = &cli.Command{
	Name:    "test",
	Aliases: []string{"t"},
	Usage:   "run arrai tests",
	Action:  doTest,
}

func doTest(c *cli.Context) error {
	tools.SetArgs(c)
	file := c.Args().Get(0)
	return testPath(file, os.Stdout)
}

// testPath runs and reports on all tests in the subtree of the given path.
func testPath(path string, w io.Writer) error {
	if path == "" {
		path = "."
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if !strings.Contains(path, string([]rune{os.PathSeparator})) {
			return fmt.Errorf(`"%s": not a command and not found as a file in the current directory`, path)
		}
		return fmt.Errorf(`"%s": file not found`, path)
	}

	results, err := test.Test(path)
	if err != nil {
		return err
	}

	test.Report(results)

	return nil
}

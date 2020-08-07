package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/tools"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var runCommand = &cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "evaluate an arrai file",
	Action:  run,
	Flags: []cli.Flag{
		outFlag,
	},
}

func run(c *cli.Context) error {
	tools.SetArgs(c)
	file := c.Args().Get(0)

	ctx := context.Background()
	ctx = ctxfs.SourceFsOnto(ctx, afero.NewOsFs())
	ctx = ctxfs.RuntimeFsOnto(ctx, afero.NewOsFs())

	return evalFile(ctx, file, os.Stdout, c.Value("out").(string))
}

func evalFile(ctx context.Context, path string, w io.Writer, out string) error {
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
	return evalExpr(ctx, path, string(buf), w, out)
}

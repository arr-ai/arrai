package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/arrai/tools"
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

	return evalFile(arraictx.InitRunCtx(context.Background()), file, os.Stdout, c.Value("out").(string))
}

func evalFile(ctx context.Context, path string, w io.Writer, out string) error {
	if err := runFileExists(ctx, path); err != nil {
		return err
	}
	buf, err := ctxfs.ReadFile(ctxfs.SourceFsFrom(ctx), path)
	if err != nil {
		return err
	}

	switch filepath.Ext(path) {
	case ".arraiz", ".zip":
		return runBundled(ctx, buf, w, out)
	}

	return evalExpr(ctx, path, string(buf), w, out)
}

func runFileExists(ctx context.Context, path string) error {
	if exists, err := tools.FileExists(ctx, path); err != nil {
		return err
	} else if !exists {
		if !strings.Contains(path, string([]rune{os.PathSeparator})) {
			return fmt.Errorf(`"%s": not a command and not found as a file in the current directory`, path)
		}
		return fmt.Errorf(`"%s": file not found`, path)
	}
	return nil
}

func runBundled(ctx context.Context, buf []byte, w io.Writer, out string) error {
	ctx, err := syntax.WithBundleRun(ctx, buf)
	if err != nil {
		return err
	}

	ctx, mainFile, path := syntax.GetMainBundleSource(ctx)
	return evalExpr(ctx, path, string(mainFile), w, out)
}

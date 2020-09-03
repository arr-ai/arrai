package main

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/arrai/tools"
	"github.com/urfave/cli/v2"
)

const bundledType = ".arraiz"

var bundleCommand = &cli.Command{
	Name:    "bundle",
	Aliases: []string{"b"},
	Usage:   "bundles an arr.ai script and its imports into an .arraiz file that can be `arrai run`.",
	Action:  bundle,
	Flags: []cli.Flag{
		outFlag,
	},
}

func bundle(c *cli.Context) error {
	tools.SetArgs(c)
	file := c.Args().Get(0)
	return bundleFiles(
		arraictx.InitRunCtx(context.Background()),
		file, os.Stdout, c.Value("out").(string),
	)
}

func bundleFiles(ctx context.Context, path string, w io.Writer, out string) (err error) {
	if err := runFileExists(ctx, path); err != nil {
		return err
	}

	if out != "" {
		if ext := filepath.Ext(out); ext != bundledType {
			out += bundledType
		}

		f, err := ctxfs.SourceFsFrom(ctx).Create(out)
		if err != nil {
			return err
		}
		w = f
	}

	buf, err := ctxfs.ReadFile(ctxfs.SourceFsFrom(ctx), path)
	if err != nil {
		return err
	}

	if ctx, err = syntax.SetupBundle(ctx, path, buf); err != nil {
		return err
	}

	if _, err = syntax.Compile(ctx, path, string(buf)); err != nil {
		return err
	}

	return syntax.OutputArraiz(ctx, w)
}

package main

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/urfave/cli/v2"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/syntax"
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
	file := c.Args().Get(0)
	return bundleFilesTo(
		arraictx.InitCliCtx(context.Background(), c),
		file, os.Stdout, c.Value("out").(string),
	)
}

func bundleFiles(ctx context.Context, path string, w io.Writer) error {
	return bundleFilesTo(ctx, path, w, "")
}

func bundleFilesTo(ctx context.Context, path string, w io.Writer, out string) (err error) {
	if err := fileExists(ctx, path); err != nil {
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

	buf, err := afero.ReadFile(ctxfs.SourceFsFrom(ctx), path)
	if err != nil {
		return err
	}

	if ctx, err = syntax.SetupBundle(ctx, path, buf); err != nil {
		return err
	}

	if _, err = syntax.Compile(importcache.WithNewImportCache(ctx), path, string(buf)); err != nil {
		return err
	}

	return syntax.OutputArraiz(ctx, w)
}

package bundle

import (
	"context"
	"io"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/arr-ai/arrai/pkg/cliutil"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/syntax"
)

const bundledType = ".arraiz"

// BundledScripts bundle scripts and writes the byte output to the provided writer.
func BundledScripts(ctx context.Context, path string, w io.Writer) error {
	return BundledScriptsTo(ctx, path, w, "")
}

// BundledScriptsTo bundle scripts and outputs it to a file.
func BundledScriptsTo(ctx context.Context, path string, w io.Writer, out string) (err error) {
	if err := cliutil.FileExists(ctx, path); err != nil {
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

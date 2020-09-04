package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/arrai/tools"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
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
	tools.SetArgs(c)
	source := c.Args().Get(0)

	ctx := arraictx.InitRunCtx(context.Background())

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

	return handleEvalOut(ctx, value, w, out)
}

func handleEvalOut(ctx context.Context, value rel.Value, w io.Writer, out string) error {
	if out != "" {
		return outputValue(ctx, value, out)
	}

	var s string
	switch v := value.(type) {
	case rel.String:
		s = v.String()
	case rel.Bytes:
		s = v.String()
	case rel.Set:
		if !v.IsTrue() {
			s = ""
		} else {
			s = rel.Repr(v)
		}
	default:
		s = rel.Repr(v)
	}
	fmt.Fprintf(w, "%s", s)
	if s != "" && !strings.HasSuffix(s, "\n") {
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	return nil
}

func outputValue(ctx context.Context, value rel.Value, out string) error {
	parts := strings.SplitN(out, ":", 2)
	if len(parts) == 1 {
		parts = []string{"", parts[0]}
	}
	mode := parts[0]
	arg := parts[1]

	fs := ctxfs.RuntimeFsFrom(ctx)
	switch mode {
	case "file", "f", "":
		return outputFile(value, arg, fs, false)
	case "dir", "d":
		if t, is := value.(rel.Dict); is {
			if err := outputTupleDir(t, arg, fs, true); err != nil {
				return err
			}
			return outputTupleDir(t, arg, fs, false)
		}
		return fmt.Errorf("result not a dict: %v", value)
	}
	return fmt.Errorf("invalid --out flag: %s", out)
}

func outputTupleDir(t rel.Dict, dir string, fs afero.Fs, dryRun bool) error {
	if _, err := fs.Stat(dir); os.IsNotExist(err) {
		if err := fs.Mkdir(dir, 0755); err != nil {
			return err
		}
	}
	for e := t.DictEnumerator(); e.MoveNext(); {
		k, v := e.Current()
		name, is := k.(rel.String)
		if !is {
			return fmt.Errorf("dir output dict key must be a non-empty string")
		}
		subpath := path.Join(dir, name.String())
		switch content := v.(type) {
		case rel.Dict:
			if err := outputTupleDir(content, subpath, fs, dryRun); err != nil {
				return err
			}
		case rel.Bytes, rel.String:
			if err := outputFile(content, subpath, fs, dryRun); err != nil {
				return err
			}
		case rel.Set:
			if content.IsTrue() {
				return fmt.Errorf("dir output entry must be dict, string or byte array")
			}
			if err := outputFile(content, subpath, fs, dryRun); err != nil {
				return err
			}
		}
	}
	return nil
}

func outputFile(content rel.Value, path string, fs afero.Fs, dryRun bool) error {
	var bytes []byte
	switch content := content.(type) {
	case rel.Bytes:
		bytes = content.Bytes()
	case rel.String:
		bytes = []byte(content.String())
	default:
		if _, is := content.(rel.Set); !(is && !content.IsTrue()) {
			return fmt.Errorf("file output not string or byte array: %v", content)
		}
		bytes = []byte{}
	}

	if dryRun {
		return nil
	}

	f, err := fs.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}

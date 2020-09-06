package main

import (
	"bytes"
	"context"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/arr-ai/arrai/pkg/arrai"
	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/arrai/tools"
	"github.com/arr-ai/wbnf/parser"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var compileCommand = &cli.Command{
	Name:    "compile",
	Aliases: []string{"compile"},
	Usage:   "compile arrai scripts into a runnable binary",
	Action:  compile,
	Flags: []cli.Flag{
		outFlag,
	},
}

func compile(c *cli.Context) error {
	tools.SetArgs(c)
	file := c.Args().Get(0)

	return compileFile(arraictx.InitRunCtx(context.Background()), file, c.Value("out").(string))
}

func compileFile(ctx context.Context, path, out string) error {
	if err := runFileExists(ctx, path); err != nil {
		return err
	}

	bundledScripts := bytes.Buffer{}
	if err := bundleFiles(ctx, path, &bundledScripts, ""); err != nil {
		return err
	}

	if out == "" {
		out = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	f, err := ctxfs.SourceFsFrom(ctx).Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := buildBinary(ctx, bundledScripts.Bytes(), f); err != nil {
		return err
	}
	return nil
}

func createGoFile(ctx context.Context, bundledScripts []byte) ([]byte, error) {
	template, err := Asset("internal/build/main.arrai")
	if err != nil {
		return nil, err
	}

	templateFn, err := syntax.EvaluateExpr(ctx, "", string(template))
	if err != nil {
		return nil, err
	}

	bundledBytes := make([]rel.Value, 0, len(bundledScripts))
	for _, b := range bundledScripts {
		bundledBytes = append(bundledBytes, rel.NewNumber(float64(b)))
	}

	res, err := rel.NewCallExpr(parser.Scanner{}, templateFn, rel.NewArray(bundledBytes...)).Eval(ctx, rel.EmptyScope)
	if err != nil {
		return nil, err
	}

	mainGo := bytes.Buffer{}
	if err = arrai.OutputValue(ctx, res, &mainGo, ""); err != nil {
		return nil, err
	}
	return mainGo.Bytes(), nil
}

func buildBinary(ctx context.Context, bundledScripts []byte, out afero.File) error {
	goFile, err := createGoFile(ctx, bundledScripts)
	if err != nil {
		return err
	}

	_, module, err := syntax.GetModuleFromBundle(ctx, bundledScripts)
	if err != nil {
		return nil
	}

	fs := ctxfs.SourceFsFrom(ctx)

	buildDir, err := afero.TempDir(fs, "", path.Base(module)+"*")
	if err != nil {
		return err
	}
	defer fs.RemoveAll(buildDir)

	mainFilePath := filepath.Join(buildDir, "main.go")
	f, err := fs.Create(mainFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(goFile); err != nil {
		return err
	}

	cmds := []*exec.Cmd{
		exec.CommandContext(ctx, "go", "mod", "init", module),
		exec.CommandContext(ctx, "go", "mod", "tidy"),
		exec.CommandContext(ctx, "go", "build", "-o", "main", mainFilePath),
	}

	for _, c := range cmds {
		c.Dir = buildDir
		if err = c.Run(); err != nil {
			return err
		}
	}

	file, err := ctxfs.ReadFile(fs, filepath.Join(buildDir, "main"))
	if err != nil {
		return err
	}

	if _, err = out.Write(file); err != nil {
		return err
	}

	// 0751 for rwxr-x--x the same as golang binary
	return fs.Chmod(out.Name(), 0751)
}

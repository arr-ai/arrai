package main

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func assertEvalOutputs(t *testing.T, expected, source string) bool { //nolint:unparam
	var sb strings.Builder
	return assert.NoError(t, evalImpl(arraictx.InitRunCtx(context.Background()), source, &sb, "")) &&
		assert.Equal(t, expected, strings.TrimRight(sb.String(), "\n"))
}

func assertEvalCreates(t *testing.T, expected map[string]string, source, out string) bool { //nolint:unparam
	var sb strings.Builder

	memFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")
	ctx := ctxfs.RuntimeFsOnto(arraictx.InitRunCtx(context.Background()), memFs)

	stdoutOK := assert.NoError(t, evalImpl(ctx, source, &sb, out)) &&
		assert.Equal(t, "", strings.TrimRight(sb.String(), "\n"))
	outOK := true
	for path, content := range expected {
		fs, err := memFs.Open(path)
		if assert.NoError(t, err) {
			data, err := ioutil.ReadAll(fs)
			if assert.NoError(t, err) && assert.Equal(t, content, string(data), "path=%s", path) {
				continue
			}
		}
		outOK = false
	}
	return stdoutOK && outOK
}

func TestEvalNumberULP(t *testing.T) {
	t.Parallel()
	assertEvalOutputs(t, `0.3`, `0.1 + 0.1 + 0.1`)
}

func TestEvalString(t *testing.T) {
	t.Parallel()
	assertEvalOutputs(t, ``, `""`)
	assertEvalOutputs(t, ``, `{}`)
	assertEvalOutputs(t, `abc`, `"abc"`)
}

func TestEvalComplex(t *testing.T) {
	t.Parallel()
	assertEvalOutputs(t, `[42, 'abc']`, `[42, "abc"]`)
	assertEvalOutputs(t, `{42, 'abc'}`, `{"abc", 42}`)
}

func TestEvalOutFile(t *testing.T) {
	t.Parallel()
	assertEvalCreates(t, map[string]string{
		"/output": "hello",
	}, `"hello"`, "file:/output")
	assertEvalCreates(t, map[string]string{
		"/output": "hello again",
	}, `"hello again"`, "f:/output")
	assertEvalCreates(t, map[string]string{
		"/output": "hello yet again",
	}, `"hello yet again"`, ":/output")
	assertEvalCreates(t, map[string]string{
		"/output": "we must stop meeting like this",
	}, `"we must stop meeting like this"`, "/output")
}

func TestEvalOutDir(t *testing.T) {
	t.Parallel()
	assertEvalCreates(t, map[string]string{
		"/files/foo": "hello\n",
		"/files/bar": "goodbye\n",
	}, `{"foo": "hello\n", "bar": "goodbye\n"}`, "dir:/files")
	assertEvalCreates(t, map[string]string{
		"/files/foo":     "hello\n",
		"/files/bar/baz": "goodbye\n",
	}, `{"foo": "hello\n", "bar/baz": "goodbye\n"}`, "d:/files")
}

package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileFile(t *testing.T) {
	t.Parallel()

	testCompileScript(t, "1 + 2", "github.com/example/named", "3")
	testCompileScript(t, "2 + 3", "", "5")
	testCompileScript(t,
		"let [_, ...tail] = //os.args; tail",
		"", "['test', '123', 'abc']", "test", "123", "abc",
	)
}

func testCompileScript(t *testing.T, script, mod, expected string, args ...string) {
	ctx := arraictx.InitRunCtx(context.Background())
	fs := ctxfs.SourceFsFrom(ctx)
	temp, err := afero.TempDir(fs, "", "*")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, fs.RemoveAll(temp))
	}()
	filename := filepath.Join(temp, "test.arrai")
	f, err := fs.Create(filename)
	require.NoError(t, err)
	_, err = f.Write([]byte(script))
	require.NoError(t, err)
	f.Close()

	if mod != "" {
		goMod, err := fs.Create(filepath.Join(temp, "go.mod"))
		require.NoError(t, err)
		_, err = goMod.Write([]byte(fmt.Sprintf("module %s\n", mod)))
		require.NoError(t, err)
		goMod.Close()
	}

	outFile := filepath.Join(temp, "test")
	assert.NoError(t, compileFile(ctx, filename, outFile))
	testExec(t, outFile, expected+"\n", args...)
}

func testExec(t *testing.T, path, expected string, args ...string) {
	path, err := filepath.Abs(path)
	assert.NoError(t, err)
	c := exec.Command(path, args...)
	actual := bytes.Buffer{}
	c.Stdout = &actual

	assert.NoError(t, c.Run())

	assert.Equal(t, expected, actual.String())
}

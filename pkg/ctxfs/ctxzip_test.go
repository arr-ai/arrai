package ctxfs

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testKey int

const key testKey = iota

func TestOutputZip(t *testing.T) {
	t.Parallel()

	files := map[string]string{
		"/test/test.txt":      "123",
		"/test.txt":           "321",
		"/more/more/test.txt": "test",
	}

	ctx := WithZipFs(context.Background(), key)

	for name, content := range files {
		require.NoError(t, ZipCreate(ctx, key, name, []byte(content)))
	}

	actual := &bytes.Buffer{}
	assert.NoError(t, OutputZip(ctx, key, actual))

	ZipEqualToFiles(t, actual.Bytes(), files)
}

func TestZipCreate(t *testing.T) {
	t.Parallel()

	ctx := WithZipFs(context.Background(), key)

	filePath := "test/test/test"
	err := ZipCreate(ctx, key, filePath, []byte{})
	assert.EqualError(t, err, "path has to be absolute in UNIX path format: test/test/test")

	filePath = "/test/test.txt"
	content := []byte("123")

	assert.NoError(t, ZipCreate(ctx, key, filePath, content))
	assertFileEqual(t, ctx.Value(key).(*zipFs).fs, filePath, content)

	newContent := []byte("321")
	assert.NoError(t, ZipCreate(ctx, key, filePath, newContent))
	assertFileEqual(t, ctx.Value(key).(*zipFs).fs, filePath, content)
}

func assertFileEqual(t *testing.T, fs afero.Fs, filePath string, content []byte) {
	f, err := fs.Open(filePath)
	assert.NoError(t, err)

	_, err = f.Stat()
	assert.NoError(t, err)

	buf, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, content, buf)
}

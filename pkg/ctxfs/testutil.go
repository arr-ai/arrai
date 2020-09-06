package ctxfs

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/afero/zipfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ZipEqualToFiles is a test utility that compares a zip buffer to a map of files,
// whose keys are filepaths and values are content of the files.
func ZipEqualToFiles(t *testing.T, buf []byte, files map[string]string) {
	r := bytes.NewReader(buf)
	zipR, err := zip.NewReader(r, r.Size())
	require.NoError(t, err)
	//FIXME: need to check for unexpected files in zip, afero.Walk for some reason does not work
	fs := zipfs.New(zipR)
	copy := files
	for name, content := range files {
		f, err := fs.Open(name)
		assert.NoError(t, err)
		buf, err := ioutil.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(t, content, string(buf))
		delete(copy, name)
	}
	assert.Zero(t, len(copy))
}

// CreateTestMemMapFs creates a memory fs from provided files
func CreateTestMemMapFs(t *testing.T, files map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()
	for name, content := range files {
		if name == "" {
			continue
		}
		name, err := filepath.Abs(name)
		require.NoError(t, err)
		file, err := fs.Create(name)
		require.NoError(t, err)
		_, err = file.Write([]byte(content))
		require.NoError(t, err)
	}
	return fs
}

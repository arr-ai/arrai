//nolint:unparam
package arrai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/syntax"
)

func TestCreateFile(t *testing.T) {
	t.Parallel()

	testAssertFiles(t,
		".",
		`
		{
			'file.txt': 'hello'
		}
		`,
		map[string]string{
			"file.txt": "hello",
		},
		afero.NewMemMapFs(),
	)
	testAssertFiles(t,
		".",
		`
		{
			'file.txt': (
				file: 'hello'
			)
		}
		`,
		map[string]string{
			"file.txt": "hello",
		},
		afero.NewMemMapFs(),
	)
	testAssertFiles(t,
		".",
		`
		{
			'files': {
				'file.txt': 'hello'
			}
		}
		`,
		map[string]string{
			"files/file.txt": "hello",
		},
		afero.NewMemMapFs(),
	)
	testAssertFiles(t,
		".",
		`
		{
			'files': (
				dir: {
					'file.txt': 'hello'
				}
			)
		}
		`,
		map[string]string{
			"files/file.txt": "hello",
		},
		afero.NewMemMapFs(),
	)
}

func TestCreateEmptyFile(t *testing.T) {
	t.Parallel()
	fs := afero.NewMemMapFs()
	testAssertFiles(t,
		".",
		`
		{
			'dir': (dir: {}),
			'file': (file: {}),
		}
		`,
		map[string]string{
			"file": "",
		},
		fs,
	)
	// the testAssertFiles ignore directories, it needs manual check
	fi, err := fs.Stat("dir")
	assert.NoError(t, err)
	assert.True(t, fi.IsDir())
}

func TestCreateDeepNestedFile(t *testing.T) {
	t.Parallel()

	testAssertFiles(t,
		".",
		`
		{
			'path': {
				'to': {
					'file.txt': 'hiii'
				}
			}
		}
		`,
		map[string]string{
			"path/to/file.txt": "hiii",
		},
		afero.NewMemMapFs(),
	)
}

func TestCreateMultipleFiles(t *testing.T) {
	t.Parallel()

	testAssertFiles(t,
		".",
		`
		{
			'path': {
				'to': {
					'file.txt': 'hiii'
				}
			},
			'hello.txt': 'hello'
		}
		`,
		map[string]string{
			"path/to/file.txt": "hiii",
			"hello.txt":        "hello",
		},
		afero.NewMemMapFs(),
	)
}

func TestRemoveIfExists(t *testing.T) {
	t.Parallel()

	testAssertFiles(t,
		".",
		`{'dir': {'to': {'replace': (ifExists: 'remove')}}}`,
		map[string]string{
			"dir/to/notreplace/dummy.txt": "dummy",
		},
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"dir/to/replace/file1.txt":    "replace",
			"dir/to/replace/file2.txt":    "RePlAcE",
			"dir/to/notreplace/dummy.txt": "dummy",
		}),
	)

	testAssertFiles(t,
		".",
		`{'dir': {'to': {'replace': {'file3.txt': (ifExists: 'remove')}}}}`,
		map[string]string{
			"dir/to/replace/file4.txt":    "retained",
			"dir/to/notreplace/dummy.txt": "dummy",
		},
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"dir/to/replace/file3.txt":    "removed",
			"dir/to/replace/file4.txt":    "retained",
			"dir/to/notreplace/dummy.txt": "dummy",
		}),
	)

	testAssertFiles(t,
		".",
		`{'dir': {'to': {'replace': {'file333.txt': (ifExists: 'remove')}}}}`,
		map[string]string{
			"dir/to/notreplace/dummy.txt": "dummy",
		},
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"dir/to/notreplace/dummy.txt": "dummy",
		}),
	)

	testAssertFilesError(t, errFileAndDirMustNotExist.Error(),
		".",
		`{'dir': {'to': {'replace': (
			ifExists: 'remove',
			dir: {'file': 'foo'},
		)}}}`,
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"dir/to/replace/file3.txt":    "rEPlaCEd",
			"dir/to/replace/file4.txt":    "retained",
			"dir/to/notreplace/dummy.txt": "dummy",
		}),
	)

	testAssertFilesError(t, errFileAndDirMustNotExist.Error(),
		".",
		`{'dir': {'to': {'replace': (
			ifExists: 'remove',
			file: 'foo',
		)}}}`,
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"dir/to/replace/file3.txt":    "rEPlaCEd",
			"dir/to/replace/file4.txt":    "retained",
			"dir/to/notreplace/dummy.txt": "dummy",
		}),
	)
}

func TestReplaceIfExists(t *testing.T) {
	t.Parallel()

	testAssertFiles(t,
		".",
		`
		{
			'dir': {
				'to': {
					'replace': (
						ifExists: 'replace',
						dir: {
							'file3.txt': 'rEPlaCEd'
						}
					)
				}
			}
		}
		`,
		map[string]string{
			"dir/to/replace/file3.txt":    "rEPlaCEd",
			"dir/to/notreplace/dummy.txt": "dummy",
		},
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"dir/to/replace/file1.txt":    "replace",
			"dir/to/replace/file2.txt":    "RePlAcE",
			"dir/to/notreplace/dummy.txt": "dummy",
		}),
	)

	testAssertFiles(t,
		".",
		`
		{
			'dir': {
				'to': {
					'replace': {
						'file3.txt': (
							ifExists: 'replace',
							file: 'replaced again'
						)
					}
				}
			}
		}
		`,
		map[string]string{
			"dir/to/replace/file3.txt":    "replaced again",
			"dir/to/notreplace/dummy.txt": "dummy",
		},
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"dir/to/replace/file3.txt":    "rEPlaCEd",
			"dir/to/notreplace/dummy.txt": "dummy",
		}),
	)
}

func TestIgnoreIfExists(t *testing.T) {
	t.Parallel()

	testAssertFiles(t,
		".",
		`
		{
			'ignore': (
				dir: {'file': 'not hello'},
				ifExists: 'ignore',
			),
			'dontignore': {
				'file': 'sike, not hello anymore'
			}
		}
		`,
		map[string]string{
			"ignore/file":     "hellooo",
			"dontignore/file": "sike, not hello anymore",
		},
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"ignore/file":     "hellooo",
			"dontignore/file": "hello again",
		}),
	)

	testAssertFiles(t,
		".",
		`
		{
			'ignore': {
				'file': (
					ifExists: 'ignore',
					file: 'oooleh'
				)
			},
			'dontignore': {
				'file': 'sike, not hello anymore'
			}
		}
		`,
		map[string]string{
			"ignore/file":     "hellooo",
			"dontignore/file": "sike, not hello anymore",
		},
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"ignore/file":     "hellooo",
			"dontignore/file": "hello again",
		}),
	)
}

func TestFailIfExists(t *testing.T) {
	t.Parallel()

	testAssertFilesError(t,
		fmt.Sprintf("%s: 'exists.txt' exists", ifExistsConfig),
		".",
		`
		{
			'exists.txt': (
				ifExists: 'fail',
				file: 'hiiii'
			)
		}
		`,
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"exists.txt": "hiiiiasd",
		}),
	)

	testAssertFilesError(t,
		fmt.Sprintf("%s: 'exists/exists.txt' exists", ifExistsConfig),
		".",
		`
		{
			'exists': {
				'exists.txt': (
					ifExists: 'fail',
					file: 'hiiii'
				)
			}
		}
		`,
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"exists/exists.txt": "hiiiiasd",
		}),
	)
}

func TestMergeIfExists(t *testing.T) {
	t.Parallel()

	testAssertFiles(t,
		".",
		`
		{
			'merge_this': (
				dir: {
					'file1.txt': 'number 1',
					'file3.txt': 'number 3',
				},
				ifExists: 'merge'
			)
		}
		`,
		map[string]string{
			"merge_this/file1.txt": "number 1",
			"merge_this/file2.txt": "number 2",
			"merge_this/file3.txt": "number 3",
		},
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"merge_this/file1.txt": "not number 1",
			"merge_this/file2.txt": "number 2",
		}),
	)

	testAssertFilesError(t,
		fmt.Sprintf("%s: '%s' config must not have '%s' field", ifExistsConfig, fileField, ifExistsMerge),
		".",
		`
		{
			'merge_this': (
				file: 'watch this',
				ifExists: 'merge'
			)
		}
		`,
		ctxfs.CreateTestMemMapFs(t, map[string]string{
			"merge_this": "watch that",
		}),
	)
}

func TestFail(t *testing.T) {
	t.Parallel()

	testAssertFilesError(t,
		fmt.Sprintf(
			"%s: value 'random value' is not valid value. It has to be one of %s",
			ifExistsConfig,
			strings.Join([]string{ifExistsMerge, ifExistsRemove, ifExistsReplace, ifExistsIgnore, ifExistsFail}, ", "),
		),
		".",
		`{'test': (
			ifExists: 'random value',
			'file': '123'
		)}`,
		afero.NewMemMapFs(),
	)

	testAssertFilesError(t,
		errFileOrDirMustExist.Error(),
		".",
		`
		{'test': (ifExists: 'replace')}`,
		afero.NewMemMapFs(),
	)
	testAssertFilesError(t,
		errFileOrDirMustExist.Error(),
		".",
		`{'test': (
			ifExists: 'replace',
			file: 'random',
			dir: {
				'random': 'random'
			}
		)}`,
		afero.NewMemMapFs(),
	)
}

func testAssertFiles(t *testing.T, dir, script string, files map[string]string, fs afero.Fs) {
	ctx := ctxfs.RuntimeFsOnto(context.Background(), fs)
	val, err := syntax.EvaluateExpr(ctx, ".", script)
	require.NoError(t, err)

	require.NoError(t, OutputValue(ctx, val, nil, fmt.Sprintf("dir:%s", dir)))

	abs := ctxfs.ToUnixPath(syntax.MustAbs(t, dir))
	assert.NoError(t,
		afero.Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			path = ctxfs.ToUnixPath(path)

			relPath := strings.Trim(strings.TrimPrefix(ctxfs.ToUnixPath(path), abs), "/")
			content, exists := files[relPath]
			assert.True(t, exists, fmt.Sprintf("unexpected file exists: path=%s, relPath=%s", path, relPath))

			f, err := afero.ReadFile(fs, path)
			assert.NoError(t, err)

			assert.Equal(t, content, string(f))

			if exists {
				delete(files, relPath)
			}
			return nil
		}),
	)
	assert.Zero(t, len(files), fmt.Sprintf("not all files exist: %v", files))
}

func testAssertFilesError(t *testing.T, expectedError, dir, script string, fs afero.Fs) {
	ctx := ctxfs.RuntimeFsOnto(context.Background(), fs)
	val, err := syntax.EvaluateExpr(ctx, ".", script)
	require.NoError(t, err)
	assert.EqualError(t, OutputValue(ctx, val, nil, fmt.Sprintf("dir:%s", dir)), expectedError)
}

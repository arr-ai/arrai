package test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

func RunTests(ctx context.Context, w io.Writer, path string) error {
	if path == "" || path == "." {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	files, err := getTestFiles(ctx, path)
	if err != nil {
		return err
	}

	for i := range files {
		err = runFile(ctx, &files[i])
		if err != nil {
			return err
		}
	}

	err = Report(w, files)
	return err
}

// Finds all *_test.arrai files in given path (recursively), reads them and returns a testFile array with them.
func getTestFiles(ctx context.Context, path string) ([]testFile, error) {
	var files []testFile
	fs := ctxfs.SourceFsFrom(ctx)

	err := afero.Walk(fs, path, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() {
			// Skip hidden dirs.
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			// Not a file, continue walking the directory tree.
			return nil
		}

		if strings.HasSuffix(path, "_test.arrai") == false { //nolint:gosimple
			return nil
		}

		bytes, readErr := afero.ReadFile(fs, path)
		if readErr != nil {
			return fmt.Errorf("failed reading test file '%s': %v", path, readErr)
		}

		files = append(files, testFile{path: path, source: string(bytes)})
		return nil
	})

	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no test files (ending in '_test.arrai') were found in path: %v", path)
	}
	return files, nil
}

// runFile runs all tests in testFile.source and fills testFile.results and testFile.wallTime.
func runFile(ctx context.Context, file *testFile) error {
	start := time.Now()
	result, err := syntax.EvaluateExpr(ctx, file.path, file.source)
	file.wallTime = time.Since(start)

	if err != nil {
		return fmt.Errorf("failed evaluating tests file '%s': %v", file.path, err)
	}

	file.results = make([]testResult, 0)
	ForeachLeaf(result, "", func(val rel.Value, path string) {
		result := testResult{
			name: path,
		}

		if isLiteralTrue(val) {
			result.outcome = Passed
		} else if isLiteralFalse(val) {
			result.outcome = Failed
			result.message = "Expected: true. Actual: false."
		} else {
			result.outcome = Invalid
			result.message = fmt.Sprintf("Could not determine test outcome due to non-boolean result of type '%T': %s",
				val, val.String())

			if _, ok := val.(rel.GenericSet); ok {
				result.message = fmt.Sprintf("Sets are not allowed as test containers. Please use tuples, " +
					"dictionaries or arrays.")
			}
		}

		file.results = append(file.results, result)
	})

	return nil
}

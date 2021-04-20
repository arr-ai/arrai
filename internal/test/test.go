package test

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

func RunTestsInPath(ctx context.Context, w io.Writer, path string) error {
	if path == "" {
		path = "."
	}

	files, err := FindTestFiles(ctx, w, path)
	if err != nil {
		return err
	}

	err = RunTests(ctx, &files)
	if err != nil {
		return err
	}

	err = Report(w, files)
	return err
}

func FindTestFiles(ctx context.Context, w io.Writer, path string) ([]TestFile, error) {
	var files []TestFile
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

			// Not a file, continue walking into it.
			return nil
		}

		if strings.HasSuffix(path, "_test.arrai") == false {
			return nil
		}

		bytes, readErr := afero.ReadFile(fs, path)
		if readErr != nil {
			fmt.Fprintf(w, "\nFailed reading test file %s\n", path)
			return readErr
		}

		file := TestFile{
			path:   path,
			source: string(bytes),
		}
		files = append(files, file)
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

// RunTestsInPath runs all tests in the subtree and adds the results to the []TestFile.
func RunTests(ctx context.Context, testFiles *[]TestFile) error {
	for _, file := range *testFiles {
		start := time.Now()
		result, err := syntax.EvaluateExpr(ctx, file.path, file.source)
		if err != nil {
			return fmt.Errorf("failed evaluating tests file '%s': %v", file.path, err)
		}
		file.wallTime = time.Since(start)

		file.results = make([]TestResult, 0)
		ForeachLeaf(result, "", func(val rel.Value, path string) {
			result := TestResult{
				name: path,
			}

			if isLiteralTrue(val) {
				result.outcome = Passed
			} else if isLiteralFalse(val) {
				result.outcome = Failed
				result.message = fmt.Sprint("Expected: true. Actual: false.")
			} else {
				result.outcome = Invalid
				result.message = fmt.Sprintf("Could not determine test outcome due to non-boolean result of type '%T': %s", val, val.String())

				if _, ok := val.(rel.GenericSet); ok {
					result.message = fmt.Sprintf("Sets are not allowed as test containers. Please use tuples, dictionaries or arrays.")
				}
			}

			file.results = append(file.results, result)
		})

		*testFiles = append(*testFiles, file)
	}

	return nil
}

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

// RunTests runs all arr.ai tests in a given path. It returns an error if the path is invalid, contains no test files or
// has invalid arr.ai code in any of them.
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

// getTestFiles finds all *_test.arrai files in given path (recursively), reads them and returns a TestFile array with
// them. It skips over hidden directories. It returns an error if any filesystem operation failed, or if no files were
// found.
func getTestFiles(ctx context.Context, path string) ([]TestFile, error) {
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

		files = append(files, TestFile{Path: path, Source: string(bytes)})
		return nil
	})

	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no test files (ending in '_test.arrai') were found in Path: %v", path)
	}
	return files, nil
}

// runFile runs all tests in TestFile.Source and fills TestFile.Results and TestFile.WallTime. It returns an error if
// the arr.ai code failed to evaluate.
func runFile(ctx context.Context, file *TestFile) error {
	expr, err := syntax.Compile(ctx, file.Path, file.Source)
	if err != nil {
		return fmt.Errorf("failed compiling tests file '%s': %v", file.Path, err)
	}

	start := time.Now()
	results, err := RunExpr(ctx, expr)
	file.WallTime = time.Since(start)
	file.Results = results

	if err != nil {
		return fmt.Errorf("failed evaluating tests file '%s': %v", file.Path, err)
	}

	return nil
}

// RunExpr runs all tests in the provided rel.Expr and returns a slice of TestResult. It returns an error if
// the arr.ai code failed to evaluate.
func RunExpr(ctx context.Context, expr rel.Expr) ([]TestResult, error) {
	result, err := expr.Eval(ctx, rel.Scope{})
	if err != nil {
		return nil, err
	}

	results := make([]TestResult, 0)
	ForeachLeaf(result, "", func(val rel.Value, path string) {
		result := TestResult{
			Name: path,
		}

		if isLiteralTrue(val) {
			result.Outcome = Passed
		} else if isLiteralFalse(val) {
			result.Outcome = Failed
			result.Message = "Expected: true. Actual: false."
		} else {
			result.Outcome = Invalid
			result.Message = fmt.Sprintf("Could not determine test Outcome due to non-boolean result of type %s: %s",
				rel.ValueTypeAsString(val), val.String())

			if _, ok := val.(rel.GenericSet); ok {
				result.Message = fmt.Sprintf("Sets are not allowed as test containers. Please use tuples, " +
					"dictionaries or arrays.")
			}
		}

		results = append(results, result)
	})

	return results, nil
}

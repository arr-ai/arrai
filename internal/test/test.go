package test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

// Test runs all tests in the subtree of path and returns the results.
func Test(ctx context.Context, w io.Writer, path string) ([]TestFile, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "_test.arrai") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no test files (filenames ending in '_test.arrai') were found in path: %v", path)
	}

	testFiles := make([]TestFile, 0)

	for _, file := range files {
		bytes, err := ctxfs.ReadFile(ctxfs.SourceFsFrom(ctx), file)
		if err != nil {
			fmt.Fprintf(w, "\nFailed reading test file %s\n", file)
			return nil, err
		}

		start := time.Now()
		result, err := syntax.EvaluateExpr(ctx, file, string(bytes))
		if err != nil {
			fmt.Fprintf(w, "\nFailed evaluating tests file %s\n", file)
			return nil, err
		}

		testFile := TestFile{
			name:     file,
			wallTime: time.Since(start),
			results:  make([]TestResult, 0),
		}

		ForeachLeaf(result, "<root>", func(val rel.Value, path string) {
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

			testFile.results = append(testFile.results, result)
		})

		testFiles = append(testFiles, testFile)
	}

	return testFiles, nil
}

type TestFile struct {
	name     string
	wallTime time.Duration
	results  []TestResult
}

type TestResult struct {
	name    string
	outcome TestOutcome
	message string
}

type TestOutcome int

const (
	Failed TestOutcome = iota
	Invalid
	Ignored
	Passed
)

func ForeachLeaf(val rel.Value, path string, leafAction func(val rel.Value, path string)) {
	path = strings.TrimPrefix(path, "<root>.")

	if isLiteralTrue(val) || isLiteralFalse(val) {
		leafAction(val, path)
		return
	}

	switch v := val.(type) {
	case rel.Array:
		for i, item := range v.Values() {
			ForeachLeaf(item, fmt.Sprintf("%s(%d)", path, i), leafAction)
		}
	case rel.Dict:
		for _, entry := range v.OrderedEntries() {
			key := entry.MustGet("@")
			keyStr := key.String()
			if _, ok := key.(rel.String); ok {
				keyStr = "'" + keyStr + "'"
			}
			ForeachLeaf(entry.MustGet(rel.DictValueAttr), fmt.Sprintf("%s(%s)", path, keyStr), leafAction)
		}
	case rel.Tuple:
		for e := v.Enumerator(); e.MoveNext(); {
			name, attr := e.Current()
			ForeachLeaf(attr, fmt.Sprintf("%s.%s", path, name), leafAction)
		}
	default:
		leafAction(val, path)
	}
}

var emptyTuple = rel.NewTuple()

func isLiteralTrue(val rel.Value) bool {
	v, ok := val.(rel.GenericSet)
	return ok && v.Count() == 1 && v.Has(emptyTuple)
}

func isLiteralFalse(val rel.Value) bool {
	v, ok := val.(rel.GenericSet)
	return ok && v.Count() == 0
}

package test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Report writes a formatted output of all the test files and their test results, and returns and error if the test run
// failed.
func Report(w io.Writer, testFiles []File) error {
	stats := calcStats(testFiles)

	for _, testFile := range testFiles {
		reportFile(w, testFile, stats.maxNameLen)
	}

	reportStats(w, stats)

	if stats.runFailed {
		return fmt.Errorf("test run "+color, 41, "FAILED")
	}

	return nil
}

const color = "\033[38;5;255;%d;1m%s\033[0m"

// reportFile writes a formatted output of a file and all its test results, ordered by outcome. The maxName parameter
// is used to aligned the results, and should contain the length of the longest Result.Name inside File.Results.
func reportFile(w io.Writer, testFile File, maxName int) {
	results := testFile.Results

	// Sort test results by outcome then by test name
	sort.Slice(results, func(i, j int) bool {
		left, right := results[i], results[j]
		if left.Outcome == right.Outcome {
			if left.Name == right.Name {
				return true
			}
			return left.Name < right.Name
		}
		return left.Outcome > right.Outcome
	})

	message.NewPrinter(language.English).
		Fprintf(w, "\n=======  %s (%dms)\n", relPath(testFile.Path), testFile.WallTime.Milliseconds())
	for _, result := range results {
		reportTest(w, result, maxName)
	}
}

// relPath makes a best-effort attempt to compute the relative path of the given absolute path. If it fails, it returns
// the absolute path untouched.
func relPath(absPath string) string {
	cwd, wdErr := os.Getwd()
	relPath, relErr := filepath.Rel(cwd, absPath)

	if wdErr != nil || relErr != nil {
		return absPath
	}

	return relPath
}

// reportTest writes a formatted output of a single test result (PASS/FAIL/SKIP/??) with the optional included Message.
func reportTest(w io.Writer, test Result, maxName int) {
	switch test.Outcome {
	case Failed:
		fmt.Fprintf(w, color, 41, "FAIL")
	case Invalid:
		fmt.Fprintf(w, color, 31, " ?? ")
	case Ignored:
		fmt.Fprintf(w, color, 93, "SKIP")
	case Passed:
		fmt.Fprintf(w, color, 32, "PASS")
	}

	fmt.Fprintf(w, "  %s\n", test.Name+strings.Repeat(" ", maxName-utf8.RuneCountInString(test.Name)))

	if test.Message != "" {
		fmt.Fprintf(w, "      %s\n", strings.ReplaceAll(test.Message, "\n", "\n      "))
	}
}

// reportStats writes a formatted single line representation of the aggregated test statistics.
func reportStats(w io.Writer, stats testStats) {
	p := message.NewPrinter(language.English)
	p.Fprintf(w, "\n=======  Summary\n")

	if stats.failed > 0 {
		p.Fprintf(w, "%d failed, ", stats.failed)
	}

	if stats.invalid > 0 {
		p.Fprintf(w, "%d invalid, ", stats.invalid)
	}

	if stats.ignored > 0 {
		p.Fprintf(w, "%d ignored, ", stats.ignored)
	}

	p.Fprintf(w, "%d passed of %d total tests. Took %dms.\n", stats.passed, stats.total, stats.wallTime.Milliseconds())
}

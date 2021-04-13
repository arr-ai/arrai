package test

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

func Report(w io.Writer, testFiles []TestFile) error {
	stats := calcStats(testFiles)

	for _, testFile := range testFiles {
		reportFile(w, testFile, stats.maxName)
	}

	reportStats(w, stats)

	if stats.runFailed {
		return fmt.Errorf("test run \u001B[38;5;255;41;1mFAILED\u001B[0m")
	}

	return nil
}

func reportFile(w io.Writer, testFile TestFile, maxName int) {
	results := testFile.results

	// Sort test results by outcome then by test name
	sort.Slice(results, func(i, j int) bool {
		left, right := results[i], results[j]
		if left.outcome == right.outcome {
			if left.name == right.name {
				return true
			}
			return left.name < right.name
		}
		return left.outcome > right.outcome
	})

	message.NewPrinter(language.English).Fprintf(w, "\n=======  %s (%dms)\n", testFile.name, testFile.wallTime.Milliseconds())
	for _, result := range results {
		reportTest(w, result, maxName)
	}
}

func calcStats(testFiles []TestFile) testStats {
	var stats testStats

	for _, testFile := range testFiles {
		stats.wallTime += testFile.wallTime

		for _, result := range testFile.results {
			if count := utf8.RuneCountInString(result.name); count > stats.maxName {
				stats.maxName = count
			}

			if count := utf8.RuneCountInString(testFile.name); count > stats.maxName {
				stats.maxFile = count
			}

			stats.total++

			switch result.outcome {
			case Invalid:
				stats.invalid++
			case Passed:
				stats.passed++
			case Ignored:
				stats.ignored++
			case Failed:
				stats.failed++
			}

			stats.runFailed = stats.failed > 0 || stats.invalid > 0
		}
	}

	return stats
}

type testStats struct {
	maxName int
	maxFile int

	total   int
	invalid int
	passed  int
	ignored int
	failed  int

	runFailed bool
	wallTime  time.Duration
}

func reportTest(w io.Writer, result TestResult, maxName int) {
	const color = "\033[38;5;255;%d;1m%s\033[0m"

	switch result.outcome {
	case Failed:
		fmt.Fprintf(w, color, 41, "FAIL")
	case Invalid:
		fmt.Fprintf(w, color, 31, " ?? ")
	case Ignored:
		fmt.Fprintf(w, color, 93, "SKIP")
	case Passed:
		fmt.Fprintf(w, color, 32, "PASS")
	}

	fmt.Fprintf(w, "  %s\n", result.name+strings.Repeat(" ", maxName-utf8.RuneCountInString(result.name)))

	if result.message != "" {
		fmt.Fprintf(w, "      %s\n", strings.ReplaceAll(result.message, "\n", "\n      "))
	}

}

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

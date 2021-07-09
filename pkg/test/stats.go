package test

import (
	"time"
	"unicode/utf8"
)

// calcStats calculates aggregate statistics for all test Results in the provided test files. Most importantly, it
// determines if the test run has succeeded or failed. A test run fails if there are any testResults with a 'failed' or
// 'invalid' Outcome.
func calcStats(testFiles []TestFile) testStats {
	var stats testStats

	for _, testFile := range testFiles {
		stats.wallTime += testFile.WallTime

		for _, result := range testFile.Results {
			if count := utf8.RuneCountInString(result.Name); count > stats.maxNameLen {
				stats.maxNameLen = count
			}

			if count := utf8.RuneCountInString(testFile.Path); count > stats.maxFileLen {
				stats.maxFileLen = count
			}

			stats.total++

			switch result.Outcome {
			case Invalid:
				stats.invalid++
			case Passed:
				stats.passed++
			case Ignored:
				stats.ignored++
			case Failed:
				stats.failed++
			}
		}
	}

	stats.runFailed = stats.failed > 0 || stats.invalid > 0

	return stats
}

type testStats struct {
	runFailed  bool
	wallTime   time.Duration
	maxNameLen int
	maxFileLen int

	total   int
	invalid int
	passed  int
	ignored int
	failed  int
}

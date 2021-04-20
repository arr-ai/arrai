package test

import (
	"time"
	"unicode/utf8"
)

func calcStats(testFiles []TestFile) testStats {
	var stats testStats

	for _, testFile := range testFiles {
		stats.wallTime += testFile.wallTime

		for _, result := range testFile.results {
			if count := utf8.RuneCountInString(result.name); count > stats.maxNameLen {
				stats.maxNameLen = count
			}

			if count := utf8.RuneCountInString(testFile.path); count > stats.maxFileLen {
				stats.maxFileLen = count
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

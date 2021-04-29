//nolint:lll
package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var passFile = testFile{
	path:     "file1",
	source:   "src1",
	wallTime: 100,
	results: []testResult{
		{name: "Test1", outcome: Passed},
		{name: "Test22", outcome: Passed},
		{name: "Test333", outcome: Passed},
		{name: "Test4444", outcome: Passed},
	},
}
var failFile = testFile{
	path:     "file22",
	source:   "src2",
	wallTime: 50,
	results: []testResult{
		{name: "Test1", outcome: Passed},
		{name: "Test2", outcome: Failed},
		{name: "Test3", outcome: Ignored},
	},
}
var invalidFile = testFile{
	path:     "file3",
	source:   "src3",
	wallTime: 3,
	results: []testResult{
		{name: "Test1", outcome: Invalid},
	},
}

func TestCalcStats_AllPass(t *testing.T) {
	t.Parallel()
	require.Equal(t,
		testStats{runFailed: false, wallTime: 100, maxNameLen: 8, maxFileLen: 5, total: 4, invalid: 0, passed: 4, ignored: 0, failed: 0},
		calcStats([]testFile{passFile}))
}

func TestCalcStats_Mixed(t *testing.T) {
	t.Parallel()
	require.Equal(t,
		testStats{runFailed: true, wallTime: 153, maxNameLen: 8, maxFileLen: 6, total: 8, invalid: 1, passed: 5, ignored: 1, failed: 1},
		calcStats([]testFile{passFile, failFile, invalidFile}))
}

func TestCalcStats_Failed(t *testing.T) {
	t.Parallel()
	require.Equal(t,
		testStats{runFailed: true, wallTime: 50, maxNameLen: 5, maxFileLen: 6, total: 3, invalid: 0, passed: 1, ignored: 1, failed: 1},
		calcStats([]testFile{failFile}))
}

func TestCalcStats_Invalid(t *testing.T) {
	t.Parallel()
	require.Equal(t,
		testStats{runFailed: true, wallTime: 3, maxNameLen: 5, maxFileLen: 5, total: 1, invalid: 1, passed: 0, ignored: 0, failed: 0},
		calcStats([]testFile{invalidFile}))
}

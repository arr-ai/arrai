//nolint:lll
package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var passFile = File{
	Path:     "file1",
	Source:   "src1",
	WallTime: 100,
	Results: []Result{
		{Name: "Test1", Outcome: Passed},
		{Name: "Test22", Outcome: Passed},
		{Name: "Test333", Outcome: Passed},
		{Name: "Test4444", Outcome: Passed},
	},
}
var failFile = File{
	Path:     "file22",
	Source:   "src2",
	WallTime: 50,
	Results: []Result{
		{Name: "Test1", Outcome: Passed},
		{Name: "Test2", Outcome: Failed},
		{Name: "Test3", Outcome: Ignored},
	},
}
var invalidFile = File{
	Path:     "file3",
	Source:   "src3",
	WallTime: 3,
	Results: []Result{
		{Name: "Test1", Outcome: Invalid},
	},
}

func TestCalcStats_AllPass(t *testing.T) {
	t.Parallel()
	require.Equal(t,
		testStats{runFailed: false, wallTime: 100, maxNameLen: 8, maxFileLen: 5, total: 4, invalid: 0, passed: 4, ignored: 0, failed: 0},
		calcStats([]File{passFile}))
}

func TestCalcStats_Mixed(t *testing.T) {
	t.Parallel()
	require.Equal(t,
		testStats{runFailed: true, wallTime: 153, maxNameLen: 8, maxFileLen: 6, total: 8, invalid: 1, passed: 5, ignored: 1, failed: 1},
		calcStats([]File{passFile, failFile, invalidFile}))
}

func TestCalcStats_Failed(t *testing.T) {
	t.Parallel()
	require.Equal(t,
		testStats{runFailed: true, wallTime: 50, maxNameLen: 5, maxFileLen: 6, total: 3, invalid: 0, passed: 1, ignored: 1, failed: 1},
		calcStats([]File{failFile}))
}

func TestCalcStats_Invalid(t *testing.T) {
	t.Parallel()
	require.Equal(t,
		testStats{runFailed: true, wallTime: 3, maxNameLen: 5, maxFileLen: 5, total: 1, invalid: 1, passed: 0, ignored: 0, failed: 0},
		calcStats([]File{invalidFile}))
}

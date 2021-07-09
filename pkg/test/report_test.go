package test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ===== report

func TestReport_TwoFilesPass(t *testing.T) {
	t.Parallel()

	err := reportShouldContain(t, generateFiles(),
		"Path/to/file1.arrai", "1,234ms", "PASS", "testA",
		"Path/to/file2.arrai", "4,321ms", "PASS", "testB",
		"2 passed", "2 total", "5,555ms")
	require.NoError(t, err)
}

func TestReport_TwoFilesFail(t *testing.T) {
	t.Parallel()

	files := generateFiles()
	files[1].Results[0] = TestResult{Name: "testB", Outcome: Failed, Message: "Uh oh"}
	err := reportShouldContain(t, files,
		"Path/to/file1.arrai", "1,234ms", "PASS", "testA",
		"Path/to/file2.arrai", "4,321ms", "FAIL", "testB",
		"1 failed", "1 passed", "2 total", "5,555ms")
	require.Error(t, err)
}

// ===== reportFile

func TestReportFile_WithTest(t *testing.T) {
	t.Parallel()
	generateFiles()[0].reportShouldContain(t, "Path/to/file1.arrai", "1,234ms", "PASS", "testA")
}

func TestReportFile_SortsCorrectly(t *testing.T) {
	t.Parallel()
	file := generateFiles()[0]
	file.Results = []TestResult{
		{Name: "failedB", Outcome: Failed},
		{Name: "failedA", Outcome: Failed},
		{Name: "ignoredB", Outcome: Ignored},
		{Name: "ignoredA", Outcome: Ignored},
		{Name: "invalidB", Outcome: Invalid},
		{Name: "invalidA", Outcome: Invalid},
		{Name: "passedA", Outcome: Passed},
		{Name: "passedA", Outcome: Passed, Message: "Duplicate"},
	}
	file.reportShouldContain(t, "passedA", "Duplicate", "ignoredA", "ignoredB", "invalidA", "invalidB",
		"failedA", "failedB")
}

// ===== reportTest

func TestReportTest_Passed(t *testing.T) {
	t.Parallel()
	TestResult{Name: "test1", Outcome: Passed}.
		reportShouldContain(t, "PASS", "test1")
}

func TestReportTest_Failed(t *testing.T) {
	t.Parallel()
	TestResult{Name: "test2", Outcome: Failed, Message: "Failed!"}.
		reportShouldContain(t, "FAIL", "test2", "Failed!")
}

func TestReportTest_Invalid(t *testing.T) {
	t.Parallel()
	TestResult{Name: "test3", Outcome: Invalid, Message: "Invalid!"}.
		reportShouldContain(t, "??", "test3", "Invalid!")
}

func TestReportTest_Ignored(t *testing.T) {
	t.Parallel()
	TestResult{Name: "test4", Outcome: Ignored}.
		reportShouldContain(t, "SKIP", "test4")
}

// ===== reportStats

func TestReportStats_Passed(t *testing.T) {
	t.Parallel()
	testStats{wallTime: time.Millisecond * 11111, total: 2, passed: 2}.
		reportShouldContain(t, "2 passed", "2 total", "11,111ms")
}

func TestReportStats_Failed(t *testing.T) {
	t.Parallel()
	testStats{wallTime: time.Millisecond * 888, total: 4, passed: 1, failed: 1, invalid: 1, ignored: 1}.
		reportShouldContain(t, "1 failed", "1 invalid", "1 ignored", "1 passed", "4 total", "888ms")
}

// ===== Helpers

func reportShouldContain(t *testing.T, files []TestFile, keyPhrases ...string) error {
	buf := bytes.Buffer{}
	err := Report(&buf, files)
	shouldContain(t, buf.String(), keyPhrases...)
	return err
}

func (file TestFile) reportShouldContain(t *testing.T, keyPhrases ...string) {
	buf := bytes.Buffer{}
	reportFile(&buf, file, 10)
	shouldContain(t, buf.String(), keyPhrases...)
}

func (test TestResult) reportShouldContain(t *testing.T, keyPhrases ...string) {
	buf := bytes.Buffer{}
	reportTest(&buf, test, 10)
	shouldContain(t, buf.String(), keyPhrases...)
}

func (stats testStats) reportShouldContain(t *testing.T, keyPhrases ...string) {
	buf := bytes.Buffer{}
	reportStats(&buf, stats)
	shouldContain(t, buf.String(), keyPhrases...)
}

// shouldContain verifies that the provided keyPhrases exist in str in the order specified.
func shouldContain(t *testing.T, str string, keyPhrases ...string) {
	for _, phrase := range keyPhrases {
		index := strings.Index(str, phrase)
		if index == -1 {
			require.Failf(t, "shouldContain() failed",
				"Key phrase '%s' was not found in remaining unmatched string: %s", strconv.Quote(phrase), strconv.Quote(str))
		}
		str = str[index+len(phrase):]
	}
}

func generateFiles() []TestFile {
	return []TestFile{
		{
			Path: "/Path/to/file1.arrai", WallTime: 1234 * time.Millisecond,
			Results: []TestResult{{Name: "testA", Outcome: Passed}},
		},
		{
			Path: "/Path/to/file2.arrai", WallTime: 4321 * time.Millisecond,
			Results: []TestResult{{Name: "testB", Outcome: Passed}},
		},
	}
}

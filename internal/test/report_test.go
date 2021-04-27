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
		"/path/to/file1.arrai", "1,234ms", "PASS", "testA",
		"/path/to/file2.arrai", "4,321ms", "PASS", "testB",
		"2 passed", "2 total", "5,555ms")
	require.NoError(t, err)
}

func TestReport_TwoFilesFail(t *testing.T) {
	t.Parallel()

	files := generateFiles()
	files[1].results[0] = testResult{name: "testB", outcome: Failed, message: "Uh oh"}
	err := reportShouldContain(t, files,
		"/path/to/file1.arrai", "1,234ms", "PASS", "testA",
		"/path/to/file2.arrai", "4,321ms", "FAIL", "testB",
		"1 failed", "1 passed", "2 total", "5,555ms")
	require.Error(t, err)
}

// ===== reportFile

func TestReportFile_WithTest(t *testing.T) {
	t.Parallel()
	generateFiles()[0].reportShouldContain(t, "/path/to/file1.arrai", "1,234ms", "PASS", "testA")
}

func TestReportFile_SortsCorrectly(t *testing.T) {
	t.Parallel()
	file := generateFiles()[0]
	file.results = []testResult{
		{name: "failedB", outcome: Failed},
		{name: "failedA", outcome: Failed},
		{name: "ignoredB", outcome: Ignored},
		{name: "ignoredA", outcome: Ignored},
		{name: "invalidB", outcome: Invalid},
		{name: "invalidA", outcome: Invalid},
		{name: "passedA", outcome: Passed},
		{name: "passedA", outcome: Passed, message: "Duplicate"},
	}
	file.reportShouldContain(t, "passedA", "Duplicate", "ignoredA", "ignoredB", "invalidA", "invalidB",
		"failedA", "failedB")
}

// ===== reportTest

func TestReportTest_Passed(t *testing.T) {
	t.Parallel()
	testResult{name: "test1", outcome: Passed}.
		reportShouldContain(t, "PASS", "test1")
}

func TestReportTest_Failed(t *testing.T) {
	t.Parallel()
	testResult{name: "test2", outcome: Failed, message: "Failed!"}.
		reportShouldContain(t, "FAIL", "test2", "Failed!")
}

func TestReportTest_Invalid(t *testing.T) {
	t.Parallel()
	testResult{name: "test3", outcome: Invalid, message: "Invalid!"}.
		reportShouldContain(t, "??", "test3", "Invalid!")
}

func TestReportTest_Ignored(t *testing.T) {
	t.Parallel()
	testResult{name: "test4", outcome: Ignored}.
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

func reportShouldContain(t *testing.T, files []testFile, keyPhrases ...string) error {
	buf := bytes.Buffer{}
	err := Report(&buf, files)
	shouldContain(t, buf.String(), keyPhrases...)
	return err
}

func (file testFile) reportShouldContain(t *testing.T, keyPhrases ...string) {
	buf := bytes.Buffer{}
	reportFile(&buf, file, 10)
	shouldContain(t, buf.String(), keyPhrases...)
}

func (test testResult) reportShouldContain(t *testing.T, keyPhrases ...string) {
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

func generateFiles() []testFile {
	return []testFile{
		{
			path: "/path/to/file1.arrai", wallTime: 1234 * time.Millisecond,
			results: []testResult{{name: "testA", outcome: Passed}},
		},
		{
			path: "/path/to/file2.arrai", wallTime: 4321 * time.Millisecond,
			results: []testResult{{name: "testB", outcome: Passed}},
		},
	}
}

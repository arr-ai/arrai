// +build !windows

package test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/ctxrootcache"
)

func TestRunTests_Pass(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{"/test/1_test.arrai": "(test1: 1 = 1)"}
	buf := &bytes.Buffer{}
	err := RunTests(withFs(t, fsContent), buf, "/test")

	require.NoError(t, err)
	shouldContain(t, buf.String(), "PASS", "test1")
}

func TestRunTests_EmptyPath(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{"/test/1_test.arrai": "(test1: 1 = 1)"}
	buf := &bytes.Buffer{}
	err := RunTests(withFs(t, fsContent), buf, "")

	require.NoError(t, err)
	shouldContain(t, buf.String(), "PASS", "test1")
}

func TestRunTests_Fail(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{"/test/1_test.arrai": "(test1: 1 = 1, test2: 2 = 3)"}
	buf := &bytes.Buffer{}
	err := RunTests(withFs(t, fsContent), buf, "/test")

	require.Error(t, err)
	shouldContain(t, buf.String(), "PASS", "test1", "FAIL", "test2")
}

func TestRunTests_Invalid(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{"/test/1_test.arrai": "(test1: 1 = 1, test2: 2)"}
	buf := &bytes.Buffer{}
	err := RunTests(withFs(t, fsContent), buf, "/test")

	require.Error(t, err)
	shouldContain(t, buf.String(), "PASS", "test1", "??", "test2")
}

func TestRunTests_BadArrai(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{"/test/1_test.arrai": "x"}
	buf := &bytes.Buffer{}
	err := RunTests(withFs(t, fsContent), buf, "/test")

	require.Error(t, err)
}

func TestRunTests_InvalidDir(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{}
	buf := &bytes.Buffer{}
	err := RunTests(withFs(t, fsContent), buf, "/nowhere")

	require.Error(t, err)
}

func TestGetTestFiles_NoFiles(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{"/test/placeholder.txt": ""}
	_, err := getTestFiles(withFs(t, fsContent), "/test")
	require.Error(t, err)
}

func TestGetTestFiles_InvalidDir(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{"/test/placeholder.txt": ""}
	_, err := getTestFiles(withFs(t, fsContent), "/nowhere")
	require.Error(t, err)
}

func TestGetTestFiles_OneFile(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{"/test/1_test.arrai": "source1"}
	files, err := getTestFiles(withFs(t, fsContent), "/test")
	require.NoError(t, err)
	require.Equal(t, []File{{Path: "/test/1_test.arrai", Source: "source1"}}, files)
}

func TestGetTestFiles_PathIsFile(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{
		"/test/1_test.arrai": "source1",
		"/test/2_test.arrai": "source2",
	}
	files, err := getTestFiles(withFs(t, fsContent), "/test/1_test.arrai")
	require.NoError(t, err)
	require.Equal(t, []File{{Path: "/test/1_test.arrai", Source: "source1"}}, files)
}

func TestGetTestFiles_NestedDir(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{
		"/test/1_test.arrai":                "source1",
		"/test/must/go/deeper/2_test.arrai": "source2",
	}
	files, err := getTestFiles(withFs(t, fsContent), "/test")
	require.NoError(t, err)
	require.Equal(t, []File{
		{Path: "/test/1_test.arrai", Source: "source1"},
		{Path: "/test/must/go/deeper/2_test.arrai", Source: "source2"},
	}, files)
}

func TestGetTestFiles_SkipHiddenDir(t *testing.T) {
	t.Parallel()

	fsContent := map[string]string{
		"/test/1_test.arrai":                 "source1",
		"/test/.must/go/deeper/2_test.arrai": "source2",
		"/test/must/.go/deeper/3_test.arrai": "source3",
		"/test/must/go/.deeper/4_test.arrai": "source4",
	}
	files, err := getTestFiles(withFs(t, fsContent), "/test")
	require.NoError(t, err)
	require.Equal(t, []File{{Path: "/test/1_test.arrai", Source: "source1"}}, files)
}

func withFs(t *testing.T, files map[string]string) context.Context {
	err := os.Chdir("/")
	require.NoError(t, err)
	fs := ctxfs.CreateTestMemMapFs(t, files)
	ctx := ctxfs.SourceFsOnto(context.Background(), fs)
	return ctxrootcache.WithRootCache(ctx)
}

func TestRunFile_InvalidArrai(t *testing.T) {
	t.Parallel()

	file := File{Source: "invalid arr.ai code"}
	err := runFile(context.Background(), &file)
	require.Error(t, err)
}

func TestRunFile_AssertFails(t *testing.T) {
	t.Parallel()

	file := File{Source: "//test.assert.equal(1, 2)"}
	err := runFile(context.Background(), &file)
	require.Error(t, err)
}

func TestRunFile_TwoPass(t *testing.T) {
	t.Parallel()

	file := File{Source: "(test1: 1 = 1, test2: //test.assert.equal(2, 2))"}
	err := runFile(context.Background(), &file)
	require.NoError(t, err)
	require.NotZero(t, file.WallTime)
	require.ElementsMatch(t, file.Results, []Result{
		{Name: "test1", Outcome: Passed},
		{Name: "test2", Outcome: Passed}})
}

func TestRunFile_OneFailOnePass(t *testing.T) {
	t.Parallel()

	file := File{Source: "(test1: 1 < 1, test2: 5 < 7)"}
	err := runFile(context.Background(), &file)
	require.NoError(t, err)
	require.NotZero(t, file.WallTime)
	require.ElementsMatch(t, file.Results, []Result{
		{Name: "test1", Outcome: Failed, Message: "Expected: true. Actual: false."},
		{Name: "test2", Outcome: Passed}})
}

func TestRunFile_OneInvalidOnePass(t *testing.T) {
	t.Parallel()

	file := File{Source: "(test1: 1, test2: 5 < 7)"}
	err := runFile(context.Background(), &file)
	require.NoError(t, err)
	require.NotZero(t, file.WallTime)
	require.ElementsMatch(t, file.Results, []Result{
		{Name: "test1", Outcome: Invalid,
			Message: "Could not determine test Outcome due to non-boolean result of type number: 1"},
		{Name: "test2", Outcome: Passed}})
}

func TestRunFile_TestInSet(t *testing.T) {
	t.Parallel()

	file := File{
		Path:   "some_test.arrai",
		Source: "(test1: 1 = 1, category1: { 5 < 7 })",
	}
	err := runFile(context.Background(), &file)
	require.NoError(t, err)
	require.NotZero(t, file.WallTime)
	require.ElementsMatch(t, file.Results, []Result{
		{Name: "test1", Outcome: Passed},
		{Name: "category1", Outcome: Invalid, Message: "Sets are not allowed as test containers. " +
			"Please use tuples, dictionaries or arrays."}})
}

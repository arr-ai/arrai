package syntax

import (
	"bytes"
	"testing"
)

func TestStdOsStdin(t *testing.T) {
	// Not parallelisable
	stdin := stdOsStdinVar.reader
	defer func() { stdOsStdinVar.reset(stdin) }()

	// Access twice to ensure caching behaves properly.
	stdOsStdinVar.reset(bytes.NewBuffer([]byte("")))
	AssertCodesEvalToSameValue(t, `{}`, `//os.stdin`)
	AssertCodesEvalToSameValue(t, `{}`, `//os.stdin`)

	stdOsStdinVar.reset(bytes.NewBuffer([]byte("abc")))
	AssertCodesEvalToSameValue(t, `<<97, 98, 99>>`, `//os.stdin`)
	AssertCodesEvalToSameValue(t, `<<97, 98, 99>>`, `//os.stdin`)
}

func TestStdOsExists(t *testing.T) {
	AssertCodesEvalToSameValue(t, `true`, `//os.exists('std_os_test.go')`)
	AssertCodesEvalToSameValue(t, `false`, `//os.exists('doesntexist.anywhere')`)
}

func TestStdOsTree(t *testing.T) {
	t.Parallel()

	// modTime is set to -1 to be non-deterministic.
	AssertCodesEvalToSameValue(t, `{
		(name: "std_os_test", path: "std_os_test", isDir: true, size: 160, modTime: -1),
		(name: ".empty", path: "std_os_test/.empty", isDir: false, size: 0, modTime: -1),
		(name: "README.md", path: "std_os_test/README.md", isDir: false, size: 84, modTime: -1),
		(name: "no files", path: "std_os_test/no files", isDir: true, size: 96, modTime: -1),
		(name: "full", path: "std_os_test/no files/full", isDir: true, size: 160, modTime: -1),
		(name: "README.md", path: "std_os_test/no files/full/README.md", isDir: false, size: 73, modTime: -1),
 		(name: "root.ln", path: "std_os_test/no files/full/root.ln", isDir: false, size: 18, modTime: -1),
		(name: "empty", path: "std_os_test/no files/full/empty", isDir: true, size: 64, modTime: -1),
	}`, `//os.tree('std_os_test') => . +> (modTime: -1)`)

	AssertCodesEvalToSameValue(t, `{
		(name: "empty", path: "std_os_test/no files/full/empty/", isDir: true, size: 64, modTime: -1),
	}`, `//os.tree('std_os_test/no files/full/empty/') => . +> (modTime: -1)`)

	AssertCodesEvalToSameValue(t, `{'.'}`, `//os.tree('.') => .path where . = '.'`)

	AssertCodesEvalToSameValue(t, `{
		(name: "README.md", path: "std_os_test/README.md", isDir: false, size: 84, modTime: -1),
	}`, `//os.tree('std_os_test/README.md') => . +> (modTime: -1)`)

	AssertCodeErrors(t, ``, `//os.tree(['std_os_test'])`)
	AssertCodeErrors(t, ``, `//os.tree('doesntexist')`)
}

func TestStdOsIsATty(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `false`, `//os.isatty(0)`)
	AssertCodesEvalToSameValue(t, `false`, `//os.isatty(1)`)

	AssertCodeErrors(t, "isatty arg must be a number, not rel.String", `//os.isatty("0")`)
	AssertCodeErrors(t, "isatty arg must be an integer, not 0.1", `//os.isatty(0.1)`)
	AssertCodeErrors(t, "isatty not implemented for 2", `//os.isatty(2)`)
	AssertCodeErrors(t, "isatty not implemented for -1", `//os.isatty(-1)`)
}

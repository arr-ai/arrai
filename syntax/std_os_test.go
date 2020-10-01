package syntax

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/tools"
)

func TestStdOsArgs(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[]`, `//os.args`)

	tools.Arguments = []string{"a", "b"}

	AssertCodesEvalToSameValue(t, `['a', 'b']`, `//os.args`)

	tools.Arguments = []string{"a", "b", "c"}

	AssertCodesEvalToSameValue(t, `['a', 'b', 'c']`, `//os.args`)
}

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
	t.Parallel()

	AssertCodesEvalToSameValue(t, `true`, `//os.exists('std_os_test.go')`)
	AssertCodesEvalToSameValue(t, `false`, `//os.exists('doesntexist.anywhere')`)
}

func TestStdOsTree(t *testing.T) {
	t.Parallel()

	// size and mod_time are non-deterministic, so evaluate some predicate of them instead.
	predx := `. +> (mod_time: .mod_time > 0, size: .size > 0)`

	AssertCodesEvalToSameValue(t, fixWindows(`{
		(name: "std_os_test", path: "std_os_test", is_dir: true, size: true, mod_time: true),
		(name: ".empty", path: "std_os_test/.empty", is_dir: false, size: false, mod_time: true),
		(name: "README.md", path: "std_os_test/README.md", is_dir: false, size: true, mod_time: true),
		(name: "no files", path: "std_os_test/no files", is_dir: true, size: true, mod_time: true),
		(name: "full", path: "std_os_test/no files/full", is_dir: true, size: true, mod_time: true),
		(name: "README.md", path: "std_os_test/no files/full/README.md", is_dir: false, size: true, mod_time: true),
 		(name: "root.ln", path: "std_os_test/no files/full/root.ln", is_dir: false, size: true, mod_time: true),
	}`), fmt.Sprintf(`//os.tree('std_os_test') => %s`, predx))

	AssertCodesEvalToSameValue(t, `{'.'}`, `//os.tree('.') => .path where . = '.'`)

	AssertCodesEvalToSameValue(t, fixWindows(`{
		(name: "README.md", path: "std_os_test/README.md", is_dir: false, size: true, mod_time: true),
	}`), fmt.Sprintf(`//os.tree('std_os_test/README.md') => %s`, predx))

	AssertCodeErrors(t, ``, `//os.tree(['std_os_test'])`)
	AssertCodeErrors(t, ``, `//os.tree('doesntexist')`)
}

func TestStdOsIsATty(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `false`, `//os.isatty(0)`)
	AssertCodesEvalToSameValue(t, `false`, `//os.isatty(1)`)

	AssertCodeErrors(t, "isatty arg must be a number, not string", `//os.isatty("0")`)
	AssertCodeErrors(t, "isatty arg must be an integer, not 0.1", `//os.isatty(0.1)`)
	AssertCodeErrors(t, "isatty not implemented for 2", `//os.isatty(2)`)
	AssertCodeErrors(t, "isatty not implemented for -1", `//os.isatty(-1)`)
}

// fixWindows replaces all /s with \s if running on Windows.
func fixWindows(code string) string {
	if runtime.GOOS != "windows" {
		return code
	}

	// Windows uses \\ as the path separator.
	code = strings.ReplaceAll(code, `/`, `\\`)
	// Windows directories have zero size.
	code = strings.ReplaceAll(code, `is_dir: true, size: true`, `is_dir: true, size: false`)
	// Symlinks on Windows have zero size.
	code = strings.ReplaceAll(code, `.ln", is_dir: false, size: true`, `.ln", is_dir: false, size: false`)
	return code
}

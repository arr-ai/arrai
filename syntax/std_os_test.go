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

func TestStdOsIsATty(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `false`, `//os.isatty(0)`)
	AssertCodesEvalToSameValue(t, `false`, `//os.isatty(1)`)

	AssertCodeErrors(t, "isatty arg must be a number, not string", `//os.isatty("0")`)
	AssertCodeErrors(t, "isatty arg must be an integer, not 0.1", `//os.isatty(0.1)`)
	AssertCodeErrors(t, "isatty not implemented for 2", `//os.isatty(2)`)
	AssertCodeErrors(t, "isatty not implemented for -1", `//os.isatty(-1)`)
}

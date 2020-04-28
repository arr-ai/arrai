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
	AssertCodesEvalToSameValue(t, `{|@,@byte| (0,97), (1,98), (2,99)}`, `//os.stdin`)
	AssertCodesEvalToSameValue(t, `{|@,@byte| (0,97), (1,98), (2,99)}`, `//os.stdin`)
}

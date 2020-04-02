package syntax

import (
	"bytes"
	"testing"
)

func TestStdOsStdin(t *testing.T) {
	// Not parallelisable
	stdin := stdOsStdinHandle
	defer func() { stdOsStdinHandle = stdin }()

	stdOsStdinHandle = bytes.NewBuffer([]byte(""))
	AssertCodesEvalToSameValue(t, `{}`, `//.os.stdin`)

	stdOsStdinHandle = bytes.NewBuffer([]byte("abc"))
	AssertCodesEvalToSameValue(t, `{|@,@byte| (0,97), (1,98), (2,99)}`, `//.os.stdin`)
}

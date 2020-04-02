package syntax

import (
	"bytes"
	"testing"
)

func TestStdOsStdin(t *testing.T) {
	stdin := stdOsStdinHandle
	stdOsStdinHandle = bytes.NewBuffer([]byte("abc"))
	AssertCodesEvalToSameValue(t, `{|@,@byte| (0,97), (1,98), (2,99)}`, `//.os.stdin`)
	stdOsStdinHandle = stdin
}

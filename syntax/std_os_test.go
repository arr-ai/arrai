package syntax

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestStdOsGetArgs(t *testing.T) {
	os.Args = []string{"arrai", "-d", "r", "file.arrai", "arg1", "arg2", "arg3"}
	stdOsGetArgs()
	assert.Equal(t, stdOsGetArgs(), strArrToRelArr(os.Args[3:]))

	os.Args = []string{"arrai", "r", "file.arrai", "arg1", "arg2", "arg3"}
	stdOsGetArgs()
	assert.Equal(t, stdOsGetArgs(), strArrToRelArr(os.Args[2:]))
}

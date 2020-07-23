package rel

import (
	"fmt"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
)

func TestGetLastScope(t *testing.T) {
	t.Parallel()

	scope := EmptyScope.With("stuff", NewNumber(1)).With("random", NewNumber(2))
	err := WrapContext(
		WrapContext(
			WrapContext(fmt.Errorf("random"), EmptyTuple, scope),
			EmptyTuple,
			scope,
		),
		EmptyTuple,
		scope,
	)
	assert.True(t, err.(ContextErr).GetLastScope().m.Equal(scope.m))

	err = WrapContext(fmt.Errorf("random"), EmptyTuple, scope)
	assert.True(t, err.(ContextErr).GetLastScope().m.Equal(scope.m))

	err = WrapContext(fmt.Errorf("random"), EmptyTuple, EmptyScope)
	assert.True(t, err.(ContextErr).GetLastScope().m.Equal(EmptyScope.m))
}

func TestGetImportantFrames(t *testing.T) {
	t.Parallel()

	err := WrapContext(fmt.Errorf("random"), EmptyTuple, EmptyScope)
	assert.Equal(t, []ContextErr{err.(ContextErr)}, err.(ContextErr).GetImportantFrames())

	err2 := WrapContext(err, EmptyTuple, EmptyScope)
	err3 := WrapContext(err2, EmptyTuple, EmptyScope)

	assert.Equal(t, []ContextErr{err.(ContextErr)}, err3.(ContextErr).GetImportantFrames())

	err = NewContextErr(fmt.Errorf("random"), *parser.NewScannerAt("random", 5, 1), EmptyScope)
	err2 = NewContextErr(err, *parser.NewScannerAt("random", 3, 1), EmptyScope)
	err3 = NewContextErr(err2, *parser.NewScannerAt("random", 1, 1), EmptyScope)
	assert.Equal(t,
		[]ContextErr{
			err3.(ContextErr),
			err2.(ContextErr),
			err.(ContextErr),
		},
		err3.(ContextErr).GetImportantFrames(),
	)
}

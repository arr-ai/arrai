package rel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLastScope(t *testing.T) {
	t.Parallel()

	scope := EmptyScope.With("stuff", NewNumber(1)).With("random", NewNumber(2))
	err := wrapContext(
		wrapContext(
			wrapContext(fmt.Errorf("random"), EmptyTuple, scope),
			EmptyTuple,
			scope,
		),
		EmptyTuple,
		scope,
	)
	assert.True(t, err.(ContextErr).GetLastScope().m.Equal(scope.m))

	err = wrapContext(fmt.Errorf("random"), EmptyTuple, scope)
	assert.True(t, err.(ContextErr).GetLastScope().m.Equal(scope.m))

	err = wrapContext(fmt.Errorf("random"), EmptyTuple, EmptyScope)
	assert.True(t, err.(ContextErr).GetLastScope().m.Equal(EmptyScope.m))
}

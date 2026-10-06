package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentExprEvalSelfEvaluatingValuesDoNotAllocate(t *testing.T) {
	bound := map[string]Value{
		"str":   NewString([]rune("properties")),
		"num":   NewNumber(1234.5),
		"arr":   NewArray(NewNumber(1), NewNumber(2)),
		"bytes": NewBytes([]byte("xyz")),
	}
	ctx := context.Background()
	for name, v := range bound {
		name, v := name, v
		t.Run(name, func(t *testing.T) {
			scope := EmptyScope.With(name, v)
			ident := NewIdentExpr(*parser.NewScanner(name), name)

			got, err := ident.Eval(ctx, scope)
			require.NoError(t, err)
			assert.True(t, v.Equal(got))

			allocs := testing.AllocsPerRun(100, func() {
				_, _ = ident.Eval(ctx, scope)
			})
			assert.Zero(t, allocs)
		})
	}
}

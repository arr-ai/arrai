package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanLiftLiteralKeepsOneBox(t *testing.T) {
	for name, v := range map[string]Value{
		"str":   NewString([]rune("properties")),
		"num":   NewNumber(1234.5),
		"array": NewArray(NewNumber(1), NewNumber(2)),
		"bytes": NewBytes([]byte("xyz")),
	} {
		v := v
		t.Run(name, func(t *testing.T) {
			node, err := encodeExpr(NewLiteralExpr(planSrc, v))
			require.NoError(t, err)
			lifted, err := decodeExpr(node)
			require.NoError(t, err)
			require.IsType(t, LiteralExpr{}, lifted)

			ctx := context.Background()
			got, err := lifted.Eval(ctx, EmptyScope)
			require.NoError(t, err)
			assert.True(t, v.Equal(got))

			allocs := testing.AllocsPerRun(100, func() {
				if _, err := lifted.Eval(ctx, EmptyScope); err != nil {
					t.Fatal(err)
				}
			})
			assert.Zero(t, allocs)
		})
	}
}

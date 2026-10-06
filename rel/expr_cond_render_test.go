package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stringCountingExpr is a condition that is always false and counts how often
// it is rendered to a string.
type stringCountingExpr struct {
	ExprScanner
	renders *int
}

func (e stringCountingExpr) Eval(context.Context, Scope) (Value, error) { return False, nil }

func (e stringCountingExpr) String() string {
	*e.renders++
	return "cond"
}

func TestCondExprEvalDoesNotRenderConditions(t *testing.T) {
	var renders int
	scanner := *parser.NewScanner("")
	dict, err := NewDictExpr(scanner, false, true,
		NewDictEntryTupleExpr(scanner, stringCountingExpr{renders: &renders}, NewNumber(1)),
		NewDictEntryTupleExpr(scanner, NewIdentExpr(scanner, "_"), NewNumber(2)),
	)
	require.NoError(t, err)

	got, err := NewCondExpr(scanner, dict).Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.True(t, NewNumber(2).Equal(got))
	assert.Zero(t, renders)
}

func TestCondExprDiscardMayBeParenthesised(t *testing.T) {
	scanner := *parser.NewScanner("")
	underscore := NewIdentExpr(scanner, "_")
	for name, at := range map[string]Expr{
		"bare":          underscore,
		"parenthesised": NewExprExpr(scanner, underscore),
		"nested":        NewExprExpr(scanner, NewExprExpr(scanner, underscore)),
	} {
		at := at
		t.Run(name, func(t *testing.T) {
			dict, err := NewDictExpr(scanner, false, true,
				NewDictEntryTupleExpr(scanner, NewLiteralExpr(scanner, False), NewNumber(1)),
				NewDictEntryTupleExpr(scanner, at, NewNumber(2)),
			)
			require.NoError(t, err)

			got, err := NewCondExpr(scanner, dict).Eval(context.Background(), EmptyScope)
			require.NoError(t, err)
			assert.True(t, NewNumber(2).Equal(got))
		})
	}
}

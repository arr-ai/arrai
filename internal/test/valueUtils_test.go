package test

import (
	"context"
	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsLiteralTrue(t *testing.T) {
	t.Parallel()

	require.True(t, isLiteralTrue(eval("true")))
	require.True(t, isLiteralTrue(eval("{()}")))

	require.False(t, isLiteralTrue(eval("false")))
	require.False(t, isLiteralTrue(eval("1")))
	require.False(t, isLiteralTrue(eval("0")))
	require.False(t, isLiteralTrue(eval("()")))
	require.False(t, isLiteralTrue(eval("{true}")))
	require.False(t, isLiteralTrue(eval("{false}")))
	require.False(t, isLiteralTrue(eval("{(1)}")))
	require.False(t, isLiteralTrue(eval("{(),1}")))
	require.False(t, isLiteralTrue(eval("[true]")))
	require.False(t, isLiteralTrue(eval("(val: true)")))
}

func TestIsLiteralFalse(t *testing.T) {
	t.Parallel()

	require.True(t, isLiteralFalse(eval("false")))
	require.True(t, isLiteralFalse(eval("{}")))

	require.False(t, isLiteralFalse(eval("true")))
	require.False(t, isLiteralFalse(eval("1")))
	require.False(t, isLiteralFalse(eval("0")))
	require.False(t, isLiteralFalse(eval("()")))
	require.False(t, isLiteralFalse(eval("{true}")))
	require.False(t, isLiteralFalse(eval("{false}")))
	require.False(t, isLiteralFalse(eval("{(1)}")))
	require.False(t, isLiteralFalse(eval("{(),1}")))
	require.False(t, isLiteralFalse(eval("[true]")))
	require.False(t, isLiteralFalse(eval("(val: true)")))
}

func TestForeachLeaf(t *testing.T) {
	t.Parallel()

	require.Equal(t,
		shouldBe{"<root>": "true"},
		forInput("true"))
	require.Equal(t,
		shouldBe{"<root>": "false"},
		forInput("false"))
	require.Equal(t,
		shouldBe{"<root>": "42"},
		forInput("42"))

	require.Equal(t,
		shouldBe{},
		forInput("()"))
	require.Equal(t,
		shouldBe{"a": "true", "b": "false"},
		forInput("(a: true, b: false)"))
	require.Equal(t,
		shouldBe{"a.b.c": "true", "d.e.f": "false"},
		forInput("(a: (b: (c: true)), d: (e: (f: false)))"))
	require.Equal(t,
		shouldBe{"a(0)": "0", "a(1)": "1", "a(2)": "2"},
		forInput("(a: [0, 1, 2])"))
	require.Equal(t,
		shouldBe{"a(0)(0)": "0", "a(0)(1)": "1", "a(1)(0)": "10", "a(1)(1)": "11"},
		forInput("(a: [[0, 1], [10, 11]])"))
}

type shouldBe map[string]string

func forInput(source string) shouldBe {
	leaves := make(shouldBe)
	tree := eval(source)
	if tree == nil {
		return nil
	}

	ForeachLeaf(tree, "<root>", func(val rel.Value, path string) {
		valStr := val.String()
		if valStr == "{}" {
			valStr = "false"
		}
		leaves[path] = valStr
	})

	return leaves
}

func eval(source string) rel.Value {
	ctx := arraictx.InitRunCtx(context.Background())
	value, err := syntax.EvaluateExpr(ctx, "", source)
	if err != nil {
		return nil
	}
	return value
}

package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/require"
)

var noSrc = *parser.NewScanner("")

func ident(name string) Expr { return NewIdentExpr(noSrc, name) }

func num(n float64) Expr { return NewNumber(n) }

// let builds the tree `let name = rhs; body` exactly as the compiler does.
func let(name string, rhs, body Expr) Expr {
	return NewArrowExpr(noSrc, rhs, NewFunction(noSrc, IdentPattern(name), body))
}

func lambda(name string, body Expr) Expr { return NewFunction(noSrc, IdentPattern(name), body) }

func isLet(e Expr) bool {
	a, is := e.(*ArrowExpr)
	if !is {
		return false
	}
	_, is = a.fn.arg.(IdentPattern)
	return is
}

func evalTo(t *testing.T, e Expr) Value {
	t.Helper()
	v, err := e.Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	return v
}

func TestSimplifyFoldsOneShotLet(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	// let x = 1; x + 2  ==>  1 + 2
	e := Simplify(let("x", num(1), NewAddExpr(noSrc, ident("x"), num(2))))
	require.False(t, isLet(e), "%v", e)
	require.Zero(t, countIdentUses(e, "x"))
	require.True(t, evalTo(t, e).Equal(NewNumber(3)))
}

func TestSimplifyFoldsChainedLets(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	// let a = 1; let b = a; b + 1  ==>  1 + 1
	e := Simplify(let("a", num(1), let("b", ident("a"), NewAddExpr(noSrc, ident("b"), num(1)))))
	require.False(t, isLet(e), "%v", e)
	require.True(t, evalTo(t, e).Equal(NewNumber(2)))
}

func TestSimplifyKeepsMultiUseLet(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	e := Simplify(let("x", num(1), NewAddExpr(noSrc, ident("x"), ident("x"))))
	require.True(t, isLet(e), "%v", e)
}

func TestSimplifyKeepsDynamicUse(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	s := MustNewSet(NewNumber(1), NewNumber(2))
	cases := map[string]Expr{
		// let x = 3; s => x + .
		"mapped": let("x", num(3), NewDArrowExpr(noSrc, s, lambda(".", NewAddExpr(noSrc, ident("x"), ident("."))))),
		"where": let("x", num(3), NewWhereExpr(noSrc, s, lambda(".", NewCompareExpr(noSrc,
			[]Expr{ident("."), ident("x")},
			[]CompareFunc{func(a, b Value) (bool, error) { return a.Less(b), nil }},
			[]string{"<"})))),
		"lambda": let("x", num(3), lambda("y", NewAddExpr(noSrc, ident("x"), ident("y")))),
		"ifTrue": let("x", num(3), NewIfElseExpr(noSrc, ident("x"), num(1), num(0))),
		"andRhs": let("x", num(3), NewAndExpr(noSrc, num(1), ident("x"))),
		"seqBody": let("x", num(3), NewSeqArrowExpr(false)(noSrc, NewArray(NewNumber(1)),
			lambda(".", NewAddExpr(noSrc, ident("x"), ident("."))))),
	}
	for name, e := range cases {
		require.True(t, isLet(Simplify(e)), "%s: %v", name, e)
	}
}

func TestSimplifyRefusesCapture(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	// let y = z; let z = 1; y + z + z — folding y into the inner body would
	// rebind z. The inner let keeps its two uses, so its binder stays, and
	// y's single once-position use sits under it.
	inner := let("z", num(1), NewAddExpr(noSrc, NewAddExpr(noSrc, ident("y"), ident("z")), ident("z")))
	e := Simplify(let("y", ident("z"), inner))
	require.True(t, isLet(e), "%v", e)
	require.Equal(t, 1, countIdentUses(e.(*ArrowExpr).fn.body, "y"), "%v", e)
}

func TestSimplifyRefusesEffectfulRhs(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	// A call may do anything, so `let x = f(1); x + 1` keeps its order.
	e := Simplify(let("x", NewCallExpr(noSrc, ident("f"), num(1)), NewAddExpr(noSrc, ident("x"), num(1))))
	require.True(t, isLet(e), "%v", e)
}

func TestSimplifyKeepsRecursiveLet(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	f := NewRecursionExpr(noSrc, "f", lambda("n", ident("n")))
	e := Simplify(let("f", f, NewCallExpr(noSrc, ident("f"), num(1))))
	require.True(t, isLet(e), "%v", e)
}

func TestSimplifyDropsUnusedTotalLet(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	body := NewAddExpr(noSrc, num(1), num(1))
	for name, rhs := range map[string]Expr{
		"literal": num(1),
		"lambda":  lambda("y", ident("y")),
		"tuple":   NewTupleExpr(noSrc, AttrExpr{name: "a", expr: NewAddExpr(noSrc, ident("q"), num(1))}),
	} {
		e := let("x", rhs, body)
		if name == "tuple" {
			// q is unbound, so the tuple is not total: kept.
			require.True(t, isLet(Simplify(e)), "%s", name)
			continue
		}
		require.False(t, isLet(Simplify(e)), "%s: %v", name, e)
	}
}

func TestSimplifyKeepsUnusedFailingLet(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	// let x = (1).a; 2 — the rhs fails, and eager evaluation reports it.
	e := Simplify(let("x", NewDotExpr(noSrc, num(1), "a"), num(2)))
	require.True(t, isLet(e), "%v", e)
	_, err := e.Eval(context.Background(), EmptyScope)
	require.Error(t, err)
}

func TestSimplifyLeavesOpaqueNodes(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("simplifier is off in the slowpath build")
	}

	// An Expr type the walker does not know hides its uses.
	e := Simplify(let("x", num(1), opaqueExpr{ident("x")}))
	require.True(t, isLet(e), "%v", e)
}

type opaqueExpr struct{ inner Expr }

func (o opaqueExpr) Eval(ctx context.Context, s Scope) (Value, error) { return o.inner.Eval(ctx, s) }
func (o opaqueExpr) Source() parser.Scanner                           { return noSrc }
func (o opaqueExpr) String() string                                   { return "opaque" }

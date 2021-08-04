package syntax

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

// assertCompileScanner checks that the expr resulting from an arr.ai expression
// shrinkwraps it.
func assertCompileScanner(t *testing.T, source string) bool { //nolint:unparam
	pc := ParseContext{SourceDir: ".."}
	// Add some space padding, which should not become part of the source.
	ctx := arraictx.InitRunCtx(context.Background())
	ast, err := pc.Parse(ctx, parser.NewScanner(" "+source+" "))
	require.NoError(t, err, "%s", source)
	expr, err := pc.CompileExpr(ctx, ast)
	require.NoError(t, err)
	return assert.Equal(t, source, expr.Source().String())
}

func TestCompileScannerAtom(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `x`)
	assertCompileScanner(t, `1`)
	assertCompileScanner(t, `'abc'`)
	assertCompileScanner(t, `<<'abc', 100>>`)
	assertCompileScanner(t, `%x`)
}

func TestCompileScannerCompositeLiteral(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `[]`)
	assertCompileScanner(t, `[1, 2, 3]`)
	assertCompileScanner(t, `(a: 1, b: 2)`)
	assertCompileScanner(t, `{'a', 'b', 'c'}`)
}

func TestCompileScannerParenExpr(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `(a + b)`)
}

func TestCompileScannerCompositeExpr(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `[a, b, c]`)
	assertCompileScanner(t, `(:a, :b, :c)`)
	assertCompileScanner(t, `{a, b, c}`)
	assertCompileScanner(t, `{a:b, b:c, c:d}`)
}

func TestCompileScannerUnExpr(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `-b`)
	assertCompileScanner(t, `+b`)
	assertCompileScanner(t, `>> (x / 3)`)
}

func TestCompileScannerBinExpr(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `a + b`)
	assertCompileScanner(t, `() ++ b`)
	assertCompileScanner(t, `() ++ {}`)
	assertCompileScanner(t, `() ++ 1`)
	assertCompileScanner(t, `1 + {}`)

	assertCompileScanner(t, `a >> b`)
	assertCompileScanner(t, `[] >> b`)
	assertCompileScanner(t, `[] >> {}`)
	assertCompileScanner(t, `[] >> 1`)
	assertCompileScanner(t, `[] >> 1`)

	assertCompileScanner(t, `[1, 2, 3] where .@ < 2`)
}

func TestCompileScannerCompare(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `a < b`)
	assertCompileScanner(t, `a <= b < (c + 2)`)
}

func TestCompileScannerDot(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `a.b`)
	assertCompileScanner(t, `().b`)
	assertCompileScanner(t, `().b.c`)
}

func TestCompileScannerCond(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `cond {a: b}`)
	assertCompileScanner(t, `cond () {b: c}`)
}

func TestCompileScannerLet(t *testing.T) {
	t.Parallel()

	assertCompileScanner(t, `let x = 1; x + (y / (z - w))`)
}

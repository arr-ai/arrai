//nolint:unparam
package shell

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineCollectorAppendLine(t *testing.T) {
	t.Parallel()

	c := newLineCollector()

	// testing string with no opener
	c.appendLine("random string")
	assert.Equal(t, []string{"random string"}, c.lines)
	assert.True(t, c.isBalanced())
	c.reset()

	// basic testing string with opener
	c.appendLine("random {")
	assert.Equal(t, []string{"random {"}, c.lines)
	assert.Equal(t, 1, len(c.stack))
	assert.Equal(t, "}", c.stack[0].char)
	assert.Equal(t, true, c.stack[0].recursive)
	assert.False(t, c.isBalanced())
	c.appendLine("}")
	assert.Equal(t, []string{"random {", "}"}, c.lines)
	assert.Equal(t, 0, len(c.stack))
	assert.True(t, c.isBalanced())
	c.reset()

	// testing with multiple openers
	c.appendLine("{([")
	assert.Equal(t, 3, len(c.stack))
	assert.Equal(t, "}", c.stack[0].char)
	assert.Equal(t, ")", c.stack[1].char)
	assert.Equal(t, "]", c.stack[2].char)
	c.appendLine("])}")
	assert.Equal(t, 0, len(c.stack))
	assert.True(t, c.isBalanced())
	c.reset()

	// testing context based opener
	c.appendLine("nested closer $'{([")
	c.appendLine("${")
	assert.Equal(t, []string{"nested closer $'{([", "${"}, c.lines)
	assert.Equal(t, 2, len(c.stack))
	assert.Equal(t, "'", c.stack[0].char)
	assert.Equal(t, true, c.stack[0].recursive)
	assert.Equal(t, "}", c.stack[1].char)
	assert.Equal(t, true, c.stack[1].recursive)
	c.appendLine("}")
	assert.Equal(t, 1, len(c.stack))
	assert.Equal(t, "'", c.stack[0].char)
	assert.Equal(t, true, c.stack[0].recursive)
	c.appendLine("'")
	assert.Equal(t, 0, len(c.stack))
	assert.True(t, c.isBalanced())
	c.reset()

	// testing non recursive opener
	c.appendLine("'")
	c.appendLine("\"")
	assert.Equal(t, []string{"'", "\""}, c.lines)
	assert.Equal(t, 1, len(c.stack))
	c.appendLine("'")
	assert.Equal(t, 0, len(c.stack))
	assert.True(t, c.isBalanced())
	c.reset()

	// testing `` escape
	c.appendLine("`stuff``")
	assert.Equal(t, []string{"`stuff``"}, c.lines)
	assert.Equal(t, 1, len(c.stack))
	assert.False(t, c.isBalanced())
	c.appendLine("`")
	assert.Equal(t, 0, len(c.stack))
	assert.True(t, c.isBalanced())
	c.reset()

	// testing escape
	c.appendLine("'\\\"'")
	c.appendLine("'\\\\'")
	c.appendLine("'\\n'")
	c.appendLine("`\\`")
	assert.True(t, c.isBalanced())
}

func TestIsBalanced(t *testing.T) {
	t.Parallel()

	c := newLineCollector()
	assert.True(t, c.isBalanced())

	c.appendLine("random;")
	assert.False(t, c.isBalanced())

	c.appendLine("\\a \\b \\c")
	assert.False(t, c.isBalanced())

	c.appendLine("\\.")
	assert.False(t, c.isBalanced())

	c.appendLine("\\a random")
	assert.True(t, c.isBalanced())

	c.appendLine("c:")
	assert.False(t, c.isBalanced())

	c.appendLine("\"")
	assert.False(t, c.isBalanced())

	c.appendLine("\"")
	assert.True(t, c.isBalanced())
	c.reset()

	c.appendLine("\"\\\"\"")
	assert.True(t, c.isBalanced())
	c.reset()

	c.appendLine("\"\\\"x\"")
	assert.True(t, c.isBalanced())
	c.reset()

	c.appendLine("\"\\\"xx\"")
	assert.True(t, c.isBalanced())
	c.reset()

	c.appendLine("\"\\\"xxx\"")
	assert.True(t, c.isBalanced())
	c.reset()

	c.appendLine(`let f = \x \y x + y; f(3, 4)`)
	assert.True(t, c.isBalanced())
	c.reset()

	c.appendLine(`let f = \x \y`)
	assert.False(t, c.isBalanced())
}

func TestGetLastToken(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "//str", getLastToken([]rune("//str")))
	assert.Equal(t, "//", getLastToken([]rune("//")))
	assert.Equal(t, "///", getLastToken([]rune("///")))
	assert.Equal(t, "//", getLastToken([]rune("//seq.contains(//")))
	assert.Equal(t, "//arch", getLastToken([]rune("//seq.contains(//arch")))
	assert.Equal(t, "tuple.", getLastToken([]rune("//seq.contains(tuple.")))
	assert.Equal(t, "x.", getLastToken([]rune("x.")))
	assert.Equal(t, "x", getLastToken([]rune("x")))
	assert.Equal(t, "", getLastToken([]rune("//seq.contains(")))
	assert.Equal(t, "", getLastToken([]rune("")))
}

func TestTabCompletionStdlib(t *testing.T) {
	t.Parallel()
	stdlib := syntax.StdScope().MustGet(".").(rel.Tuple)
	stdlibNames := stdlib.Names().OrderedNames()

	assertTabCompletion(t, append(stdlibNames, "{"), 0, "//\t", nil)
	assertTabCompletion(t, append(stdlibNames, "{"), 0, "//seq.contains(//\t", nil)
	prefix := "s"

	assertTabCompletionWithPrefix(t, prefix, stdlibNames, "//%s\t", nil)
	assertTabCompletionWithPrefix(t, prefix, stdlibNames, "x(//%s\t", nil)
	assertTabCompletionWithPrefix(t, prefix, stdlibNames, "x(//%s\t + random)", nil)

	lib := "seq"
	strlib := stdlib.MustGet(lib).(rel.Tuple).Names().OrderedNames()
	assertTabCompletionWithPrefix(t, prefix, strlib, "//"+lib+".%s\t", nil)
	for i := 0; i < len(strlib); i++ {
		strlib[i] = "." + strlib[i]
	}
	assertTabCompletionWithPrefix(t, "", strlib, "//"+lib+"%s\t", nil)
}

func TestTrimExpr(t *testing.T) {
	t.Parallel()

	sh := newShellInstance(newLineCollector(), syntax.StdScope())

	realExpr, residue := sh.trimExpr(`x.`)
	assert.Equal(t, "x", realExpr)
	assert.Equal(t, ".", residue)

	realExpr, residue = sh.trimExpr(`abc(`)
	assert.Equal(t, "abc", realExpr)
	assert.Equal(t, "(", residue)

	realExpr, residue = sh.trimExpr(`x`)
	assert.Equal(t, "x", realExpr)
	assert.Equal(t, "", residue)

	realExpr, residue = sh.trimExpr(`x.'`)
	assert.Equal(t, "x", realExpr)
	assert.Equal(t, ".'", residue)

	realExpr, residue = sh.trimExpr(`abc("`)
	assert.Equal(t, "abc", realExpr)
	assert.Equal(t, "(\"", residue)

	realExpr, residue = sh.trimExpr("x(`")
	assert.Equal(t, "x", realExpr)
	assert.Equal(t, "(`", residue)

	realExpr, residue = sh.trimExpr(`x.'random.random`)
	assert.Equal(t, "x", realExpr)
	assert.Equal(t, ".'random.random", residue)

	realExpr, residue = sh.trimExpr(`abc("abc(`)
	assert.Equal(t, "abc", realExpr)
	assert.Equal(t, "(\"abc(", residue)

	realExpr, residue = sh.trimExpr("x(abc.(bca`")
	assert.Equal(t, "x", realExpr)
	assert.Equal(t, "(abc.(bca`", residue)

	realExpr, residue = sh.trimExpr(`abc.ab`)
	assert.Equal(t, "abc", realExpr)
	assert.Equal(t, ".ab", residue)

	//FIXME: unable to handle this case
	// realExpr, residue = sh.trimExpr("x(abc).\"(bca`")
	// assert.Equal(t, "x(abc)", realExpr)
	// assert.Equal(t, ".\"(bca`", residue)
}

func TestCompletionCurrentExpr(t *testing.T) {
	t.Parallel()

	assertTabCompletion(t, []string{".a"}, 0, "(a: 1)\t", nil)
	assertTabCompletion(t, []string{".a", ".b"}, 0, "(a: 1, b: 2)\t", nil)
	assertTabCompletion(t, []string{".a", ".b"}, 0, "(a: 1, b: 2)\t + 123", nil)
	assertTabCompletion(t, []string{"a"}, 1, "(a: 1).\t", nil)
	assertTabCompletion(t, []string{"a", "b"}, 1, "(a: 1, b: 2).\t", nil)
	assertTabCompletion(t, []string{"a", "b"}, 1, "(a: 1, b: 2).\t + 123", nil)
	assertTabCompletion(t, []string{".c"}, 0, "(a: (c: 3), b: 2).a\t", nil)
	assertTabCompletion(t, []string{`'random string'`}, 1, "(`random string`: 1).\t", nil)
	assertTabCompletion(t, []string{".a"}, 0, "x\t", map[string]string{"x": "(a: 1)"})
	assertTabCompletion(t,
		[]string{"a", `'b"b'`, `"c'c"`, "'d`d'", "\"e\\\"e'e`ee\""}, 1,
		"(a: 1, 'b\"b': 2, \"c'c\": 3, \"d`d\": 4, \"e\\\"e'e`ee\": 5).\t", nil)
	assertTabCompletion(t,
		[]string{"", "a", "aa"}, 2,
		"x.a\t", map[string]string{"x": "(a: 1, aa: 2, aaa: 3)"})
	assertTabCompletion(t,
		[]string{"", "a", "aa", ".a"}, 2,
		"x.a\t", map[string]string{"x": "(a: (a: 1), aa: 2, aaa: 3)"})
	assertTabCompletion(t,
		[]string{"", "a", "aa", ".a", ".b"}, 2,
		"x.a\t", map[string]string{"x": "(a: (a: 1, b: 2), aa: 2, aaa: 3)"})

	assertTabCompletion(t, []string{`('a')`}, 0, "{`a`: 1}\t", nil)
	assertTabCompletion(t, []string{`('a')`, `('b')`}, 0, "{`a`: 1, `b`: 2}\t", nil)
	assertTabCompletion(t, []string{`('a')`, `('b')`}, 0, "{`a`: 1, `b`: 2}\t + 123", nil)
	assertTabCompletion(t, []string{`'a')`}, 1, "{`a`: 1}(\t", nil)
	assertTabCompletion(t, []string{`'a')`, `'b')`}, 1, "{`a`: 1, `b`: 2}(\t", nil)
	assertTabCompletion(t, []string{`'a')`, `'b')`}, 1, "{`a`: 1, `b`: 2}(\t + 123", nil)
	assertTabCompletion(t, []string{`'c')`}, 1, "{`a`: {`c`: 3}, `b`: 2}(`a`)(\t", nil)
	assertTabCompletion(t, []string{`'random string')`}, 1, "{`random string`: 1}(\t", nil)
	assertTabCompletion(t, []string{`('a')`}, 0, "x\t", map[string]string{"x": "{`a`: 1}"})
	assertTabCompletion(t,
		[]string{"('a')", `('b"b')`, `("c'c")`, "('d`d')", "(\"e\\\"e'e`ee\")"}, 0,
		"{'a': 1, 'b\"b': 2, \"c'c\": 3, \"d`d\": 4, \"e\\\"e'e`ee\": 5}\t", nil)
	assertTabCompletion(t,
		[]string{`(2)`, `('string')`, `([1, 2, 3])`}, 0,
		"{`string`: 1, 2: 20, [1, 2, 3]: 30}\t", nil)
	assertTabCompletion(t,
		[]string{`bc')`, `bcd')`, `bd')`}, 3,
		"{`abc`: 1, `abcd`: 2, `abd`: 3, `bd`: 4}('a\t", nil)
	assertTabCompletion(t,
		[]string{`('abc')`}, 0,
		"{`abc`: {`abc`: 1}, `abcd`: 2, `abd`: 3, `bd`: 4}('abc')\t", nil)

	assertTabCompletion(t, []string{`.a`}, 0, "let x = (a: 1); x\t", nil)
	assertTabCompletion(t, []string{`('a')`}, 0, "let x = {`a`: 1}; x\t", nil)
	assertTabCompletion(t, []string{`.a`}, 0, "x\t", map[string]string{"x": `(a: {"b": (c: 3)})`})
	assertTabCompletion(t, []string{`('b')`}, 0, "x.a\t", map[string]string{"x": `(a: {"b": (c: 3)})`})
	assertTabCompletion(t, []string{`.c`}, 0, "x.a(`b`)\t", map[string]string{"x": `(a: {"b": (c: 3)})`})
}

func assertTabCompletionWithPrefix(
	t *testing.T,
	prefix string,
	choices []string,
	format string,
	scopeValues map[string]string,
) {
	var libWithPrefix []string
	for _, c := range choices {
		if strings.HasPrefix(c, prefix) {
			libWithPrefix = append(libWithPrefix, strings.TrimPrefix(c, prefix))
		}
	}
	assertTabCompletion(t, libWithPrefix, len(prefix), fmt.Sprintf(format, prefix), scopeValues)
}

func assertTabCompletion(t *testing.T,
	expectedPredictions []string,
	expectedLength int,
	line string,
	scopeValues map[string]string,
) {
	scope := syntax.StdScope()
	for name, expr := range scopeValues {
		val, err := syntax.EvaluateExpr("", expr)
		require.NoError(t, err)
		scope = scope.With(name, val)
	}
	sh := newShellInstance(newLineCollector(), scope)
	predictions, length := sh.Do([]rune(line), strings.Index(line, "\t"))
	strPredictions := make([]string, 0, len(predictions))
	for _, p := range predictions {
		strPredictions = append(strPredictions, string(p))
	}
	assert.Equal(t, expectedPredictions, strPredictions)
	assert.Equal(t, expectedLength, length)
}

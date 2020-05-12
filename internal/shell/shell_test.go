//nolint:unparam
package shell

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "//", getLastToken([]rune("//str.contains(//")))
	assert.Equal(t, "//arch", getLastToken([]rune("//str.contains(//arch")))
	assert.Equal(t, "tuple.", getLastToken([]rune("//str.contains(tuple.")))
	assert.Equal(t, "", getLastToken([]rune("//str.contains(")))
	assert.Equal(t, "", getLastToken([]rune("")))
}

func TestTabCompletionStdlib(t *testing.T) {
	t.Parallel()
	stdlib := syntax.StdScope().MustGet(".").(rel.Tuple)
	stdlibNames := stdlib.Names().OrderedNames()

	assertTabCompletion(t, append(stdlibNames, "{"), 0, "//\t", nil)
	assertTabCompletion(t, append(stdlibNames, "{"), 0, "//str.contains(//\t", nil)
	prefix := "s"

	assertTabCompletionWithPrefix(t, prefix, stdlibNames, "//%s\t", nil)
	assertTabCompletionWithPrefix(t, prefix, stdlibNames, "x(//%s\t", nil)
	assertTabCompletionWithPrefix(t, prefix, stdlibNames, "x(//%s\t + random)", nil)

	lib := "str"
	strlib := stdlib.MustGet(lib).(rel.Tuple).Names().OrderedNames()
	assertTabCompletionWithPrefix(t, prefix, strlib, "//"+lib+".%s\t", nil)
}

func assertTabCompletionWithPrefix(
	t *testing.T,
	prefix string,
	choices []string,
	format string,
	scopeValues map[string]rel.Expr,
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
	scopeValues map[string]rel.Expr,
) {
	scope := syntax.StdScope()
	for name, expr := range scopeValues {
		scope = scope.With(name, expr)
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

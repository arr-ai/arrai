package syntax

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFmtPrettyForDict(t *testing.T) {
	t.Parallel()
	simpleDict, err := EvaluateExpr(".", `{'a':1, 'b':2, 'c':3}`)
	assert.Nil(t, err)
	str := FormatString(simpleDict, 0)
	assert.Equal(t, true, strings.Contains(str, "\n\ta: 1"))
	assert.Equal(t, true, strings.Contains(str, "\n\tb: 2"))
	assert.Equal(t, true, strings.Contains(str, "\n\tc: 3"))

	complexDict, err := EvaluateExpr(".", `{'a':1, 'b':2, 'c':(d:11, e:22, f:{111, 222})}`)
	assert.Nil(t, err)
	str = FormatString(complexDict, 0)
	assert.Equal(t, true, strings.Contains(str, "\n\ta: 1"))
	assert.Equal(t, true, strings.Contains(str, "\n\tb: 2"))
	assert.Equal(t, true, strings.Contains(str, "\n\tc: ("))
	assert.Equal(t, true, strings.Contains(str, "\n\t\td: 11"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\te: 22"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\tf: {"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\t\t111"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\t\t222"))
}

func TestFmtPrettyForSet(t *testing.T) {
	t.Parallel()
	simpleSet, err := EvaluateExpr(".", `{24, 25, 26}`)
	assert.Nil(t, err)
	str := FormatString(simpleSet, 0)
	assert.Equal(t, true, strings.Contains(str, "\n\t24"))
	assert.Equal(t, true, strings.Contains(str, "\n\t25"))
	assert.Equal(t, true, strings.Contains(str, "\n\t26"))

	complexSet, err := EvaluateExpr(".", `{(a: 1),(b: 2),(c:{11,22,33})}`)
	str = FormatString(complexSet, 0)
	assert.Nil(t, err)
	assert.Equal(t, true, strings.Contains(str, "\n\t\ta: 1"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\tb: 2"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\t\t11"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\t\t11"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\t\t22"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\t\t33"))
}

func TestFmtPrettyForTuple(t *testing.T) {
	t.Parallel()
	simpleSet, err := EvaluateExpr(".", `(a:1, b:2, c:3)`)
	assert.Nil(t, err)
	str := FormatString(simpleSet, 0)
	assert.Equal(t, true, strings.Contains(str, "\n\ta: 1"))
	assert.Equal(t, true, strings.Contains(str, "\n\tb: 2"))
	assert.Equal(t, true, strings.Contains(str, "\n\tc: 3"))

	complexTuple, err := EvaluateExpr(".", `(a:1, b:(d:11, e:12, f:{1, 2}), c:3)`)
	assert.Nil(t, err)
	str = FormatString(complexTuple, 0)
	assert.Equal(t, true, strings.Contains(str, "\n\ta: 1"))
	assert.Equal(t, true, strings.Contains(str, "\n\tc: 3"))
	assert.Equal(t, true, strings.Contains(str, "\n\tb: ("))
	assert.Equal(t, true, strings.Contains(str, "\n\t\td: 11"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\te: 12"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\tf: {"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\t\t1"))
	assert.Equal(t, true, strings.Contains(str, "\n\t\t\t2"))
}

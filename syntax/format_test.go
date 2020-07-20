package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFmtPrettyForDict(t *testing.T) { //nolint:dupl
	t.Parallel()
	simpleDict, err := EvaluateExpr(".", `{'c':3, 'a':1, 'b':2}`)
	assert.Nil(t, err)
	str, err := FormatString(simpleDict, 0)
	assert.Nil(t, err)
	assert.Equal(t, "{\n  a: 1,\n  b: 2,\n  c: 3\n}", str)

	complexDict, err := EvaluateExpr(".", `{'a':1, 'c':(d:11, e:22, f:{111, 222}), 'b':2}`)
	assert.Nil(t, err)
	str, err = FormatString(complexDict, 0)
	assert.Nil(t, err)
	assert.Equal(t,
		"{\n  a: 1,\n  b: 2,\n  c: (\n    d: 11,\n    e: 22,\n    f: {\n      111,\n      222\n    }\n  )\n}",
		str)
}

func TestFmtPrettyForSet(t *testing.T) { //nolint:dupl
	t.Parallel()
	simpleSet, err := EvaluateExpr(".", `{26, 24, 25}`)
	assert.Nil(t, err)
	str, err := FormatString(simpleSet, 0)
	assert.Nil(t, err)
	assert.Equal(t, "{\n  24,\n  25,\n  26\n}", str)

	complexSet, err := EvaluateExpr(".", `{(a: 1),(b: 2),(c:{11,22,33})}`)
	assert.Nil(t, err)
	str, err = FormatString(complexSet, 0)
	assert.Nil(t, err)
	assert.Equal(t,
		"{\n  (\n    a: 1\n  ),\n  (\n    b: 2\n  ),\n  (\n    c: {\n      11,\n      22,\n      33\n    }\n  )\n}",
		str)
}

func TestFmtPrettyForTuple(t *testing.T) { //nolint:dupl
	t.Parallel()
	simpleSet, err := EvaluateExpr(".", `(a:1, c:3, b:2)`)
	assert.Nil(t, err)
	str, err := FormatString(simpleSet, 0)
	assert.Nil(t, err)
	assert.Equal(t, "(\n  a: 1,\n  b: 2,\n  c: 3\n)", str)

	complexTuple, err := EvaluateExpr(".", `(a:1, b:(d:11, e:12, f:{1, 2}), c:3)`)
	assert.Nil(t, err)
	str, err = FormatString(complexTuple, 0)
	assert.Nil(t, err)
	assert.Equal(t,
		"(\n  a: 1,\n  b: (\n    d: 11,\n    e: 12,\n    f: {\n      1,\n      2\n    }\n  ),\n  c: 3\n)",
		str)
}

func TestFmtPrettyForArray(t *testing.T) { //nolint:dupl
	t.Parallel()
	array, err := EvaluateExpr(".", `[1, 2, 3, 5, 6, 4, 10]`)
	assert.Nil(t, err)
	str, err := FormatString(array, 0)
	assert.Nil(t, err)
	assert.Equal(t, "[1, 2, 3, 5, 6, 4, 10]", str)
}

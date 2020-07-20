package syntax

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFmtPrettyForDict(t *testing.T) { //nolint:dupl
	t.Parallel()
	simpleDict, err := EvaluateExpr(".", `{'a':1, 'b':2, 'c':3}`)
	assert.Nil(t, err)
	str := FormatString(simpleDict, 0)
	assert.Equal(t, true, contains(str, "\n%sa: 1", 1))
	assert.Equal(t, true, contains(str, "\n%sb: 2", 1))
	assert.Equal(t, true, contains(str, "\n%sc: 3", 1))

	complexDict, err := EvaluateExpr(".", `{'a':1, 'b':2, 'c':(d:11, e:22, f:{111, 222})}`)
	assert.Nil(t, err)
	str = FormatString(complexDict, 0)
	assert.Equal(t, true, contains(str, "\n%sa: 1", 1))
	assert.Equal(t, true, contains(str, "\n%sb: 2", 1))
	assert.Equal(t, true, contains(str, "\n%sc: (", 1))
	assert.Equal(t, true, contains(str, "\n%sd: 11", 2))
	assert.Equal(t, true, contains(str, "\n%se: 22", 2))
	assert.Equal(t, true, contains(str, "\n%sf: {", 2))
	assert.Equal(t, true, contains(str, "\n%s111", 3))
	assert.Equal(t, true, contains(str, "\n%s222", 3))
}

func TestFmtPrettyForSet(t *testing.T) { //nolint:dupl
	t.Parallel()
	simpleSet, err := EvaluateExpr(".", `{24, 25, 26}`)
	assert.Nil(t, err)
	str := FormatString(simpleSet, 0)
	assert.Equal(t, true, contains(str, "\n%s24", 1))
	assert.Equal(t, true, contains(str, "\n%s25", 1))
	assert.Equal(t, true, contains(str, "\n%s26", 1))

	complexSet, err := EvaluateExpr(".", `{(a: 1),(b: 2),(c:{11,22,33})}`)
	str = FormatString(complexSet, 0)
	assert.Nil(t, err)
	assert.Equal(t, true, contains(str, "\n%sa: 1", 2))
	assert.Equal(t, true, contains(str, "\n%sb: 2", 2))
	assert.Equal(t, true, contains(str, "\n%s11", 3))
	assert.Equal(t, true, contains(str, "\n%s11", 3))
	assert.Equal(t, true, contains(str, "\n%s22", 3))
	assert.Equal(t, true, contains(str, "\n%s33", 3))
}

func TestFmtPrettyForTuple(t *testing.T) { //nolint:dupl
	t.Parallel()
	simpleSet, err := EvaluateExpr(".", `(a:1, b:2, c:3)`)
	assert.Nil(t, err)
	str := FormatString(simpleSet, 0)
	assert.Equal(t, true, contains(str, "\n%sa: 1", 1))
	assert.Equal(t, true, contains(str, "\n%sb: 2", 1))
	assert.Equal(t, true, contains(str, "\n%sc: 3", 1))

	complexTuple, err := EvaluateExpr(".", `(a:1, b:(d:11, e:12, f:{1, 2}), c:3)`)
	assert.Nil(t, err)
	str = FormatString(complexTuple, 0)
	assert.Equal(t, true, contains(str, "\n%sa: 1", 1))
	assert.Equal(t, true, contains(str, "\n%sc: 3", 1))
	assert.Equal(t, true, contains(str, "\n%sb: (", 1))
	assert.Equal(t, true, contains(str, "\n%sd: 11", 2))
	assert.Equal(t, true, contains(str, "\n%se: 12", 2))
	assert.Equal(t, true, contains(str, "\n%sf: {", 2))
	assert.Equal(t, true, contains(str, "\n%s1", 3))
	assert.Equal(t, true, contains(str, "\n%s2", 3))
}

func contains(s, substrf string, identsNum int) bool {
	return strings.Contains(s, fmt.Sprintf(substrf, getIndents(identsNum)))
}

package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTupleMatcher(t *testing.T) {
	t.Parallel()

	var a int
	matcher := NewTupleMatcher(
		map[string]Matcher{
			"a": MatchInt(func(i int) { a = i }),
		},
		Lit(EmptyTuple),
	)

	assert.Equal(t, true, matcher.Match(NewTuple(NewAttr("a", NewNumber(1)))))
	assert.Equal(t, 1, a)
}

func TestTupleMatcherMissingAttr(t *testing.T) {
	t.Parallel()

	a := -1
	matcher := NewTupleMatcher(
		map[string]Matcher{
			"a": MatchInt(func(i int) { a = i }),
		},
		Lit(EmptyTuple),
	)

	assert.Equal(t, false, matcher.Match(NewTuple()))
	assert.Equal(t, -1, a)
}

func TestTupleMatcherRest(t *testing.T) {
	t.Parallel()

	var a int
	matcher := NewTupleMatcher(
		map[string]Matcher{
			"a": MatchInt(func(i int) { a = i }),
		},
		Lit(NewTuple(NewAttr("b", NewNumber(2)))),
	)
	v := NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(2)))

	assert.Equal(t, true, matcher.Match(v))
	assert.Equal(t, 1, a)
}

func TestTupleMatcherMissingRest(t *testing.T) {
	t.Parallel()

	var a int
	matcher := NewTupleMatcher(
		map[string]Matcher{
			"a": MatchInt(func(i int) { a = i }),
		},
		Lit(NewTuple(NewAttr("b", NewNumber(2)))),
	)
	v := NewTuple(NewAttr("a", NewNumber(1)))

	assert.Equal(t, false, matcher.Match(v))
	assert.Equal(t, 1, a)
}

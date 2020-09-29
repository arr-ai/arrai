package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindMatchNumber(t *testing.T) {
	t.Parallel()

	val := NewNumber(1)

	m := Bind(&val)

	assert.Equal(t, true, m.Match(NewNumber(2)))
}

func TestBindMatchString(t *testing.T) {
	t.Parallel()

	val := NewString([]rune("foo"))
	target := val.(Set)

	m := Bind(&val)

	assert.Equal(t, true, m.Match(NewString([]rune("foo"))))
	assert.Equal(t, &target, m.target)
}

func TestIntSetMatcherAllInts(t *testing.T) {
	t.Parallel()

	values := []int{}
	matcher := NewSetMatcher(MatchInt(func(i int) { values = append(values, i) }))
	set := MustNewSet(NewNumber(1), NewNumber(2), NewNumber(3))

	assert.Equal(t, true, matcher.Match(set))
	assert.ElementsMatch(t, []int{1, 2, 3}, values)
}

func TestIntSetMatcherMixedTypes(t *testing.T) {
	t.Parallel()

	values := []int{}
	matcher := NewSetMatcher(MatchInt(func(i int) { values = append(values, i) }))
	set := MustNewSet(NewNumber(1), NewNumber(2), None)

	assert.Equal(t, false, matcher.Match(set))
	assert.Equal(t, []int{}, values)
}

func TestSetMatcherNonSet(t *testing.T) {
	t.Parallel()

	values := []int{}
	matcher := NewSetMatcher(MatchInt(func(i int) { values = append(values, i) }))

	assert.Equal(t, false, matcher.Match(NewTuple()))
	assert.Equal(t, []int{}, values)
}

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

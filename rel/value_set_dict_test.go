package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/frozen"
)

func TestDictEntryTupleLess(t *testing.T) {
	t.Parallel()
	// Test for a panic when comparing Strings using !=.
	assert.NotPanics(t, func() {
		MustNewDict(false,
			NewDictEntryTuple(NewString([]rune("b")), NewNumber(2)),
			NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)),
		).(Dict).OrderedEntries()
	})
}

func TestDictEntryTupleOrdered(t *testing.T) {
	t.Parallel()

	entries := MustNewDict(true,
		NewDictEntryTuple(NewString([]rune("b")), NewNumber(2)),
		NewDictEntryTuple(NewString([]rune("a")), NewNumber(2)),
		NewDictEntryTuple(NewString([]rune("b")), NewNumber(1)),
		NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)),
	).(Dict).OrderedEntries()

	AssertEqualValues(t, NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)), entries[0])
	AssertEqualValues(t, NewDictEntryTuple(NewString([]rune("a")), NewNumber(2)), entries[1])
	AssertEqualValues(t, NewDictEntryTuple(NewString([]rune("b")), NewNumber(1)), entries[2])
	AssertEqualValues(t, NewDictEntryTuple(NewString([]rune("b")), NewNumber(2)), entries[3])
}

// TestDictHash128AgreesWithGenericSetEqual guards against a hash/equals
// contract violation: Dict.Equal treats a Dict as equal to a generic Set of
// the same (@, @value) entry tuples, and an empty Dict as equal to
// EmptySet{}. Hash128 must agree in both cases, or a Set/Map keyed on Dict
// values can silently fail to recognize an existing equal entry depending
// on which representation it happens to be constructed as
// (see https://github.com/arr-ai/arrai/issues/PLACEHOLDER).
func TestDictHash128AgreesWithGenericSetEqual(t *testing.T) {
	t.Parallel()

	key := NewString([]rune("k"))
	value := NewNumber(1)
	d := MustNewDict(false, NewDictEntryTuple(key, value)).(Dict)
	entryTuple := NewTuple(NewAttr("@", key), NewAttr(DictValueAttr, value))
	gs := GenericSet{set: frozen.NewSet[Value](entryTuple)}

	assert.True(t, d.Equal(gs), "a Dict must equal an equivalent generic Set of entry tuples")
	assert.Equal(t, d.Hash128(), gs.Hash128(), "equal Dict/Set values must hash the same")

	assert.True(t, Dict{}.Equal(EmptySet{}), "an empty Dict must equal EmptySet{}")
	assert.Equal(t, Dict{}.Hash128(), EmptySet{}.Hash128(), "an empty Dict must hash like EmptySet{}")
}

func TestDictLess(t *testing.T) {
	t.Parallel()

	kv := func(k, v float64) DictEntryTuple {
		return NewDictEntryTuple(NewNumber(k), NewNumber(v))
	}
	assertLess := func(a, b Set) {
		assert.True(t, a.Less(b))
		assert.False(t, b.Less(a))
	}
	assertLess(MustNewDict(true, kv(1, 42)), MustNewDict(true, kv(1, 43)))
	assertLess(MustNewDict(true, kv(1, 42)), MustNewDict(true, kv(1, 43), kv(2, 44)))
	assertLess(MustNewDict(true, kv(1, 42)), MustNewDict(true, kv(1, 42), kv(1, 44)))
	assertLess(MustNewDict(true, kv(1, 41), kv(1, 42)), MustNewDict(true, kv(1, 42)))
	assertLess(MustNewDict(true, kv(1, 42)), MustNewDict(true, kv(1, 43), kv(2, 42)))
	assertLess(MustNewDict(true, kv(1, 42), kv(2, 43)), MustNewDict(true, kv(1, 42), kv(3, 43)))

	assertSame := func(a, b Set) {
		assert.False(t, a.Less(b))
		assert.False(t, b.Less(a))
	}
	assertSame(MustNewDict(true, kv(1, 43), kv(2, 42)), MustNewDict(true, kv(1, 43), kv(2, 42)))
}

func TestDictCallAll(t *testing.T) {
	t.Parallel()

	kv := func(k, v float64) DictEntryTuple {
		return NewDictEntryTuple(NewNumber(k), NewNumber(v))
	}
	dict := MustNewDict(false, kv(1, 10), kv(2, 20), kv(3, 30))
	ctx := arraictx.InitRunCtx(context.Background())
	AssertEqualValues(t, MustNewSet(NewNumber(10)), mustCallAll(ctx, dict, NewNumber(1)))
	AssertEqualValues(t, MustNewSet(NewNumber(20)), mustCallAll(ctx, dict, NewNumber(2)))
	AssertEqualValues(t, MustNewSet(NewNumber(30)), mustCallAll(ctx, dict, NewNumber(3)))
	AssertEqualValues(t, None, mustCallAll(ctx, dict, NewNumber(4)))

	dict = MustNewDict(true, kv(1, 10), kv(1, 11), kv(2, 20), kv(3, 30))

	AssertEqualValues(t, MustNewSet(NewNumber(10), NewNumber(11)), mustCallAll(ctx, dict, NewNumber(1)))
	AssertEqualValues(t, MustNewSet(NewNumber(20)), mustCallAll(ctx, dict, NewNumber(2)))
	AssertEqualValues(t, MustNewSet(NewNumber(30)), mustCallAll(ctx, dict, NewNumber(3)))
	AssertEqualValues(t, None, mustCallAll(ctx, dict, NewNumber(4)))
}

func TestDictWithMultipleEntriesOfSameValue(t *testing.T) {
	t.Parallel()

	d1 := MustNewDict(true,
		NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)),
		NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)),
	)
	assert.Equal(t, 1, d1.Count())
	for e := d1.(Dict).DictEnumerator(); e.MoveNext(); {
		k, v := e.Current()
		AssertEqualValues(t, NewString([]rune("a")), k)
		AssertEqualValues(t, NewNumber(1), v)
	}

	d2 := MustNewDict(true,
		NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)),
	)
	d2 = d2.With(NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)))
	assert.Equal(t, 1, d2.Count())
	for e := d2.(Dict).DictEnumerator(); e.MoveNext(); {
		k, v := e.Current()
		AssertEqualValues(t, NewString([]rune("a")), k)
		AssertEqualValues(t, NewNumber(1), v)
	}
}

package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMode(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		OnlyOnLHS,
		createMode([]int{2, 3}, []int{2, 3}, []int{0, 1}, []int{}),
	)
	assert.Equal(t,
		OnlyOnLHS|InBoth,
		createMode([]int{2, 3}, []int{2, 3}, []int{0, 1, 2, 3}, []int{}),
	)
	assert.Equal(t,
		OnlyOnRHS,
		createMode([]int{2, 3}, []int{2, 3}, []int{}, []int{0, 1}),
	)
	assert.Equal(t,
		OnlyOnRHS|InBoth,
		createMode([]int{2, 3}, []int{2, 3}, []int{}, []int{0, 1, 2, 3}),
	)
	assert.Equal(t,
		OnlyOnLHS|OnlyOnRHS,
		createMode([]int{2, 3}, []int{2, 3}, []int{0, 1}, []int{0, 1}),
	)
	assert.Equal(t,
		InBoth,
		createMode([]int{2, 3}, []int{2, 3}, []int{}, []int{2, 3}),
	)
	assert.Equal(t,
		InBoth,
		createMode([]int{2, 3}, []int{2, 3}, []int{2, 3}, []int{}),
	)
	assert.Equal(t,
		OnlyOnLHS|InBoth|OnlyOnRHS,
		createMode([]int{2, 3}, []int{2, 3}, []int{0, 1, 2, 3}, []int{0, 1}),
	)
	assert.Equal(t,
		OnlyOnLHS|InBoth|OnlyOnRHS,
		createMode([]int{2, 3}, []int{2, 3}, []int{0, 1}, []int{0, 1, 2, 3}),
	)
	assert.Panics(t, func() {
		createMode([]int{2, 3, 4}, []int{2, 3}, []int{}, []int{})
	})
	assert.Panics(t, func() {
		createMode([]int{2, 3, 4}, []int{2, 3}, []int{0, 1}, []int{0, 1, 2, 3})
	})
}

func row(numbers ...int) Values {
	v := make(Values, 0, len(numbers))
	for _, n := range numbers {
		v = append(v, NewNumber(float64(n)))
	}
	return v
}

func TestGroupBy(t *testing.T) {
	t.Parallel()

	row1 := row(1, 1, 2)
	row2 := row(1, 1, 3)
	row3 := row(1, 2, 3)
	pr := newPositionalRelation(3, row1, row2, row3)

	testGroup := func(grouper valueProjector, groups int) {
		index := pr.groupBy(grouper)
		assert.Equal(t, groups, index.count(), "group count for %v", grouper)
		// Probe with each original row: the index must return a group
		// containing that row.
		for e := pr.Range(); e.Next(); {
			r := e.Values()
			b := index.bucketOfRow(r, grouper)
			if !assert.NotNil(t, b, "row %v has no group for %v", r, grouper) {
				continue
			}
			found := false
			for _, id := range b.rows {
				if index.row(id).equalValues(r) {
					found = true
					break
				}
			}
			assert.True(t, found, "row %v missing from its own group", r)
		}
		// The memoised index is the same one.
		assert.Same(t, index, pr.groupBy(grouper))
	}

	testGroup(valueProjector{}, 1)
	testGroup(valueProjector{0}, 1)
	testGroup(valueProjector{1}, 2)
	testGroup(valueProjector{0, 1}, 2)
	testGroup(valueProjector{2, 0}, 2)
	testGroup(valueProjector{1, 2}, 3)
	testGroup(valueProjector{0, 1, 2}, 3)
}

func TestDerivedViewInheritsGroupIndex(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("slowpath rebuilds indexes per view")
	}
	row1 := row(1, 1, 2)
	row2 := row(1, 1, 3)
	row3 := row(1, 2, 3)
	row4 := row(2, 2, 4)
	pr := newPositionalRelation(3, row1, row2, row3, row4)
	p := valueProjector{0}
	parent := pr.groupBy(p)
	require.Equal(t, 2, parent.count())
	require.False(t, parent.filtered)

	view, err := pr.Where(func(v Values) (bool, error) {
		return v[0].(Number).Float64() == 1, nil
	})
	require.NoError(t, err)
	require.Equal(t, 3, view.Count())
	require.Same(t, pr, view.parent)

	got := view.groupBy(p)
	assert.True(t, got.filtered, "derived view rebuilt the parent index")
	assert.Equal(t, 1, got.count(), "filtered index should keep only key=1")
	b := got.bucketOfRow(row1, p)
	require.NotNil(t, b)
	assert.Equal(t, 3, len(b.rows), "three rows have col0=1")
}

func TestProjectDotsSharesStoreWhenInjective(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("slowpath copies projections")
	}
	s := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(10))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewNumber(20))),
	)
	r := s.(Relation)
	got, ok := r.projectDots([]string{"a"}, []string{"a"})
	require.True(t, ok)
	p := got.(Relation)
	assert.Same(t, r.rows, p.rows, "injective project must not copy the store")
	assert.Equal(t, r.Count(), p.Count())
	assert.True(t, p.Has(NewTuple(NewAttr("a", NewNumber(1)))))
	assert.False(t, p.Has(NewTuple(NewAttr("a", NewNumber(9)))))
}

func TestProjectDotsCopiesWhenNotInjective(t *testing.T) {
	t.Parallel()
	s := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(10))),
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(20))),
	)
	r := s.(Relation)
	got, ok := r.projectDots([]string{"a"}, []string{"a"})
	require.True(t, ok)
	p := got.(Relation)
	assert.NotSame(t, r.rows, p.rows, "non-injective project must copy to dedup")
	assert.Equal(t, 1, p.Count())
}

func TestSequenceAtKeyIsSeeded(t *testing.T) {
	t.Parallel()
	pr := newPositionalRelation(2, row(0, 1), row(1, 2))
	r := newRelation(NamesSlice{"@", ArrayItemAttr}, valueProjector{0, 1}, pr)
	assert.True(t, r.rows.hasCandidateKey(valueProjector{0}), "@ of a sequence shape is a key")
}

func TestEmpiricalCandidateKeyIsCached(t *testing.T) {
	t.Parallel()
	pr := newPositionalRelation(2, row(1, 10), row(2, 20), row(3, 10))
	uniq := valueProjector{0}
	dup := valueProjector{1}
	require.Equal(t, 3, pr.groupBy(uniq).count())
	require.Equal(t, 2, pr.groupBy(dup).count())
	assert.True(t, pr.hasCandidateKey(uniq), "col0 is unique, should be cached as a key")
	assert.False(t, pr.hasCandidateKey(dup), "col1 has duplicates")
}

// The With fast path appends to a store only when the view is the store's
// frontier; older views and siblings must copy out, never see the new row,
// and never go quadratic on a simple chain.
func TestPositionalRelationWithChain(t *testing.T) {
	t.Parallel()

	r := newPositionalRelation(2)
	generations := []*positionalRelation{r}
	for i := 0; i < 100; i++ {
		r = r.With(row(i, i*2))
		generations = append(generations, r)
	}
	for i, g := range generations {
		assert.Equal(t, i, g.Count(), "generation %d", i)
		if i > 0 {
			assert.True(t, g.Has(row(i-1, (i-1)*2)))
		}
		assert.False(t, g.Has(row(100, 200)))
	}
	// A duplicate add returns the same view.
	assert.Same(t, r, r.With(row(7, 14)))

	// Branch off an old generation: the sibling copies out and neither
	// branch sees the other's rows.
	branch := generations[50].With(row(1000, 2000))
	assert.Equal(t, 51, branch.Count())
	assert.True(t, branch.Has(row(1000, 2000)))
	assert.False(t, r.Has(row(1000, 2000)))
	assert.True(t, r.Has(row(99, 198)))
	assert.False(t, branch.Has(row(99, 198)))
}

func TestPositionalRelationWhereSharesStore(t *testing.T) {
	t.Parallel()

	r := newPositionalRelation(2, row(1, 10), row(2, 20), row(3, 30), row(4, 40))
	even, err := r.Where(func(v Values) (bool, error) {
		return int(v[0].(Number).Float64())%2 == 0, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, even.Count())
	assert.Same(t, r.store, even.store, "where must share the arena")
	assert.True(t, even.Has(row(2, 20)))
	assert.False(t, even.Has(row(1, 10)))

	// Adding a row that is in the arena but excluded by the selection must
	// re-add it, not find the hidden copy.
	readded := even.With(row(1, 10))
	assert.Equal(t, 3, readded.Count())
	assert.True(t, readded.Has(row(1, 10)))
	assert.Equal(t, 2, even.Count(), "the source view must not change")

	// A filtered view forces With to copy out.
	grown := even.Without(row(2, 20)).With(row(5, 50))
	assert.Equal(t, 2, grown.Count())
	assert.True(t, grown.Has(row(4, 40)))
	assert.True(t, grown.Has(row(5, 50)))
	assert.False(t, grown.Has(row(2, 20)))

	// Views of one store and of different stores compare by content.
	all, err := r.Where(func(Values) (bool, error) { return true, nil })
	assert.NoError(t, err)
	assert.Same(t, r, all)
	rebuilt := newPositionalRelation(2, row(4, 40), row(3, 30), row(2, 20), row(1, 10))
	assert.True(t, r.EqualPositionalRelation(rebuilt))
	assert.Equal(t, r.Hash128(), rebuilt.Hash128())
	assert.False(t, r.EqualPositionalRelation(even))
}

func TestPositionalRelationProject(t *testing.T) {
	t.Parallel()

	r := newPositionalRelation(3, row(1, 10, 7), row(2, 20, 7), row(3, 20, 7))
	// Identity: same view.
	assert.Same(t, r, r.Project(valueProjector{0, 1, 2}))
	// Permutation: distinct rows stay distinct.
	perm := r.Project(valueProjector{2, 0, 1})
	assert.Equal(t, 3, perm.Count())
	assert.True(t, perm.Has(row(7, 1, 10)))
	// Narrowing deduplicates.
	narrowed := r.Project(valueProjector{1, 2})
	assert.Equal(t, 2, narrowed.Count())
	assert.True(t, narrowed.Has(row(10, 7)))
	assert.True(t, narrowed.Has(row(20, 7)))
	constant := r.Project(valueProjector{2})
	assert.Equal(t, 1, constant.Count())
}

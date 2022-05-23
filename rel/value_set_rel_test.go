package rel

import (
	"testing"

	"github.com/arr-ai/frozen"
	"github.com/stretchr/testify/assert"
)

func TestProjectionBasedOnNames(t *testing.T) {
	t.Parallel()

	r := Relation{attrs: NamesSlice{"a", "b", "c"}, p: valueProjector{0, 1, 2}}
	assert.Equal(t, valueProjector{2, 1, 0}, r.projectionBasedOnNames(NamesSlice{"c", "b", "a"}))
	assert.Equal(t, valueProjector{0, 1, 2}, r.projectionBasedOnNames(NamesSlice{"a", "b", "c"}))
	assert.Equal(t, valueProjector{0, 1}, r.projectionBasedOnNames(NamesSlice{"a", "b"}))
	assert.Equal(t, valueProjector{}, r.projectionBasedOnNames(NamesSlice{}))
	assert.Equal(t, valueProjector{0, 0, 2}, r.projectionBasedOnNames(NamesSlice{"a", "a", "c"}))

	r = Relation{attrs: NamesSlice{"a", "b", "c"}, p: valueProjector{2, 0, 1}}
	assert.Equal(t, valueProjector{1, 0, 2}, r.projectionBasedOnNames(NamesSlice{"c", "b", "a"}))
	assert.Equal(t, valueProjector{2, 0, 1}, r.projectionBasedOnNames(NamesSlice{"a", "b", "c"}))
	assert.Equal(t, valueProjector{2, 0}, r.projectionBasedOnNames(NamesSlice{"a", "b"}))
	assert.Equal(t, valueProjector{}, r.projectionBasedOnNames(NamesSlice{}))
	assert.Equal(t, valueProjector{2, 2, 1}, r.projectionBasedOnNames(NamesSlice{"a", "a", "c"}))

	assert.Panics(t, func() { r.projectionBasedOnNames(NamesSlice{"d"}) })
}

func TestRelationString(t *testing.T) {
	t.Parallel()

	r := Relation{
		attrs: NamesSlice{"c", "b", "a"},
		p:     valueProjector{0, 1, 2},
		rows: &positionalRelation{
			set: frozen.NewSet(
				row(1, 1, 2),
				row(1, 2, 3),
				row(2, 1, 2),
			),
		},
	}
	assert.Equal(t, "{|a, b, c| (2, 1, 1), (2, 1, 2), (3, 2, 1)}", r.String())

	r = Relation{
		attrs: NamesSlice{"c", "b", "a"},
		p:     valueProjector{2, 0, 1},
		rows: &positionalRelation{
			set: frozen.NewSet(
				row(1, 1, 2),
				row(1, 2, 3),
				row(2, 1, 2),
			),
		},
	}
	assert.Equal(t, "{|a, b, c| (1, 1, 2), (1, 2, 2), (2, 1, 3)}", r.String())
}

func TestRelationUnion(t *testing.T) {
	t.Parallel()

	r1 := Relation{
		attrs: NamesSlice{"a", "b"},
		p:     valueProjector{0, 1},
		rows: &positionalRelation{
			set: frozen.NewSet(row(1, 3)),
		},
	}

	r2 := Relation{
		attrs: NamesSlice{"b", "a"},
		p:     valueProjector{0, 1},
		rows: &positionalRelation{
			set: frozen.NewSet(row(1, 3)),
		},
	}

	// this ensures that even if NamesSlice is in different order, as long both Relations have the same names, the union
	// of the Relations should be the same type of Relation.
	AssertEqualValues(t,
		Relation{
			attrs: NamesSlice{"a", "b"},
			p:     valueProjector{0, 1},
			rows: &positionalRelation{
				set: frozen.NewSet(row(1, 3), row(3, 1)),
			},
		},
		Union(r1, r2),
	)
}

func TestRelationHas(t *testing.T) {
	t.Parallel()

	r := Relation{
		attrs: NamesSlice{"c", "b", "a"},
		p:     valueProjector{0, 1, 2},
		rows: &positionalRelation{
			set: frozen.NewSet(row(1, 3, 2)),
		},
	}
	assert.True(t,
		r.Has(
			NewTuple(
				NewAttr("a", NewNumber(2)),
				NewAttr("b", NewNumber(3)),
				NewAttr("c", NewNumber(1)),
			),
		),
	)
	assert.False(t,
		r.Has(
			NewTuple(
				NewAttr("a", NewNumber(0)),
				NewAttr("b", NewNumber(0)),
				NewAttr("c", NewNumber(0)),
			),
		),
	)
	assert.False(t,
		r.Has(
			NewTuple(
				NewAttr("a", NewNumber(0)),
				NewAttr("b", NewNumber(0)),
			),
		),
	)
}

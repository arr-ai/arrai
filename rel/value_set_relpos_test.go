package rel

import (
	"testing"

	"github.com/arr-ai/frozen"
	"github.com/stretchr/testify/assert"
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

	prb := &positionalRelationBuilder{&frozen.SetBuilder[any]{}}
	prb.Add(row1)
	prb.Add(row2)
	prb.Add(row3)
	pr := prb.Finish()

	kv := func(k any, v frozen.Set[any]) frozen.KeyValue[any, frozen.Set[any]] {
		return frozen.KV[any, frozen.Set[any]](k, v)
	}
	s := func(rows ...any) frozen.Set[any] {
		return frozen.NewSet[any](rows...)
	}
	testGroup := func(grouper valueProjector, grouped frozen.Map[any, frozen.Set[any]]) {
		assert.True(t, pr.groupBy(grouper).Equal(grouped))
		assert.True(t, pr.meta.indices.MustGet(grouper).(frozen.Map[any, frozen.Set[any]]).Equal(grouped))
	}

	testGroup(valueProjector{}, frozen.NewMap(kv(row(), s(row1, row2, row3))))

	testGroup(valueProjector{0}, frozen.NewMap(kv(row(1), s(row1, row2, row3))))

	testGroup(
		valueProjector{1},
		frozen.NewMap(
			kv(row(1), s(row1, row2)),
			kv(Values{NewNumber(2)}, s(row3)),
		),
	)

	testGroup(
		valueProjector{0, 1},
		frozen.NewMap(
			kv(row(1, 1), s(row1, row2)),
			kv(row(1, 2), s(row3)),
		),
	)

	testGroup(
		valueProjector{2, 0},
		frozen.NewMap(
			kv(row(2, 1), s(row1)),
			kv(row(3, 1), s(row2, row3)),
		),
	)

	testGroup(
		valueProjector{1, 2},
		frozen.NewMap(
			kv(row(1, 2), s(row1)),
			kv(row(1, 3), s(row2)),
			kv(row(2, 3), s(row3)),
		),
	)

	testGroup(
		valueProjector{0, 1, 2},
		frozen.NewMap(
			kv(row1, s(row1)),
			kv(row2, s(row2)),
			kv(row3, s(row3)),
		),
	)
}

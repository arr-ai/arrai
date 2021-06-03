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

func TestGroupBy(t *testing.T) {
	t.Parallel()
	row := func(numbers ...int) Values {
		v := make(Values, 0, len(numbers))
		for _, n := range numbers {
			v = append(v, NewNumber(float64(n)))
		}
		return v
	}

	row1 := row(1, 1, 2)
	row2 := row(1, 1, 3)
	row3 := row(1, 2, 3)

	prb := &positionalRelationBuilder{&frozen.SetBuilder{}}
	prb.Add(row1)
	prb.Add(row2)
	prb.Add(row3)
	pr := prb.Finish()

	testGroup := func(grouper valueProjector, grouped frozen.Map) {
		assert.True(t, pr.groupBy(grouper).Equal(grouped))
		assert.True(t, pr.meta.indices.MustGet(grouper).(frozen.Map).Equal(grouped))
	}

	testGroup(valueProjector{}, frozen.NewMap(frozen.KV(row(), frozen.NewSet(row1, row2, row3))))

	testGroup(valueProjector{0}, frozen.NewMap(frozen.KV(row(1), frozen.NewSet(row1, row2, row3))))

	testGroup(
		valueProjector{1},
		frozen.NewMap(
			frozen.KV(row(1), frozen.NewSet(row1, row2)),
			frozen.KV(Values{NewNumber(2)}, frozen.NewSet(row3)),
		),
	)

	testGroup(
		valueProjector{0, 1},
		frozen.NewMap(
			frozen.KV(row(1, 1), frozen.NewSet(row1, row2)),
			frozen.KV(row(1, 2), frozen.NewSet(row3)),
		),
	)

	testGroup(
		valueProjector{2, 0},
		frozen.NewMap(
			frozen.KV(row(2, 1), frozen.NewSet(row1)),
			frozen.KV(row(3, 1), frozen.NewSet(row2, row3)),
		),
	)

	testGroup(
		valueProjector{1, 2},
		frozen.NewMap(
			frozen.KV(row(1, 2), frozen.NewSet(row1)),
			frozen.KV(row(1, 3), frozen.NewSet(row2)),
			frozen.KV(row(2, 3), frozen.NewSet(row3)),
		),
	)

	testGroup(
		valueProjector{0, 1, 2},
		frozen.NewMap(
			frozen.KV(row1, frozen.NewSet(row1)),
			frozen.KV(row2, frozen.NewSet(row2)),
			frozen.KV(row3, frozen.NewSet(row3)),
		),
	)
}

package rel

import (
	"testing"

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

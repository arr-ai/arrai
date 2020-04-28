package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIndexes(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []int{1, 2, 3, 4, 5}, getIndexes(1, 6, 1, false))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, getIndexes(1, 5, 1, true))
	assert.Equal(t, []int{2, 4}, getIndexes(2, 5, 2, false))
	assert.Equal(t, []int{1, 3, 5}, getIndexes(1, 5, 2, true))
	assert.Equal(t, []int{5, 3}, getIndexes(5, 1, -2, false))
	assert.Equal(t, []int{5, 3, 1}, getIndexes(5, 1, -2, true))
	assert.Equal(t, []int{}, getIndexes(5, 1, 1, true))
	assert.Equal(t, []int{}, getIndexes(5, 1, 0, false))
}

// func TestInitDefaultArrayIndex(t *testing.T) {
// 	t.Parallel()

// 	start, end := initDefaultArrayIndex(Number(10), Number(20), 5, 25, 1)
// 	assert.Equal(t, 10, start)
// 	assert.Equal(t, 20, end)

// 	start, end = initDefaultArrayIndex(nil, nil, 5, 25, 1)
// 	assert.Equal(t, 5, start)
// 	assert.Equal(t, 25, end)

// 	start, end = initDefaultArrayIndex(nil, nil, 5, 25, -2)
// 	assert.Equal(t, 24, start)
// 	assert.Equal(t, 4, end)

// 	start, end = initDefaultArrayIndex(nil, Number(12), 5, 25, -2)
// 	assert.Equal(t, 24, start)
// 	assert.Equal(t, 12, end)

// 	start, end = initDefaultArrayIndex(nil, Number(12), 5, 25, 2)
// 	assert.Equal(t, 5, start)
// 	assert.Equal(t, 12, end)

// 	start, end = initDefaultArrayIndex(Number(7), nil, 5, 25, -2)
// 	assert.Equal(t, 7, start)
// 	assert.Equal(t, 4, end)

// 	start, end = initDefaultArrayIndex(Number(7), nil, 5, 25, 2)
// 	assert.Equal(t, 7, start)
// 	assert.Equal(t, 25, end)

// 	start, end = initDefaultArrayIndex(nil, Number(42), 5, 25, -2)
// 	assert.Equal(t, 24, start)
// 	assert.Equal(t, 25, end)

// 	start, end = initDefaultArrayIndex(Number(42), nil, 5, 25, 2)
// 	assert.Equal(t, 24, start)
// 	assert.Equal(t, 25, end)

// 	start, end = initDefaultArrayIndex(Number(-5), nil, 5, 25, 2)
// 	assert.Equal(t, 5, start)
// 	assert.Equal(t, 25, end)

// 	start, end = initDefaultArrayIndex(nil, Number(-5), 5, 25, 2)
// 	assert.Equal(t, 5, start)
// 	assert.Equal(t, 5, end)

// 	start, end = initDefaultArrayIndex(Number(-30), nil, 5, 25, 2)
// 	assert.Equal(t, 5, start)
// 	assert.Equal(t, 25, end)

// 	start, end = initDefaultArrayIndex(nil, Number(-30), 5, 25, 2)
// 	assert.Equal(t, 5, start)
// 	assert.Equal(t, 5, end)

// 	start, end = initDefaultArrayIndex(Number(42), Number(-30), 5, 25, 2)
// 	assert.Equal(t, 24, start)
// 	assert.Equal(t, 5, end)

// 	start, end = initDefaultArrayIndex(Number(-30), Number(42), 5, 25, 2)
// 	assert.Equal(t, 5, start)
// 	assert.Equal(t, 25, end)
// }

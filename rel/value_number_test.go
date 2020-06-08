package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func sum(nums []float64) float64 {
	result := 0.0
	for _, n := range nums {
		result += n
	}
	return result
}

func TestNumberULP(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "0.3", NewNumber(sum([]float64{0.1, 0.1, 0.1})).String())
}

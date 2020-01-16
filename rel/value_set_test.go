package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayCall(t *testing.T) {
	f := NewArray(
		NewNumber(0),
		NewNumber(1),
		NewNumber(4),
		NewNumber(9),
		NewNumber(16),
		NewNumber(25),
	)
	for i := 0; i < f.Count(); i++ {
		assert.Equal(t, i*i, int(f.Call(NewNumber(float64(i))).(Number).Float64()))
	}
}

package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStringTuple(t *testing.T) {
	s := "hello"
	for e := NewString([]rune(s)).Enumerator(); e.MoveNext(); {
		match := stringTupleMatcher()
		index, char, is := match(e.Current())
		if assert.True(t, is) {
			assert.Equal(t, rune(s[index]), char)
		}
	}
}

func TestStringCall(t *testing.T) {
	s := "hello"
	f := NewString([]rune(s))
	for i, c := range s {
		assert.Equal(t, c, rune(f.Call(NewNumber(float64(i))).(Number).Float64()))
	}
}

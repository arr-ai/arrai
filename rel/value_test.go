package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCall(t *testing.T) {
	t.Parallel()

	set := NewSet(
		NewTuple(NewAttr("@", NewNumber(1)), NewAttr("@foo", NewNumber(42))),
		NewTuple(NewAttr("@", NewNumber(1)), NewAttr("@foo", NewNumber(24))),
	)

	assert.Panics(t, func() { SetCall(set, NewNumber(1)) })
	assert.Panics(t, func() { SetCall(set, NewNumber(0)) })

	set = NewSet(
		NewTuple(NewAttr("@", NewNumber(1)), NewAttr("@foo", NewNumber(42))),
		NewTuple(NewAttr("@", NewNumber(2)), NewAttr("@foo", NewNumber(24))),
	)

	assert.True(t, SetCall(set, NewNumber(1)).Equal(NewNumber(42)))
	assert.True(t, SetCall(set, NewNumber(2)).Equal(NewNumber(24)))
}

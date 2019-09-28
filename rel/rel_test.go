package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testValueFromValue(t *testing.T, v Value) {
	v2, err := NewValue(v)
	assert.NoError(t, err)
	assert.Equal(t, v2, v)
}

func testNumericValue(t *testing.T, intf interface{}) {
	n, err := NewValue(intf)
	assert.NoError(t, err)
	assert.Equal(t, 123.0, n.Export())
	testValueFromValue(t, n)
}

func TestNewValueFromUint(t *testing.T) {
	t.Parallel()
	testNumericValue(t, uint(123))
}

func TestNewValueFromUint8(t *testing.T) {
	t.Parallel()
	testNumericValue(t, uint8(123))
}

func TestNewValueFromUint16(t *testing.T) {
	t.Parallel()
	testNumericValue(t, uint16(123))
}

func TestNewValueFromUint32(t *testing.T) {
	t.Parallel()
	testNumericValue(t, uint32(123))
}

func TestNewValueFromUint64(t *testing.T) {
	t.Parallel()
	testNumericValue(t, uint64(123))
}

func TestNewValueFromInt(t *testing.T) {
	t.Parallel()
	testNumericValue(t, int(123))
}

func TestNewValueFromInt8(t *testing.T) {
	t.Parallel()
	testNumericValue(t, int8(123))
}

func TestNewValueFromInt16(t *testing.T) {
	t.Parallel()
	testNumericValue(t, int16(123))
}

func TestNewValueFromInt32(t *testing.T) {
	t.Parallel()
	testNumericValue(t, int32(123))
}

func TestNewValueFromInt64(t *testing.T) {
	t.Parallel()
	testNumericValue(t, int64(123))
}

func TestNewValueFromFloat32(t *testing.T) {
	t.Parallel()
	testNumericValue(t, float32(123))
}

func TestNewValueFromFloat64(t *testing.T) {
	t.Parallel()
	testNumericValue(t, float64(123))
}

func TestNewValueFromMap(t *testing.T) {
	t.Parallel()
	m := map[string]interface{}{
		"a": 42.0,
		"b": map[string]interface{}{"c": 299792458.0},
	}
	v, err := NewValue(m)
	if assert.NoError(t, err) {
		assert.Equal(t, m, v.Export(), "%s.Export()", v)
		testValueFromValue(t, v)
	}
}

// func TestNewValueFromSlice(t *testing.T) {
// 	t.Parallel()
// 	m := []interface{}{
// 		42.0,
// 		43.0,
// 		[]interface{}{299792458.0},
// 	}
// 	v, err := NewValue(m)
// 	if assert.NoError(t, err) {
// 		assert.Equal(t, m, v.Export(), "%s.Export()", v)
//	}
// }

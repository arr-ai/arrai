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

// TestNewValueFromUint tests NewValue(uint).
func TestNewValueFromUint(t *testing.T) {
	testNumericValue(t, uint(123))
}

// TestNewValueFromUint8 tests NewValue(uint8).
func TestNewValueFromUint8(t *testing.T) {
	testNumericValue(t, uint8(123))
}

// TestNewValueFromUint16 tests NewValue(uint16).
func TestNewValueFromUint16(t *testing.T) {
	testNumericValue(t, uint16(123))
}

// TestNewValueFromUint32 tests NewValue(uint32).
func TestNewValueFromUint32(t *testing.T) {
	testNumericValue(t, uint32(123))
}

// TestNewValueFromUint64 tests NewValue(uint64).
func TestNewValueFromUint64(t *testing.T) {
	testNumericValue(t, uint64(123))
}

// TestNewValueFromInt tests NewValue(int).
func TestNewValueFromInt(t *testing.T) {
	testNumericValue(t, int(123))
}

// TestNewValueFromInt8 tests NewValue(int8).
func TestNewValueFromInt8(t *testing.T) {
	testNumericValue(t, int8(123))
}

// TestNewValueFromInt16 tests NewValue(int16).
func TestNewValueFromInt16(t *testing.T) {
	testNumericValue(t, int16(123))
}

// TestNewValueFromInt32 tests NewValue(int32).
func TestNewValueFromInt32(t *testing.T) {
	testNumericValue(t, int32(123))
}

// TestNewValueFromInt64 tests NewValue(int64).
func TestNewValueFromInt64(t *testing.T) {
	testNumericValue(t, int64(123))
}

// TestNewValueFromFloat32 tests NewValue(float32).
func TestNewValueFromFloat32(t *testing.T) {
	testNumericValue(t, float32(123))
}

// TestNewValueFromFloat64 tests NewValue(float64).
func TestNewValueFromFloat64(t *testing.T) {
	testNumericValue(t, float64(123))
}

// TestNewValueFromMap test NewValue(map[string]interface{})
func TestNewValueFromMap(t *testing.T) {
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

// // TestNewValueFromSlice test NewValue([]interface{})
// func TestNewValueFromSlice(t *testing.T) {
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

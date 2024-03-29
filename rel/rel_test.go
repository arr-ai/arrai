package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

func testValueFromValue(t *testing.T, v Value) {
	v2, err := NewValue(v)
	assert.NoError(t, err)
	assert.Equal(t, v2, v)
}

func testNumericValue(t *testing.T, intf interface{}) {
	n, err := NewValue(intf)
	assert.NoError(t, err)
	assert.Equal(t, 123.0, n.Export(arraictx.InitRunCtx(context.Background())))
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

func TestNewValueFromuintptr(t *testing.T) {
	t.Parallel()
	testNumericValue(t, uintptr(123))
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
		assert.Equal(t, m, v.Export(arraictx.InitRunCtx(context.Background())), "%s.Export()", v)
		testValueFromValue(t, v)
	}
}

// FIXME: uncomment when hashing works
// func TestRelationHash(t *testing.T) {
// 	t.Parallel()

// 	pr1 := projectedRelation{
// 		p: []int{0, 1},
// 		r: positionalRelation{
// 			frozen.NewSet(
// 				Values{NewNumber(1), NewNumber(2)},
// 				Values{NewNumber(3), NewNumber(4)},
// 			),
// 		},
// 	}
// 	pr2 := projectedRelation{
// 		p: []int{1, 0},
// 		r: positionalRelation{
// 			frozen.NewSet(
// 				Values{NewNumber(2), NewNumber(1)},
// 				Values{NewNumber(4), NewNumber(3)},
// 			),
// 		},
// 	}

// 	name1, name2 := []string{"x", "y"}, []string{"y", "x"}
// 	rel1 := Relation{attrs: name1, rows: pr1}
// 	rel2 := Relation{attrs: name1, rows: pr2}
// 	rel3 := Relation{attrs: name2, rows: pr1}
// 	rel4 := Relation{attrs: name2, rows: pr2}

// 	assert.Equal(t, rel1.Hash(0), rel2.Hash(0))
// 	assert.Equal(t, rel3.Hash(0), rel4.Hash(0))
// }

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

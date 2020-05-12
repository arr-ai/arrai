package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

// ArrayContains check if array a contains b, and b can be rel.Value or rel.Array.
func ArrayContains(a rel.Array, b rel.Value) rel.Value {
	switch b := b.(type) {
	case rel.Array:
		return arrayContainsArray(a, b)
	case rel.Value:
		arrayEnum, _ := a.ArrayEnumerator()
		if arrayEnum != nil {
			for arrayEnum.MoveNext() {
				if arrayEnum.Current().Equal(b) {
					return rel.NewBool(true)
				}
			}
		}
	}
	return rel.NewBool(false)
}

func arrayContainsArray(a, b rel.Array) rel.Value {
	// Get index of b[0] in a
	bOffset := 0
	bVals := b.Values()
	arrayEnum, _ := a.ArrayEnumerator()
	for arrayEnum.MoveNext() {
		if bOffset < len(bVals) && arrayEnum.Current().Equal(bVals[bOffset]) {
			bOffset++
		} else {
			if bOffset > 0 && bOffset < len(bVals) {
				return rel.NewBool(false)
			}
		}
	}

	if bOffset == len(bVals) {
		return rel.NewBool(true)
	}

	return rel.NewBool(false)
}

// ArraySub substitutes all b in a with c.
func ArraySub(a rel.Array, b, c rel.Value) rel.Value {
	// bArray := rel.NewArray(b)
	// cArray := rel.NewArray(c)
	// switch b := b.(type) {
	// case rel.Array:
	// 	return arrayContainsArray(a, b)
	// case rel.Value:
	return nil
}

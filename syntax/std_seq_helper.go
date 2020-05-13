package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

// ArrayContains check if array a contains b, and b can be rel.Value or rel.Array.
func ArrayContains(a rel.Array, b rel.Value) rel.Value {
	bArray := convert2Array(b)

	bOffset := 0
	bVals := bArray.Values()
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
	return nil
}

// ArrayPrefix check if a starts with b.
func ArrayPrefix(a rel.Array, b rel.Value) rel.Value {
	bArray := convert2Array(b)

	if bArray.Count() == 0 {
		return rel.NewBool(false)
	}
	if a.Count() < bArray.Count() {
		return rel.NewBool(false)
	}

	bVals := bArray.Values()
	bOffset := 0
	arrayEnum, _ := a.ArrayEnumerator()
	for arrayEnum.MoveNext() {
		if bOffset < bArray.Count() && arrayEnum.Current().Equal(bVals[bOffset]) {
			bOffset++
			if bOffset == bArray.Count() {
				break
			}
		} else {
			return rel.NewBool(false)
		}
	}

	return rel.NewBool(true)
}

// ArraySuffix check if a starts with b.
func ArraySuffix(a rel.Array, b rel.Value) rel.Value {
	bArray := convert2Array(b)

	if bArray.Count() == 0 {
		return rel.NewBool(false)
	}
	if a.Count() < bArray.Count() {
		return rel.NewBool(false)
	}

	aVals := a.Values()
	bVals := bArray.Values()
	bOffset := bArray.Count() - 1

	for i := a.Count() - 1; i >= 0; i-- {
		if bOffset > -1 && aVals[i].Equal(bVals[bOffset]) {
			bOffset--
			if bOffset == -1 {
				break
			}
		} else {
			return rel.NewBool(false)
		}
	}

	return rel.NewBool(true)
}

func convert2Array(val rel.Value) rel.Array {
	switch val := val.(type) {
	case rel.Array:
		return val
	case rel.Value:
		valArray, _ := rel.AsArray(rel.NewArray(val))
		return valArray
	}

	panic("it support types rel.Array and rel.Value only.")
}

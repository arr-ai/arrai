package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

// ArrayContains check if array a contains b, and b can be rel.Value or rel.Array.
func ArrayContains(a rel.Array, b rel.Value) rel.Value {
	bArray := convert2Array(b)
	return rel.NewBool(findFirstSequentialSubArray(a, bArray))
}

// It is brute force approach, can be improved later if it is necessary.
func findFirstSequentialSubArray(a, b rel.Array) bool {
	bOffset := 0
	aVals := a.Values()
	bVals := b.Values()

	for i := 0; i < len(aVals); i++ {
		if bOffset < len(bVals) && aVals[i].Equal(bVals[bOffset]) {
			bOffset++
		} else {
			if bOffset > 0 && bOffset < len(bVals) {
				bOffset = 0
				i--
			}
		}
	}

	return bOffset == len(bVals)
}

// ArraySub substitutes all b in a with c.
func ArraySub(a rel.Array, b, c rel.Value) rel.Value {
	return nil
}

// ArrayJoin joins array a to b, a is joiner and b is joinee.
func ArrayJoin(a rel.Array, b rel.Value) rel.Value {
	bArray := convert2Array(b)
	if bArray.Count() == 0 {
		// if joinee is empty, the final value will be empty
		return b
	}
	if a.Count() == 0 {
		return a
	}

	var vals []rel.Value = nil
	for i, value := range bArray.Values() {
		vals = append(vals, value)
		if i+1 < bArray.Count() {
			vals = append(vals, a.Values()...)
		}
	}

	return rel.NewArray(vals...)
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

	for _, val := range aVals[a.Count()-1:] {
		if bOffset > -1 && val.Equal(bVals[bOffset]) {
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
	case rel.GenericSet:
		valArray, _ := rel.AsArray(val)
		return valArray
	case rel.Value:
		valArray, _ := rel.AsArray(rel.NewArray(val))
		return valArray
	}

	panic("it support types rel.Array, rel.GenericSet and rel.Value only.")
}

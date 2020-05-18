package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

// arrayContains check if array a contains b, and b can be rel.Value or rel.Array.
func arrayContains(a rel.Array, b rel.Value) rel.Value {
	bArray := convert2Array(b)
	return rel.NewBool(indexSubArray(a.Values(), bArray.Values()) > -1)
}

// arraySub substitutes all old in a with new.
func arraySub(a rel.Array, old, new rel.Value) rel.Value {
	// Convert to array to facilitate process
	oldArray := convert2Array(old)
	newArray := convert2Array(new)

	finalVals := make([]rel.Value, 0, a.Count())
	for start, absoluteIndex := 0, 0; start < a.Count(); {
		relativeIndex := indexSubArray(a.Values()[start:], oldArray.Values())
		if relativeIndex >= 0 {
			absoluteIndex = relativeIndex + start
			if absoluteIndex-start > 0 {
				finalVals = append(finalVals, a.Values()[start:absoluteIndex]...)
			}
			finalVals = append(finalVals, newArray.Values()...)
			start = absoluteIndex + oldArray.Count()
		} else {
			finalVals = append(finalVals, a.Values()[absoluteIndex+1:]...)
			break
		}
	}

	return rel.NewArray(finalVals...)
}

// arraySplit split a by b.
func arraySplit(a rel.Array, b rel.Value) rel.Value {
	delimiter := convert2Array(b)
	var result []rel.Value

	if !delimiter.IsTrue() {
		for _, e := range a.Values() {
			result = append(result, rel.NewArray(e))
		}
	} else {
		for start, absoluteIndex := 0, 0; start < a.Count(); {
			relativeIndex := indexSubArray(a.Values()[start:], delimiter.Values())
			if relativeIndex >= 0 {
				absoluteIndex = relativeIndex + start
				if start != absoluteIndex {
					result = append(result, rel.NewArray(a.Values()[start:absoluteIndex]...))
				}
				start = absoluteIndex + delimiter.Count()
			} else {
				if start == 0 || start < a.Count() {
					result = append(result, rel.NewArray(a.Values()[start:]...))
				}
				break
			}
		}
	}

	return rel.NewArray(result...)
}

// arrayJoin joins array b to a, b is joiner and a is joinee.
func arrayJoin(a rel.Array, b rel.Value) rel.Value {
	joiner := convert2Array(b)
	if joiner.Count() == 0 || a.Count() == 0 {
		// if joinee is empty, the final value will be empty
		return a
	}

	vals := make([]rel.Value, 0, a.Count())
	for i, value := range a.Values() {
		switch vArray := value.(type) {
		case rel.Array:
			vals = append(vals, generate1LevelArray(vArray)...)
		case rel.Value:
			vals = append(vals, value)
		}

		if i+1 < a.Count() {
			vals = append(vals, generate1LevelArray(joiner)...)
		}
	}

	return rel.NewArray(vals...)
}

// arrayPrefix check if a starts with b.
func arrayPrefix(a rel.Array, b rel.Value) rel.Value {
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

// arraySuffix check if a ends with b.
func arraySuffix(a rel.Array, b rel.Value) rel.Value {
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

	panic("it supports types rel.Array, rel.GenericSet and rel.Value only.")
}

// It is brute force approach, can be improved later if it is necessary.
func indexSubArray(a, b []rel.Value) int {
	aOffset, bOffset := 0, 0

	for ; aOffset < len(a); aOffset++ {
		if bOffset < len(b) && a[aOffset].Equal(b[bOffset]) {
			bOffset++
		} else {
			if bOffset > 0 && bOffset < len(b) {
				bOffset = 0
				aOffset--
			}
		}
		if bOffset == len(b) {
			break
		}
	}

	if aOffset < len(a) {
		return aOffset
	}
	return -1
}

// Convert [[1, 2],[3, 4]] to [1, 2, 3, 4]
func generate1LevelArray(source rel.Array) []rel.Value {
	if source.Count() == 0 {
		return nil
	}

	finalArray := make([]rel.Value, 0, source.Count())
	for _, val := range source.Values() {
		switch rVal := val.(type) {
		case rel.Array:
			finalArray = append(finalArray, generate1LevelArray(rVal)...)
		case rel.Value:
			finalArray = append(finalArray, rVal)
		}
	}

	return finalArray
}

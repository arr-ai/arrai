package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

// Checks if array subject contains sub.
func arrayContains(sub rel.Value, subject rel.Array) rel.Value {
	subArray := convert2Array(sub)
	return rel.NewBool(search(subject.Values(), subArray.Values()) > -1)
}

// Substitutes all old in subject with new.
func arraySub(old, new rel.Value, subject rel.Array) rel.Value {
	// Convert to array to facilitate process
	oldArray := convert2Array(old)
	newArray := convert2Array(new)

	result := make([]rel.Value, 0, subject.Count())
	if !old.IsTrue() {
		for _, e := range subject.Values() {
			result = append(append(result, newArray.Values()...), e)
		}
		result = append(result, newArray.Values()...)
	} else {
		subjectVals := subject.Values()
		for {
			if i := search(subjectVals, oldArray.Values()); i >= 0 {
				result = append(append(result, subjectVals[:i]...), newArray.Values()...)
				subjectVals = subjectVals[i+oldArray.Count():]
			} else {
				result = append(result, subjectVals...)
				break
			}
		}
	}

	return rel.NewArray(result...)
}

// Splits array subject by delimiter.
func arraySplit(delimiter rel.Value, subject rel.Array) rel.Value {
	delimiterArray := convert2Array(delimiter)
	var result []rel.Value

	if !delimiterArray.IsTrue() {
		for _, e := range subject.Values() {
			result = append(result, rel.NewArray(e))
		}
	} else {
		subjectVals := subject.Values()
		for {
			if i := search(subjectVals, delimiterArray.Values()); i >= 0 {
				result = append(result, rel.NewArray(subjectVals[:i]...))
				subjectVals = subjectVals[i+delimiterArray.Count():]
			} else {
				result = append(result, rel.NewArray(subjectVals...))
				break
			}
		}
	}

	return rel.NewArray(result...)
}

// Joins array joiner to subject.
// The type of subject element must be rel.Array, it can help to make sure the API output is clear and will not confuse.
// For example:
//  `//seq.join([0], [1, 2])`, it can return [1, 0, 2] or [1, 2],
//  `//seq.join([], [1, [2, 3]])`, it can return [1, 2, 3] or [1, [2, 3]].
// All of the results make sense. It can see the output can't be sure in above cases, it is not good.
func arrayJoin(joiner rel.Value, subject rel.Array) rel.Value {
	joinerArray := convert2Array(joiner)

	result := make([]rel.Value, 0, subject.Count())
	for i, value := range subject.Values() {
		if i > 0 {
			result = append(result, joinerArray.Values()...)
		}
		switch vArray := value.(type) {
		case rel.Array:
			result = append(result, vArray.Values()...)
		case rel.Value:
			panic("the type of subject element must be rel.Array")
		}
	}

	return rel.NewArray(result...)
}

// Check if array subject starts with prefix.
func arrayHasPrefix(prefix rel.Value, subject rel.Array) rel.Value {
	if !prefix.IsTrue() {
		return rel.NewBool(true)
	}
	prefixArray := convert2Array(prefix)
	if subject.Count() < prefixArray.Count() {
		return rel.NewBool(false)
	}

	prefixVals := prefixArray.Values()
	prefixOffset := 0
	arrayEnum, _ := subject.ArrayEnumerator()
	for arrayEnum.MoveNext() {
		if prefixOffset < prefixArray.Count() && arrayEnum.Current().Equal(prefixVals[prefixOffset]) {
			prefixOffset++
			if prefixOffset == prefixArray.Count() {
				break
			}
		} else {
			return rel.NewBool(false)
		}
	}

	return rel.NewBool(true)
}

// Check if array subject ends with suffix.
func arrayHasSuffix(suffix rel.Value, subject rel.Array) rel.Value {
	suffixArray := convert2Array(suffix)

	if !suffixArray.IsTrue() {
		return rel.NewBool(true)
	}
	if subject.Count() < suffixArray.Count() {
		return rel.NewBool(false)
	}

	subjectVals := subject.Values()
	suffixVals := suffixArray.Values()
	suffixOffset := suffixArray.Count() - 1

	for _, val := range subjectVals[subject.Count()-1:] {
		if suffixOffset > -1 && val.Equal(suffixVals[suffixOffset]) {
			suffixOffset--
			if suffixOffset == -1 {
				break
			}
		} else {
			return rel.NewBool(false)
		}
	}

	return rel.NewBool(true)
}

func arrayTrimPrefix(prefix rel.Value, subject rel.Array) rel.Value {
	if prefix.IsTrue() && subject.IsTrue() && arrayHasPrefix(prefix, subject).IsTrue() {
		switch result := rel.Difference(subject, prefix.(rel.Array)).(type) {
		case rel.Array:
			return result.Shift(-prefix.(rel.Array).Count())
		case rel.Set:
			return result
		}
	}
	return subject
}

func arrayTrimSuffix(suffix rel.Value, subject rel.Array) rel.Value {
	if suffix.IsTrue() && subject.IsTrue() && arrayHasSuffix(suffix, subject).IsTrue() {
		suffix := suffix.(rel.Array)
		return rel.Difference(subject, suffix.Shift(subject.Count()-suffix.Count()))
	}
	return subject
}

func convert2Array(val rel.Value) rel.Array {
	if array, is := rel.AsArray(val.(rel.Set)); is {
		return array
	}
	panic("it supports types rel.Array and rel.GenericSet only.")
}

// Searches array sub in subject and return the first indedx if found, or return -1.
// It is brute force approach, can be improved later if it is necessary.
// Case: subject=[1,2,3,4], sub=[2], return 1
// Case: subject=[1,2,3,4], sub=[2,3], return 1
// Case: subject=[1,2,3,4], sub=[2,5], return -1
func search(subject, sub []rel.Value) int {
	subjectOffset, subOffset := 0, 0

	for ; subjectOffset < len(subject); subjectOffset++ {
		if subOffset < len(sub) && subject[subjectOffset].Equal(sub[subOffset]) {
			subOffset++
		} else {
			if subOffset > 0 && subOffset < len(sub) {
				subOffset = 0
				subjectOffset--
			}
		}
		if subOffset == len(sub) {
			break
		}
	}

	if subjectOffset < len(subject) {
		// see len(sub) > 1
		return (subjectOffset + 1) - len(sub)
	}
	return -1
}

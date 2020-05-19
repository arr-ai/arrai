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

	if !oldArray.IsTrue() && !new.IsTrue() {
		return subject
	}

	result := make([]rel.Value, 0, subject.Count())
	if !old.IsTrue() {
		for _, e := range subject.Values() {
			result = append(append(result, newArray.Values()...), e)
		}
		result = append(result, newArray.Values()...)
	} else {
		for start, absoluteIndex := 0, 0; start < subject.Count(); {
			relativeIndex := search(subject.Values()[start:], oldArray.Values())
			if relativeIndex >= 0 {
				absoluteIndex = relativeIndex + start
				if absoluteIndex-start > 0 {
					result = append(result, subject.Values()[start:absoluteIndex]...)
				}
				result = append(result, newArray.Values()...)
				start = absoluteIndex + oldArray.Count()
			} else {
				result = append(result, subject.Values()[absoluteIndex+1:]...)
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
		for start, absoluteIndex := 0, 0; start < subject.Count(); {
			relativeIndex := search(subject.Values()[start:], delimiterArray.Values())
			if relativeIndex >= 0 {
				absoluteIndex = relativeIndex + start
				if start != absoluteIndex {
					result = append(result, rel.NewArray(subject.Values()[start:absoluteIndex]...))
				}
				start = absoluteIndex + delimiterArray.Count()
			} else {
				if start == 0 || start < subject.Count() {
					result = append(result, rel.NewArray(subject.Values()[start:]...))
				}
				break
			}
		}
	}

	return rel.NewArray(result...)
}

// Joins array joiner to subject.
func arrayJoin(joiner rel.Value, subject rel.Array) rel.Value {
	joinerArray := convert2Array(joiner)
	if !joinerArray.IsTrue() || !subject.IsTrue() {
		return subject
	}

	result := make([]rel.Value, 0, subject.Count())
	for i, value := range subject.Values() {
		if i > 0 {
			result = append(result, joinerArray.Values()...)
		}
		switch vArray := value.(type) {
		case rel.Array:
			result = append(result, vArray.Values()...)
		case rel.Value:
			result = append(result, value)
		}
	}

	return rel.NewArray(result...)
}

// Check if array subject starts with prefix.
func arrayHasPrefix(prefix rel.Value, subject rel.Array) rel.Value {
	prefixArray := convert2Array(prefix)

	if !prefixArray.IsTrue() && subject.IsTrue() {
		return rel.NewBool(true)
	}
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

	if !suffixArray.IsTrue() && subject.IsTrue() {
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

func convert2Array(val rel.Value) rel.Array {
	switch val := val.(type) {
	case rel.Array:
		return val
	case rel.GenericSet:
		valArray, _ := rel.AsArray(val)
		return valArray
	}

	panic("it supports types rel.Array and rel.GenericSet only.")
}

// Searches array sub in subject and return the first indedx if found, or return -1.
// It is brute force approach, can be improved later if it is necessary.
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
		return subjectOffset
	}
	return -1
}

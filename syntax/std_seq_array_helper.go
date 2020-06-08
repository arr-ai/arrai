package syntax

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

// Checks if array subject contains sub.
func arrayContains(sub rel.Value, subject rel.Array) (rel.Value, error) {
	if subArray, is := rel.AsArray(sub); is {
		return rel.NewBool(search(subject.Values(), subArray.Values()) > -1), nil
	}
	return nil, fmt.Errorf("//seq.contains: sub not an array: %v", sub)
}

// Substitutes all old in subject with new.
func arraySub(old, new rel.Value, subject rel.Array) (rel.Value, error) {
	// Convert to array to facilitate process
	oldArray, is := rel.AsArray(old)
	if !is {
		return nil, fmt.Errorf("//seq.sub: old not an array: %v", old)
	}
	newArray, is := rel.AsArray(new)
	if !is {
		return nil, fmt.Errorf("//seq.sub: new not an array: %v", new)
	}

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

	return rel.NewArray(result...), nil
}

// Splits array subject by delimiter.
func arraySplit(delimiter rel.Value, subject rel.Array) (rel.Value, error) {
	delimiterArray, is := rel.AsArray(delimiter)
	if !is {
		return nil, fmt.Errorf("//seq.split: delimiter not an array: %v", delimiter)
	}
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

	return rel.NewArray(result...), nil
}

// Joins array joiner to subject.
// The type of subject element must be rel.Array, it can help to make sure the API output is clear and will not confuse.
// For example:
//  `//seq.join([0], [1, 2])`, it can return [1, 0, 2] or [1, 2],
//  `//seq.join([], [1, [2, 3]])`, it can return [1, 2, 3] or [1, [2, 3]].
// All of the results make sense. It can see the output can't be sure in above cases, it is not good.
func arrayJoin(joiner rel.Value, subject rel.Array) (rel.Value, error) {
	joinerArray, is := rel.AsArray(joiner)
	if !is {
		return nil, fmt.Errorf("//seq.join: joiner not an array: %v", joiner)
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
			return nil, fmt.Errorf("//seq.join: the type of subject element must be rel.Array")
		}
	}

	return rel.NewArray(result...), nil
}

// Check if array subject starts with suffix.
func arrayHasPrefix(prefix rel.Value, subject rel.Array) (rel.Value, error) {
	if !prefix.IsTrue() {
		return rel.NewBool(true), nil
	}
	prefixArray, is := rel.AsArray(prefix)
	if !is {
		return nil, fmt.Errorf("//seq.has_prefix: prefix not an array: %v", prefix)
	}
	if subject.Count() < prefixArray.Count() {
		return rel.NewBool(false), nil
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
			return rel.NewBool(false), nil
		}
	}

	return rel.NewBool(true), nil
}

// Check if array subject ends with suffix.
func arrayHasSuffix(suffix rel.Value, subject rel.Array) (rel.Value, error) {
	suffixArray, is := rel.AsArray(suffix)
	if !is {
		return nil, fmt.Errorf("//seq.has_suffix: suffix not an array: %v", suffix)
	}

	if !suffixArray.IsTrue() {
		return rel.NewBool(true), nil
	}
	if subject.Count() < suffixArray.Count() {
		return rel.NewBool(false), nil
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
			return rel.NewBool(false), nil
		}
	}

	return rel.NewBool(true), nil
}

func arrayTrimPrefix(prefix rel.Value, subject rel.Array) (rel.Value, error) {
	if prefix.IsTrue() && subject.IsTrue() {
		hasPrefix, err := arrayHasPrefix(prefix, subject)
		if err != nil {
			return nil, err
		}
		if hasPrefix.IsTrue() {
			switch result := rel.Difference(subject, prefix.(rel.Array)).(type) {
			case rel.Array:
				return result.Shift(-prefix.(rel.Array).Count()), nil
			case rel.Set:
				return result, nil
			}
		}
	}
	return subject, nil
}

func arrayTrimSuffix(suffix rel.Value, subject rel.Array) (rel.Value, error) {
	if suffix.IsTrue() && subject.IsTrue() {
		hasSuffix, err := arrayHasSuffix(suffix, subject)
		if err != nil {
			return nil, err
		}
		if hasSuffix.IsTrue() {
			suffixArray := suffix.(rel.Array)
			return rel.Difference(subject, suffixArray.Shift(subject.Count()-suffixArray.Count())), nil
		}
	}
	return subject, nil
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

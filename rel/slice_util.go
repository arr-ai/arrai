package rel

import "math"

// resolveArrayIndexes returns an array of indexes to get for array.
func resolveArrayIndexes(start, end Value, step, offset, maxLen int, inclusive bool) []int {
	if maxLen == 0 {
		return []int{}
	}
	startIndex, endIndex := initDefaultArrayIndex(start, end, offset, maxLen+offset, step)

	if startIndex == endIndex {
		if inclusive {
			return []int{startIndex}
		}
		return []int{}
	}

	return getIndexes(startIndex, endIndex, step, inclusive)
}

// initDefaultArrayIndex initialize the start and end values for arrays.
func initDefaultArrayIndex(start, end Value, minLen, maxLen, step int) (startIndex int, endIndex int) {
	if start != nil {
		startIndex = resolveIndex(int(start.(Number)), minLen, maxLen)
		if startIndex == maxLen {
			startIndex--
		}
	} else {
		if step > 0 {
			startIndex = minLen
		} else {
			startIndex = maxLen - 1
		}
	}

	if end != nil {
		endIndex = resolveIndex(int(end.(Number)), minLen, maxLen)
	} else {
		// TODO: apply inclusivity to the undefined end index
		if step > 0 {
			endIndex = maxLen
		} else {
			endIndex = minLen - 1
		}
	}
	return
}

// resolveIndex solves the edge cases of index values.
func resolveIndex(i, minVal, maxVal int) int {
	if i > maxVal {
		return maxVal
	} else if i < 0 {
		if -i > maxVal {
			return minVal
		}
		return maxVal + i
	}
	return i
}

// getIndexes returns a range of numbers between start and end with the provided step.
// inclusive decides whether end can be included or not.
func getIndexes(start, end, step int, inclusive bool) []int {
	if !isValidRange(start, end, step) {
		return []int{}
	}
	if inclusive {
		if step > 0 {
			end++
		} else {
			end--
		}
	}

	length := int(math.Abs(float64(start - end)))
	if step != 1 && step != -1 {
		length = int(math.Ceil(float64(length) / math.Abs(float64(step))))
	}
	indexes := make([]int, 0, length)
	for i := 0; i < length; i++ {
		indexes = append(indexes, start+step*i)
	}

	return indexes
}

// isValidRange checks whether start, end, and step are valid values.
func isValidRange(start, end, step int) bool {
	return step != 0 && ((start > end && step < 0) || (start < end && step > 0))
}

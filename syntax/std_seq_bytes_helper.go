package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

// BytesJoin join b to a
func BytesJoin(a, b rel.Bytes) rel.Value {
	result := make([]byte, 0, a.Count())
	for index, e := range a.Bytes() {
		if index > 0 && index < a.Count() {
			result = append(result, b.Bytes()...)
		}
		result = append(result, e)
	}

	return rel.NewBytes(result)
}

// BytesContain check if a contains b.
func BytesContain(a, b rel.Bytes) rel.Value {
	return rel.NewBool(indexSubBytes(a.Bytes(), b.Bytes()) > -1)
}

// BytesSub substitute all old in a with new.
func BytesSub(a, old, new rel.Bytes) rel.Value {
	finalVals := make([]byte, 0, a.Count())

	for start, absoluteIndex := 0, 0; start < a.Count(); {
		relativeIndex := indexSubBytes(a.Bytes()[start:], old.Bytes())
		if relativeIndex >= 0 {
			absoluteIndex = relativeIndex + start
			if absoluteIndex-start > 0 {
				finalVals = append(finalVals, a.Bytes()[start:absoluteIndex]...)
			}
			finalVals = append(finalVals, new.Bytes()...)
			start = absoluteIndex + old.Count()
		} else {
			finalVals = append(finalVals, a.Bytes()[absoluteIndex+1:]...)
			break
		}
	}

	return rel.NewBytes(finalVals)
}

// It is brute force approach, can be improved later if it is necessary.
func indexSubBytes(a, b []byte) int {
	aOffset, bOffset := 0, 0

	for ; aOffset < len(a); aOffset++ {
		if bOffset < len(b) && a[aOffset] == b[bOffset] {
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

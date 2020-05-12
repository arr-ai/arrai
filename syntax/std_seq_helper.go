package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

// ContainsArray check if array a contains array b.
func ContainsArray(a rel.Array, b rel.Array) rel.Value {
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

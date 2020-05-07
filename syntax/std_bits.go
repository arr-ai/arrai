package syntax

import (
	"math"
	"math/bits"

	"github.com/arr-ai/arrai/rel"
)

func stdBits() rel.Attr {
	return rel.NewTupleAttr("bits",
		rel.NewNativeFunctionAttr("mask", mask),
		rel.NewNativeFunctionAttr("set", set),
	)
}

func set(v rel.Value) rel.Value {
	n, isNumber := v.(rel.Number)
	if !isNumber || float64(n) < 0 {
		panic("argument has to be non-negative number")
	}
	var bitmask []rel.Value
	if float64(n) == float64(int(n)) {
		bitmask = setInt(int(n))
	} else {
		bitmask = maskFloat(float64(n))
	}
	return rel.NewSet(bitmask...)
}

func setInt(v int) (bitmask []rel.Value) {
	for ; v != 0; v &= (v - 1) {
		bitmask = append(bitmask, rel.NewNumber(float64(bits.TrailingZeros64(uint64(v)))))
	}
	return
}

func maskFloat(v float64) (bitmask []rel.Value) {
	panic("unimplemented")
}

func mask(v rel.Value) rel.Value {
	if s, isSet := v.(rel.Set); isSet {
		var n float64
		for e := s.Enumerator(); e.MoveNext(); {
			n += math.Pow(2, float64(e.Current().(rel.Number)))
		}
		return rel.NewNumber(n)
	}
	panic("argument has to be a Set")
}

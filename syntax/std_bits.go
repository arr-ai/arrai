package syntax

import (
	"context"
	"fmt"
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

func set(_ context.Context, v rel.Value) (rel.Value, error) {
	n, isNumber := v.(rel.Number)
	if !isNumber || float64(n) < 0 {
		return nil, fmt.Errorf("argument not a non-negative number: %v", v)
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
	for ; v != 0; v &= v - 1 {
		bitmask = append(bitmask, rel.NewNumber(float64(bits.TrailingZeros64(uint64(v)))))
	}
	return
}

func maskFloat(v float64) (bitmask []rel.Value) {
	panic("unimplemented")
}

func mask(_ context.Context, v rel.Value) (rel.Value, error) {
	if s, isSet := v.(rel.Set); isSet {
		var total float64
		for e := s.Enumerator(); e.MoveNext(); {
			n, is := e.Current().(rel.Number)
			if !is {
				return nil, fmt.Errorf("//bits.mask: element not a number: %v", e.Current())
			}
			total += math.Pow(2, n.Float64())
		}
		return rel.NewNumber(total), nil
	}
	return nil, fmt.Errorf("arg to mask must be a set, not %s", rel.ValueTypeAsString(v))
}

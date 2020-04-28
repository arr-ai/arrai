package rel

import (
	"container/heap"

	"github.com/arr-ai/wbnf/parser"
	"github.com/pkg/errors"
)

type float64Heap []float64

func (h float64Heap) Len() int           { return len(h) }
func (h float64Heap) Less(i, j int) bool { return h[i] > h[j] }
func (h float64Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *float64Heap) Push(x interface{}) {
	*h = append(*h, x.(float64))
}

func (h *float64Heap) Pop() interface{} {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[0 : n-1]
	return x
}

// NewMedianExpr evaluates to the median of expr over all elements in a.
func NewMedianExpr(scanner parser.Scanner, a, b Expr) Expr {
	type Agg struct {
		h float64Heap
		n int
	}
	return NewReduceExpr(
		scanner, a, ExprAsFunction(b), "%s min ???",
		func(s Set) (interface{}, error) {
			if n := s.Count(); n > 0 {
				return Agg{h: make(float64Heap, 0, n/2+2), n: n}, nil
			}
			return nil, errors.Errorf("Empty set has no mean")
		},
		func(acc interface{}, v Value) (interface{}, error) {
			agg := acc.(Agg)
			switch v := v.(type) {
			case Number:
				heap.Push(&agg.h, v.Float64())
				if len(agg.h) == cap(agg.h) {
					heap.Pop(&agg.h)
				}
				return agg, nil
			}
			return nil, errors.Errorf("Non-numeric value used in mean")
		},
		func(acc interface{}) (Value, error) {
			if acc != nil {
				agg := acc.(Agg)
				a := heap.Pop(&agg.h).(float64)
				if agg.n%2 == 0 {
					b := heap.Pop(&agg.h).(float64)
					return NewNumber(0.5 * (a + b)), nil
				}
				return NewNumber(a), nil
			}
			return nil, errors.Errorf("Empty input to min")
		},
	)
}

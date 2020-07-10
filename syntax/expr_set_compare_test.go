package syntax

import (
	"fmt"
	"math/bits"
	"testing"

	"github.com/arr-ai/arrai/rel"
)

func TestSetCompare(t *testing.T) {
	t.Parallel()

	intSet := func(i int) rel.Set {
		set := rel.None
		for ; i != 0; i &= i - 1 {
			set = set.With(rel.NewNumber(float64(bits.TrailingZeros(uint(i)))))
		}
		return set
	}
	for i := 0; i < 8; i++ {
		i := i
		a := intSet(i)
		for j := 0; j < 8; j++ {
			j := j
			b := intSet(j)
			t.Run(fmt.Sprintf("%v.%v", a, b), func(t *testing.T) {
				t.Parallel()

				assertComparison := func(op string, result bool) bool { //nolint:unparam
					var expected string
					if result {
						expected = `true`
					} else {
						expected = `false`
					}
					return AssertCodesEvalToSameValue(t, expected, fmt.Sprintf("%v %s %v", a, op, b))
				}
				assertComparison(`(<)`, i&^j == 0 && i != j)
				assertComparison(`(<=)`, i&^j == 0)
				assertComparison(`(>)`, j&^i == 0 && j != i)
				assertComparison(`(>=)`, j&^i == 0)
				assertComparison(`(<>)`, (i&^j == 0 || j&^i == 0) && i != j)
				assertComparison(`(<>=)`, i&^j == 0 || j&^i == 0)

				assertComparison(`!(<)`, !(i&^j == 0 && i != j))
				assertComparison(`!(<=)`, !(i&^j == 0))
				assertComparison(`!(>)`, !(j&^i == 0 && j != i))
				assertComparison(`!(>=)`, !(j&^i == 0))
				assertComparison(`!(<>)`, !((i&^j == 0 || j&^i == 0) && i != j))
				assertComparison(`!(<>=)`, !(i&^j == 0 || j&^i == 0))
			})
		}
	}
}

func TestStringCompare(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `true `, `"abc" = "abc"  `)
	AssertCodesEvalToSameValue(t, `false`, `"abc" = "ab"   `)
	AssertCodesEvalToSameValue(t, `false`, `"abc" = 2\"abc"`)
	AssertCodesEvalToSameValue(t, `false`, `"abc" = {}     `)
}

func TestBytesCompare(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `true `, `<<"abc">> = <<"abc">>  `)
	AssertCodesEvalToSameValue(t, `false`, `<<"abc">> = <<"ab">>   `)
	AssertCodesEvalToSameValue(t, `false`, `<<"abc">> = 2\<<"abc">>`)
	AssertCodesEvalToSameValue(t, `false`, `<<"abc">> = {}         `)
}

func TestArrayCompare(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `true `, `[1, 2, 3] = [1, 2, 3]   `)
	AssertCodesEvalToSameValue(t, `false`, `[1, 2, 3] = [1, 2]      `)
	AssertCodesEvalToSameValue(t, `false`, `[1, 2, 3] = -1\[1, 2, 3]`)
	AssertCodesEvalToSameValue(t, `false`, `[1, 2, 3] = {}          `)
}

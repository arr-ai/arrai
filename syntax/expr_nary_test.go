package syntax

import (
	"fmt"
	"testing"
)

func TestComparisons(t *testing.T) {
	t.Parallel()
	ops := map[string]func(a, b int) bool{
		"=":  func(a, b int) bool { return a == b },
		"!=": func(a, b int) bool { return a != b },
		"<":  func(a, b int) bool { return a < b },
		">=": func(a, b int) bool { return a >= b },
	}
	for a, fa := range ops {
		for b, fb := range ops {
			x := 1
			for y := 0; y < 3; y++ {
				for z := 0; z < 3; z++ {
					var expected string
					if fa(x, y) && fb(y, z) {
						expected = `true`
					} else {
						expected = `false`
					}
					AssertCodesEvalToSameValue(t, expected, fmt.Sprintf("%d %s %d %s %d", x, a, y, b, z))
				}
			}
		}
	}
}

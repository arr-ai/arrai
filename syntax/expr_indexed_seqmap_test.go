package syntax

import "testing"

func TestExprIndexedSequenceMap(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[1, 3, 5]`, `[1, 2, 3] >>> \i \n i + n`)
	AssertCodesEvalToSameValue(t,
		`{
			3      : ( "key": 3      , "val": (2)      ),
			"ten"  : ( "key": "ten"  , "val": 10       ),
			"stuff": ( "key": "stuff", "val": "random" ),
		}`,
		`{"stuff": "random", "ten": 10, 3: (2)} >>> \i \n ("key": i, "val": n)`,
	)
	AssertCodeErrors(t,
		`{("a": "z"), ("b": "y")} >>> \i \n (i ++ n)`,
		`=> not applicable to unindexed type {(a: z), (b: y)}`,
	)
}

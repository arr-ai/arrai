package syntax

import (
	"testing"
)

func TestEvalCond(t *testing.T) {
	t.Parallel()
	// AssertCodesEvalToSameValue(t, `1`, `cond {(1 > 0) : 1, (2 > 3): 2, _:1 + 2,}`)
	// AssertCodesEvalToSameValue(t, `1`, `cond (1 > 0 : 1, 2 > 3: 2, *:1 + 2)`)
	// AssertCodesEvalToSameValue(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2, *:1 + 2)`)
	// AssertCodesEvalToSameValue(t, `3`, `cond (1 > 2 : 1, 2 > 3: 2, *:1 + 2)`)
	// AssertCodesEvalToSameValue(t, `1`, `cond {(1 < 2) : 1,}`)
	AssertCodesEvalToSameValue(t, `1`, `cond {(1 < 2) : 1}`)
	// AssertCodesEvalToSameValue(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2)`)
	// AssertCodesEvalToSameValue(t, `3`, `cond (* : 1 + 2)`)
	// AssertCodesEvalToSameValue(t, `1`, `cond (1 < 2: 1, * : 1 + 2)`)
	// AssertCodesEvalToSameValue(t, `3`, `cond (1 > 2: 1, * : 1 + 2)`)
	// AssertCodesEvalToSameValue(t, `3`, `let a = cond (1 > 2: 1, * : 1 + 2);a`)
	// AssertCodesEvalToSameValue(t, `3`, `let a = cond (1 > 2: 1, * : 1 + 2,);a`)
	// AssertCodesEvalToSameValue(t, `1`, `let a = cond (1 < 2: 1, * : 1 + 2);a * 1`)
	// // Multiple true conditions
	// AssertCodesEvalToSameValue(t, `1`, `cond (1 > 0 : 1, 2 < 3: 2, *:1 + 2)`)
	// AssertCodesEvalToSameValue(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2, *:1 + 2)`)
	// AssertCodesEvalToSameValue(t, `3`, `cond (1 > 2 : 1, 2 > 3: 2, 3 < 4 :3, *:2 + 2)`)
	// AssertCodesEvalToSameValue(t, `3`, `cond (1 > 2 : 1, 2 > 3: 2, 3 < 4 :3, 4 < 5 : 5, *:2 + 2)`)
	// AssertCodesEvalToSameValue(t, `3`, `cond (1 > 2 : 1, 2 > 3: 2, 3 < 4 :3, 4 < 5 : 5, 5 > 6 : 6, *:2 + 2)`)
	// AssertCodesEvalToSameValue(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2, 3 < 4 :3, 4 < 5 : 5, 5 > 6 : 6, *:2 + 2)`)

	// // Nested call
	// AssertCodesEvalToSameValue(t, `1`, `cond (cond (1 > 0 : 1) > 0 : 1, 2 < 3: 2, *:1 + 2)`)
	// AssertCodesEvalToSameValue(t, `2`, `cond (cond (1 > 2 : 1, * : 11) < 2 : 1, 2 < 3: 2, *:1 + 2)`)
	// AssertCodesEvalToSameValue(t, `20`, `let a = cond (cond (1 > 2 : 1, * : 11) < 2 : 1, 2 < 3: 2, *:1 + 2);a * 10`)

	// AssertCodesEvalToSameValue(t, ``, `cond (1 < 0 : 1, 2 > 3: 2)`)
	// AssertCodesEvalToSameValue(t, ``, `cond (1 < 0 : 1)`)
	// AssertCodesEvalToSameValue(t, ``, `cond ()`)
}

func TestEvalCondWithControlVar(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `1`, `let a = 1; a cond {(1) :1, (2) :2}`)
	// AssertCodesEvalToSameValue(t, `1`, `let [a, b, c , _] = [1,2,3,4]; [a,b]`)

	// Control var conditions
	AssertCodesEvalToSameValue(t, `1`, `let a = 1; a cond {(1) :1}`)
	AssertCodesEvalToSameValue(t, `1`, `let a = 1; a cond {(1) :1, (2) :2}`)
	AssertCodesEvalToSameValue(t, `1`, `let a = 1; a cond {(1) :1, (2) :2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `11`, `let a = 1; a cond {(1) :1 + 10, (2) : 2, _:1 + 2}`)
	// AssertCodesEvalToSameValue(t, `3`, `let a = 1; a cond {(2) :2, _:1 + 2}`)
	// AssertCodesEvalToSameValue(t, `3`, `let a = 1; a cond (2 :2, _:1 + 2)`)
	// AssertCodesEvalToSameValue(t, `13`, `let a = 1; a cond {(2) :2, _:11 + 2}`)
	// AssertCodesEvalToSameValue(t, `5`, `let a = 3; a cond {_:5}`)
	// AssertCodesEvalToSameValue(t, `3`, `let a = 3; a cond {_:1 + 2}`)

	AssertCodesEvalToSameValue(t, `1`, `let a = 1; let b = a cond {(1) :1, (2) :2, _:1 + 2}; b`)
	AssertCodesEvalToSameValue(t, `100`, `let a = 1; let b = a cond {(1) :1, (2) :2, _:1 + 2}; b * 100`)
	// //
	// AssertCodesEvalToSameValue(t, `2`, `let a = 1; (a + 1) cond {(1) :1, (2) :2, _:1 + 2}`)
	// AssertCodesEvalToSameValue(t, `3`, `let a = 1; (a + 10) cond {(1) :1, (2) :2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = 1; let b = (a + 1) cond {(1) :1, (2) :2, _:1 + 2}; b`)
	// AssertCodesEvalToSameValue(t, `300`, `let a = 1; let b = (a + 10) cond {(1) :1, (2) :2, _:1 + 2}; b * 100`)
	// Nested call
	// AssertCodesEvalToSameValue(t, "B", `let a = 2; a cond { (a cond {(1,2) : 1}): "A", (2, 3): "B", _: "C"}`)
	// AssertCodesEvalToSameValue(t, "A", `let a = 1; a cond { (cond {(2 > 1) : 1}): "A", (2, 3): "B", _: "C"}`)
	// AssertCodesEvalToSameValue(t, "A", `let a = 1; cond { (a cond {(1) : 1}) : "A", (2): "B", _: "C"}`)

	// AssertCodesEvalToSameValue(t, ``, `let a = 3; a cond {(1) :1, (2) :2 + 1}`)
	// AssertCodesEvalToSameValue(t, ``, `let a = 3; let b = a cond {(1) :1, (2) :2 + 1}; b`)
	// AssertCodesEvalToSameValue(t, ``, `let a = 3; let b = (a + 10) cond {(1) :1, (2) :2 + 1}; b`)
}

package syntax

import "testing"

func TestExprLetIdentPattern(t *testing.T) {
	AssertCodesEvalToSameValue(t, `7`, `let x = 6; 7`)
	AssertCodesEvalToSameValue(t, `42`, `let x = 6; x * 7`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let x = 1; [x, 2]`)
	AssertCodesEvalToSameValue(t, `2`, `let x = 1; let x = x + 1; x`)
}

func TestExprLetValuePattern(t *testing.T) {
	AssertCodesEvalToSameValue(t, `42`, `let 42 = 42; 42`)
	AssertCodesEvalToSameValue(t, `1`, `let 42 = 42; 1`)
	AssertCodePanics(t, `let 42 = 1; 42`)
	AssertCodePanics(t, `let 42 = 1; 1`)
}

func TestExprLetArrayPattern(t *testing.T) {
	AssertCodesEvalToSameValue(t, `1`, `let [] = []; 1`)
	AssertCodesEvalToSameValue(t, `9`, `let [a, b, c] = [1, 2, 3]; 9`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `let [a, b, c] = [1, 2, 3]; [a, b, c]`)
	AssertCodesEvalToSameValue(t, `2`, `let [a, b, c] = [1, 2, 3]; [a, b, c](1)`)
	AssertCodesEvalToSameValue(t, `2`, `let [a, b, c] = [1, 2, 3]; b`)
	AssertCodesEvalToSameValue(t, `2`, `let [a, b, c] = [1, 2, 3]; [c, b](1)`)
	AssertCodesEvalToSameValue(t, `2`, `let arr = [1, 2]; let [a, b] = arr; b`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `let [[x, y], z] = [[1, 2], 3]; [x, y, z]`)
	AssertCodesEvalToSameValue(t, `1`, `let [x, x] = [1, 1]; x`)
	AssertCodesEvalToSameValue(t, `1`, `let [x, _, _] = [1, 2, 3]; x`)
	AssertCodesEvalToSameValue(t, `2`, `let [_, x, _] = [1, 2, 3]; x`)
	AssertCodesEvalToSameValue(t, `3`, `let x = 3; let [(x)] = [3]; x`)
	AssertCodesEvalToSameValue(t, `2`, `let x = 3; let [b, (x)] = [2, 3]; b`)
	AssertCodesEvalToSameValue(t, `2`, `let x = 3; let [_, b, (x)] = [1, 2, 3]; b`)
	AssertCodesEvalToSameValue(t, `2`, `let x = 3; let [x] = [2]; x`)

	AssertCodePanics(t, `let [(x)] = [2]; x`)
	AssertCodePanics(t, `let x = 3; let [(x)] = [2]; x`)
	AssertCodePanics(t, `let [x, y] = 1; x`)
	AssertCodePanics(t, `let [x, x] = [1]; x`)
	AssertCodePanics(t, `let [x, y] = [1]; x`)
	AssertCodePanics(t, `let [x, x] = [1, 2]; x`)
	AssertCodeErrors(t, `let [_] = [1]; _`,
		"Name \"_\" not found in {} \n\n\x1b[1;37m:1:16:\x1b[0m\nlet [_]")
	AssertCodePanics(t, `let x = 3; let [b, (x)] = [2, 1]; b`)
}

func TestExprLetTuplePattern(t *testing.T) {
	AssertCodesEvalToSameValue(t, `4`, `let () = (); 4`)
	AssertCodesEvalToSameValue(t, `4`, `let (a:x, b:y) = (a:4, b:7); x`)
	AssertCodesEvalToSameValue(t, `4`, `let (a:x, b:x) = (a:4, b:4); x`)
	AssertCodesEvalToSameValue(t, `4`, `let x = 4; let (a:x) = (a:4); x`)
	AssertCodesEvalToSameValue(t, `4`, `let x = 5; let (a:x) = (a:4); x`)
	AssertCodePanics(t, `let (a:x) = (b:7, a:4); x`)
	AssertCodePanics(t, `let (a:x, a:x) = (a:4, a:4); x`)
	AssertCodePanics(t, `let (a:x, a:x) = (a:4); x`)
	AssertCodePanics(t, `let x = 5; let (a:(x)) = (a:4); x`)
	AssertCodePanics(t, `let (a:x, b:x) = (a:4, b:7); x`)
}

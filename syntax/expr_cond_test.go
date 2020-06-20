package syntax

import (
	"testing"
)

func TestEvalCond(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `1`, `cond {1 > 0 : 1, 2 > 3: 2, _:1 + 2,}`)
	AssertCodesEvalToSameValue(t, `1`, `cond {1 > 0 : 1, 2 > 3: 2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `2`, `cond {1 > 2 : 1, 2 < 3: 2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `3`, `cond {1 > 2 : 1, 2 > 3: 2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `1`, `cond {1 < 2 : 1,}`)
	AssertCodesEvalToSameValue(t, `1`, `cond {1 < 2 : 1}`)
	AssertCodesEvalToSameValue(t, `2`, `cond {1 > 2 : 1, 2 < 3: 2}`)
	AssertCodesEvalToSameValue(t, `3`, `cond {_ : 1 + 2}`)
	AssertCodesEvalToSameValue(t, `1`, `cond {1 < 2: 1, _ : 1 + 2}`)
	AssertCodesEvalToSameValue(t, `3`, `cond {1 > 2: 1, _ : 1 + 2}`)
	AssertCodesEvalToSameValue(t, `3`, `let a = cond {1 > 2: 1, _ : 1 + 2};a`)
	AssertCodesEvalToSameValue(t, `3`, `let a = cond {1 > 2: 1, _ : 1 + 2,};a`)
	AssertCodesEvalToSameValue(t, `1`, `let a = cond {1 < 2: 1, _ : 1 + 2};a * 1`)
	// // Multiple true conditions
	AssertCodesEvalToSameValue(t, `1`, `cond {1 > 0 : 1, 2 < 3: 2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `2`, `cond {1 > 2 : 1, 2 < 3: 2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `3`, `cond {1 > 2 : 1, 2 > 3: 2, 3 < 4 :3, _:2 + 2}`)
	AssertCodesEvalToSameValue(t, `3`, `cond {1 > 2 : 1, 2 > 3: 2, 3 < 4 :3, 4 < 5 : 5, _:2 + 2}`)
	AssertCodesEvalToSameValue(t, `3`, `cond {1 > 2 : 1, 2 > 3: 2, 3 < 4 :3, 4 < 5 : 5, 5 > 6 : 6, _:2 + 2}`)
	AssertCodesEvalToSameValue(t, `2`, `cond {1 > 2 : 1, 2 < 3: 2, 3 < 4 :3, 4 < 5 : 5, 5 > 6 : 6, _:2 + 2}`)

	// // Nested call
	AssertCodesEvalToSameValue(t, `1`, `cond {cond {1 > 0 : 1} > 0 : 1, 2 < 3: 2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `2`, `cond {cond {1 > 2 : 1, _ : 11} < 2 : 1, 2 < 3: 2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `20`, `let a = cond {cond {1 > 2 : 1, _ : 11} < 2 : 1, 2 < 3: 2, _:1 + 2};a * 10`)

	AssertCodesEvalToSameValue(t, `{}`, `cond {1 < 0 : 1, 2 > 3: 2}`)
	AssertCodesEvalToSameValue(t, `{}`, `cond {1 < 0 : 1}`)
	AssertCodesEvalToSameValue(t, `{}`, `cond {}`)
}

// TestEvalCondMulti executes the cases whose condition has multiple expressions.
func TestEvalCondMulti(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `1`, `cond {1 > 0 || 3 > 2: 1, 2 > 3: 2, _:1 + 2,}`)
	AssertCodesEvalToSameValue(t, `1`, `cond {0 > 1 || 3 > 2: 1, 2 > 3: 2, _:1 + 2,}`)
	AssertCodesEvalToSameValue(t, `3`, `cond {0 > 1 || 3 > 4: 1, 2 > 3: 2, _:1 + 2,}`)
	AssertCodesEvalToSameValue(t, `1`, `cond {1 > 0 && 3 > 2: 1, (2 > 3): 2, _:1 + 2,}`)
	AssertCodesEvalToSameValue(t, `1`, `cond {1 > 0 && 3 > 2: 1, (2 > 3 || 1 < 0): 2, _:1 + 2,}`)
	AssertCodesEvalToSameValue(t, `3`, `let a = cond {1 > 2 && 2 > 1: 1, _ : 1 + 2};a`)
	// Multiple true conditions
	AssertCodesEvalToSameValue(t, `1`, `cond {(1 > 0 && 3 > 2): 1, (2 > 1): 2, _:1 + 2,}`)
	AssertCodesEvalToSameValue(t, `2`, `cond {(1 > 0 && 3 < 2): 1, (2 > 1 || 1 > 0): 2, _:1 + 2,}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = cond {(1 > 2 && 2 > 1): 1, (2 > 1) : 2, _ : 1 + 2};a`)

	AssertCodesEvalToSameValue(t, `{}`, `cond {1 < 0 || 2 > 3 : 1, 2 > 3: 2}`)
	AssertCodesEvalToSameValue(t, `{}`, `cond {1 < 0 || 3 > 4 : 1}`)
}

func TestEvalCondStr(t *testing.T) {
	t.Parallel()
	AssertEvalExprString(t, "{(1>0):1,(2>3):2,_:(1+2)}", "cond {(1 > 0) : 1, (2 > 3): 2, _:1 + 2,}")
	AssertEvalExprString(t, "{(1>0):1,(2>3):2,_:(1+2)}", "cond {(1 > 0) : 1, (2 > 3): 2, _:1 + 2}")
	AssertEvalExprString(t, "{(1<2):1}", "cond {(1 < 2) : 1}")
	AssertEvalExprString(t, "{(1>2):1,(2<3):2}", "cond {(1 > 2) : 1, (2 < 3): 2}")
	AssertEvalExprString(t, "{_:(1+2)}", "cond {_: 1 + 2}")
	AssertEvalExprString(t, "{(1<2):1,_:(1+2)}", "cond {(1 < 2): 1, _ : 1 + 2}")
}

func TestEvalCondWithControlVar(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `let a = 1; cond a {(1 + 2) :1, 2 :2}`)
	// // Control var conditions
	AssertCodesEvalToSameValue(t, `1`, `let a = 1; cond a {1 :1}`)
	AssertCodesEvalToSameValue(t, `1`, `let a = 1; cond a {1 :1, 2:2}`)
	AssertCodesEvalToSameValue(t, `1`, `let a = 1; cond a {1 :1, 2 :2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `11`, `let a = 1; cond a {(1) :1 + 10, (2) : 2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `3`, `let a = 1; cond a {234 :2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `3`, `let a = 1; cond a {(234) :2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = 234; cond a {234 :2, _:11 + 2}`)
	AssertCodesEvalToSameValue(t, `5`, `let a = 3; cond a {_:5}`)
	AssertCodesEvalToSameValue(t, `3`, `let a = 3; cond a {_:1 + 2}`)

	AssertCodesEvalToSameValue(t, `1`, `let a = 1; let b = cond a {1 :1, 2 :2, _:1 + 2}; b`)
	AssertCodesEvalToSameValue(t, `100`, `let a = 1; let b = cond a {1 :1, 2 :2, _:1 + 2}; b * 100`)
	// // //
	AssertCodesEvalToSameValue(t, `2`, `let a = 1; cond (a + 100) {1 :1, 101 :2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `3`, `let a = 1; cond a + 10 {1 :1, 2 :2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = 1; let b = cond a + 1 {1 :1, 2 :2, _:1 + 2}; b`)
	AssertCodesEvalToSameValue(t, `300`, `let a = 1; let b = cond a + 10 {1 :1, 2 :2, _:1 + 2}; b * 100`)
	// Nested call
	AssertCodesEvalToSameValue(t, `"B"`, `let a = 2; cond a {(cond a {(1,2) : 1}): "A", (2, 3): "B", _: "C"}`)
	AssertCodesEvalToSameValue(t, `"A"`, `let a = 1; cond a {(cond {(2 > 1) : 1}): "A", (2, 3): "B", _: "C"}`)
	AssertCodesEvalToSameValue(t, `"A"`, `let a = 1; cond {(cond a {(1) : 1}) : "A", (2): "B", _: "C"}`)

	AssertCodesEvalToSameValue(t, `{}`, `let a = 3; cond a {1 :1, 2 :2 + 1}`)
	AssertCodesEvalToSameValue(t, `{}`, `let a = 3; let b = cond a {1 :1, 2 :2 + 1}; b`)
	AssertCodesEvalToSameValue(t, `{}`, `let a = 3; let b = cond a + 10 {1 :1, 2 :2 + 1}; b`)

	AssertCodesEvalToSameValue(t, `1`, `let x = 1; cond 2 { 2: x }`)
}

func TestEvalCondWithControlVarMulti(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `1`, `let a = 1; cond a {(1 + 0,2 + 0) :1}`)
	AssertCodesEvalToSameValue(t, `1`, `let a = 2; cond a {(1 + 0,2 + 0,3 + 0) :1, (2 + 0):2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = 2; cond a {(1 + 0,2 + 4,3 + 5) :1, (2 + 0):2}`)
	AssertCodesEvalToSameValue(t, `11`, `let [a,b] = [2,4]; cond a {(1 + 0,b -2,3 + 5) :11, (2):2}`)
	AssertCodesEvalToSameValue(t, `1`, `let a = 3; cond a {(1 + 0,2 + 0,3 + 0) :1, (2 + 0) :2, _:1 + 2}`)

	AssertCodesEvalToSameValue(t, `2`, `let a = 2; cond a {(1) :1 + 10, (2,3) : 2, _:1 + 2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = 2; cond a {(1) :1 + 10, (3,2) : 2, _:1 + 2}`)

	AssertCodesEvalToSameValue(t, `"med"`, `let a = 2;
	cond a {
		(1):"lo",
		(2,3): "med",
		_: "hi"}`)
}

// TestEvalCondMultiStr executes the cases whose condition has multiple expressions.
func TestEvalCondMultiStr(t *testing.T) {
	t.Parallel()

	AssertEvalExprString(t, "((control_var:1),{((1>0))&&((2>1)):1})", "cond (1) {(1 > 0 && 2 > 1) : 1}")
	AssertEvalExprString(t, "((control_var:1),{((1>0))||((2>1)):1})", "cond (1) {(1 > 0 || 2 > 1) : 1}")
	AssertEvalExprString(t, "((control_var:1),{((1>0))||((2>1)):1,_:11})", "cond (1) {(1 > 0 || 2 > 1) : 1, _ : 11}")
}

func TestEvalCondWithControlVarStr(t *testing.T) {
	t.Parallel()

	AssertEvalExprString(t, "((control_var:1),{1:1})", "cond (1) {1 : 1}")
	AssertEvalExprString(t, "((control_var:1),{1:1})", "cond (1) {1 : 1,}")
	AssertEvalExprString(t, "((control_var:1),{1:1,(2+1):3})", "cond (1) {1 : 1, (2 + 1) : 3}")
	AssertEvalExprString(t, "((control_var:1),{1:1,(2+1):3})", "cond (1) {(1) : 1, (2 + 1) : 3,}")
	AssertEvalExprString(t, "((control_var:1),{1:1,(2+1):3,_:4})", "cond (1) {1 : 1, (2 + 1) : 3, _ : 4}")
	AssertEvalExprString(t, "((control_var:1),{1:1,(2+1):3,_:4})", "cond (1) {(1) : 1, (2 + 1) : 3, _ : 4,}")

	AssertEvalExprString(t, "(1->(\\a((control_var:a),{1:1})))",
		"let a = 1; cond a {(1) : 1}")
	AssertEvalExprString(t, "(1->(\\a((control_var:a),{(1+2):1,_:(1+2)})))",
		"let a = 1; cond a {(1 + 2): 1, _ : 1 + 2}")
	AssertEvalExprString(t, "(2->(\\a(((control_var:a),{(1+2):1,_:(1+2)})->(\\b(b*1)))))",
		"let a = 2; let b = cond a {(1 + 2): 1, _ : 1 + 2}; b * 1")
	AssertEvalExprString(t, "(3->(\\a((control_var:(a+2)),{(1+2):1,_:(1+2)})))",
		"let a = 3; cond (a + 2) {(1 + 2): 1, _ : 1 + 2}")
}

func TestEvalCondWithControlVarMultiStr(t *testing.T) {
	t.Parallel()
	AssertEvalExprString(t, "((control_var:1),{[1,2]:1})", "cond (1) {(1,2) :1}")
	AssertEvalExprString(t, "((control_var:2),{1:(1+10),[2,3]:2,_:(1+2)})", "cond (2) {(1) :1 + 10, (2,3) : 2, _:1 + 2}")
}

func TestEvalCondPatternMatchingWithControlVar(t *testing.T) { //nolint:dupl
	t.Parallel()

	AssertCodesEvalToSameValue(t, `2`, `let a = 'A' ; cond a {'A':2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = 'ABC' ; cond a {'ABC':2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = "A" ; cond a {"A":2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = "ABC" ; cond a {"ABC":2}`)
	AssertCodesEvalToSameValue(t, `2`, `let [a,b,c] = [10,100,1000]; cond 1100 {(b + c):2}`)

	AssertCodesEvalToSameValue(t, `2`, `let a = [1, 2]; cond a {[1, 2]: 2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = ['a', 'b']; cond a {['a', 'b']: 2}`)
	AssertCodesEvalToSameValue(t, `6`, `let a = (x:4); cond a {(x:x): x + 2}`)
	AssertCodesEvalToSameValue(t, `6`, `cond (x:4) {(x:x): x + 2}`)

	AssertCodesEvalToSameValue(t, `8`, `let a = (a:3); cond a {(a:x): x + 5,_:2}`)
	AssertCodesEvalToSameValue(t, `8`, `let a = {"a":3}; cond a {{"a":x}: x + 5,_:2}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = {"a":3}; cond a {1 : x + 5,_:2}`)

	AssertCodesEvalToSameValue(t, `10`, `cond ({'a': {'b': 10}}) {{'a': {'b': b}}: b, _:42}`)
	AssertCodesEvalToSameValue(t, `42`, `cond ({'a': {'z': 10}}) {{'a': {'b': b}}: b, _:42}`)

	AssertCodesEvalToSameValue(t, `10`, `cond (a: (b: 10)) {(a: (:b)): b, _:42}`)
	AssertCodesEvalToSameValue(t, `42`, `cond (a: (z: 10)) {(a: (:b)): b, _:42}`)

	AssertCodesEvalToSameValue(t, `42`, `cond [] {[b]: b, _:42}`)
	AssertCodesEvalToSameValue(t, `42`, `cond (a: []) {(a: [b]): b, _:42}`)
	AssertCodesEvalToSameValue(t, `43`, `cond () {{1:2}: 42, _: 43}`)

	AssertCodesEvalToSameValue(t, `{}`, `let a = 2; cond a {[1,2,3]: 6}`)
	AssertCodesEvalToSameValue(t, `2`, `let a = {"a":3}; cond a {(a:x): x + 5,_:2}`)
}

func TestCondEvalWithWhere(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, "{111,222,333}",
		"{1, 2, 3} where cond . {1:11 , 2:22, 3:33} => cond . {1:111, 2:222, 3:333}")
	AssertCodesEvalToSameValue(t, "{111,222,{}}",
		"{1, 2, 3} where cond . {1:11 , 2:22, 3:33} => cond . {1:111, 2:222}")

	AssertCodesEvalToSameValue(t, "{1,5}",
		"{1, [2, 3]} where cond . {k: true, [a, b]: true} => cond . {[a, b]: a + b, k: k}")
	AssertCodesEvalToSameValue(t, "{1,5,4,11}",
		"{1, [2, 3], 4, [5, 6]} where cond . {k: true, [a, b]: true} => cond . {[a, b]: a + b, k:k}")

	AssertCodesEvalToSameValue(t, "{9,5}",
		"{[2,3,4], [2, 3]} where cond . {[a,b,c]: true, [a, b]: true} => cond . {[a,b,c]:a+b+c, [a, b]: a + b}")
	AssertCodesEvalToSameValue(t, "{9,5}",
		"{[2,3,4], [2, 3]} where cond . {[a,b,c]: true, [a, b]: true} => cond . {[a, b]: a + b , [a,b,c]:a+b+c}")

	AssertCodesEvalToSameValue(t, "{1, [2,3]}",
		"{1, [2, 3]} where cond . {k: true, [a, b]: true} => cond . { k: k , [a, b]: a + b}")
}

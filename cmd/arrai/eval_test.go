package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEvalOutputs(t *testing.T, expected, source string) bool { //nolint:unparam
	var sb strings.Builder
	return assert.NoError(t, evalImpl(source, &sb)) &&
		assert.Equal(t, expected, strings.TrimRight(sb.String(), "\n"))
}

func TestEvalNumberULP(t *testing.T) {
	assertEvalOutputs(t, `0.3`, `0.1 + 0.1 + 0.1`)
}

func TestEvalString(t *testing.T) {
	assertEvalOutputs(t, ``, `""`)
	assertEvalOutputs(t, ``, `{}`)
	assertEvalOutputs(t, `abc`, `"abc"`)
}

func TestEvalComplex(t *testing.T) {
	assertEvalOutputs(t, `[42, 'abc']`, `[42, "abc"]`)
	assertEvalOutputs(t, `{42, 'abc'}`, `{"abc", 42}`)
}
<<<<<<< HEAD
=======

func TestEvalCond(t *testing.T) {
	t.Parallel()
	assertEvalOutputs(t, `1`, `cond (1 > 0 : 1, 2 > 3: 2, *:1 + 2,)`)
	assertEvalOutputs(t, `1`, `cond (1 > 0 : 1, 2 > 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `3`, `cond (1 > 2 : 1, 2 > 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `1`, `cond (1 < 2 : 1,)`)
	assertEvalOutputs(t, `1`, `cond (1 < 2 : 1)`)
	assertEvalOutputs(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2)`)
	assertEvalOutputs(t, `3`, `cond (* : 1 + 2)`)
	assertEvalOutputs(t, `1`, `cond (1 < 2: 1, * : 1 + 2)`)
	assertEvalOutputs(t, `3`, `cond (1 > 2: 1, * : 1 + 2)`)
	assertEvalOutputs(t, `3`, `let a = cond (1 > 2: 1, * : 1 + 2);a`)
	assertEvalOutputs(t, `3`, `let a = cond (1 > 2: 1, * : 1 + 2,);a`)
	assertEvalOutputs(t, `1`, `let a = cond (1 < 2: 1, * : 1 + 2);a * 1`)
	// Multiple true conditions
	assertEvalOutputs(t, `1`, `cond (1 > 0 : 1, 2 < 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `3`, `cond (1 > 2 : 1, 2 > 3: 2, 3 < 4 :3, *:2 + 2)`)
	assertEvalOutputs(t, `3`, `cond (1 > 2 : 1, 2 > 3: 2, 3 < 4 :3, 4 < 5 : 5, *:2 + 2)`)
	assertEvalOutputs(t, `3`, `cond (1 > 2 : 1, 2 > 3: 2, 3 < 4 :3, 4 < 5 : 5, 5 > 6 : 6, *:2 + 2)`)
	assertEvalOutputs(t, `2`, `cond (1 > 2 : 1, 2 < 3: 2, 3 < 4 :3, 4 < 5 : 5, 5 > 6 : 6, *:2 + 2)`)

	// Nested call
	assertEvalOutputs(t, `1`, `cond (cond (1 > 0 : 1) > 0 : 1, 2 < 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `2`, `cond (cond (1 > 2 : 1, * : 11) < 2 : 1, 2 < 3: 2, *:1 + 2)`)
	assertEvalOutputs(t, `20`, `let a = cond (cond (1 > 2 : 1, * : 11) < 2 : 1, 2 < 3: 2, *:1 + 2);a * 10`)

	assertEvalOutputs(t, ``, `cond (1 < 0 : 1, 2 > 3: 2)`)
	assertEvalOutputs(t, ``, `cond (1 < 0 : 1)`)
	assertEvalOutputs(t, ``, `cond ()`)
}

func TestEvalCondStr(t *testing.T) {
	t.Parallel()
	assertEvalExprString(t, "((1>0):1,(2>3):2,*:(1+2))", "cond (1 > 0 : 1, 2 > 3: 2, *:1 + 2,)")
	assertEvalExprString(t, "((1>0):1,(2>3):2,*:(1+2))", "cond (1 > 0 : 1, 2 > 3: 2, *:1 + 2)")
	assertEvalExprString(t, "((1<2):1)", "cond (1 < 2 : 1)")
	assertEvalExprString(t, "((1>2):1,(2<3):2)", "cond (1 > 2 : 1, 2 < 3: 2)")
	assertEvalExprString(t, "(*:(1+2))", "cond (*: 1 + 2)")
	assertEvalExprString(t, "((1<2):1,*:(1+2))", "cond (1 < 2: 1, * : 1 + 2)")
}

// TestEvalCondMulti executes the cases whose condition has multiple expressions.
func TestEvalCondMulti(t *testing.T) {
	t.Parallel()
	assertEvalOutputs(t, `1`, `cond (1 > 0 || 3 > 2: 1, 2 > 3: 2, *:1 + 2,)`)
	assertEvalOutputs(t, `1`, `cond (0 > 1 || 3 > 2: 1, 2 > 3: 2, *:1 + 2,)`)
	assertEvalOutputs(t, `3`, `cond (0 > 1 || 3 > 4: 1, 2 > 3: 2, *:1 + 2,)`)
	assertEvalOutputs(t, `1`, `cond (1 > 0 && 3 > 2: 1, 2 > 3: 2, *:1 + 2,)`)
	assertEvalOutputs(t, `1`, `cond ((1 > 0 && 3 > 2): 1, (2 > 3) || (1 < 0): 2, *:1 + 2,)`)
	assertEvalOutputs(t, `3`, `let a = cond (1 > 2 && 2 > 1: 1, * : 1 + 2);a`)
	// Multiple true conditions
	assertEvalOutputs(t, `1`, `cond (1 > 0 && 3 > 2: 1, 2 > 1: 2, *:1 + 2,)`)
	assertEvalOutputs(t, `2`, `cond ((1 > 0 && 3 < 2): 1, (2 > 1) || (1 > 0): 2, *:1 + 2,)`)
	assertEvalOutputs(t, `2`, `let a = cond (1 > 2 && 2 > 1: 1, (2 > 1) : 2, * : 1 + 2);a`)

	assertEvalOutputs(t, ``, `cond (1 < 0 || 2 > 3 : 1, 2 > 3: 2)`)
	assertEvalOutputs(t, ``, `cond (1 < 0 || 3 > 4 : 1)`)
}

// TestEvalCondMultiStr executes the cases whose condition has multiple expressions.
func TestEvalCondMultiStr(t *testing.T) {
	t.Parallel()
	assertEvalExprString(t, "((control_var:1),(((1>0))&&((2>1)):1))", "(1) cond (1 > 0 && 2 > 1 : 1)")
	assertEvalExprString(t, "((control_var:1),(((1>0))||((2>1)):1))", "(1) cond (1 > 0 || 2 > 1 : 1)")
	assertEvalExprString(t, "((control_var:1),(((1>0))||((2>1)):1,*:11))", "(1) cond (1 > 0 || 2 > 1 : 1, * : 11)")
}

//nolint:dupl
func TestEvalCondWithControlVar(t *testing.T) {
	t.Parallel()
	// Control var conditions
	assertEvalOutputs(t, `1`, `let a = 1; a cond (1 :1)`)
	assertEvalOutputs(t, `1`, `let a = 1; a cond (1 :1, 2 :2)`)
	assertEvalOutputs(t, `1`, `let a = 1; a cond (1 :1, 2 :2, *:1 + 2)`)
	assertEvalOutputs(t, `11`, `let a = 1; a cond (1 :1 + 10, 2 : 2, *:1 + 2)`)
	assertEvalOutputs(t, `3`, `let a = 1; a cond (2 :2, *:1 + 2)`)
	assertEvalOutputs(t, `13`, `let a = 1; a cond (2 :2, *:11 + 2)`)
	assertEvalOutputs(t, `5`, `let a = 3; a cond (*:5)`)
	assertEvalOutputs(t, `3`, `let a = 3; a cond (*:1 + 2)`)

	assertEvalOutputs(t, `1`, `let a = 1; let b = a cond (1 :1, 2 :2, *:1 + 2); b`)
	assertEvalOutputs(t, `100`, `let a = 1; let b = a cond (1 :1, 2 :2, *:1 + 2); b * 100`)
	//
	assertEvalOutputs(t, `2`, `let a = 1; (a + 1) cond (1 :1, 2 :2, *:1 + 2)`)
	assertEvalOutputs(t, `3`, `let a = 1; (a + 10) cond (1 :1, 2 :2, *:1 + 2)`)
	assertEvalOutputs(t, `2`, `let a = 1; let b = (a + 1) cond (1 :1, 2 :2, *:1 + 2); b`)
	assertEvalOutputs(t, `300`, `let a = 1; let b = (a + 10) cond (1 :1, 2 :2, *:1 + 2); b * 100`)
	// Nested call
	assertEvalOutputs(t, "B", `let a = 2; a cond ( a cond ((1,2) : 1): "A", (2, 3): "B", *: "C")`)
	assertEvalOutputs(t, "A", `let a = 1; a cond ( cond (2 > 1 : 1): "A", (2, 3): "B", *: "C")`)
	assertEvalOutputs(t, "A", `let a = 1; cond ( a cond (1 : 1) : "A", 2: "B", *: "C")`)

	assertEvalOutputs(t, ``, `let a = 3; a cond (1 :1, 2 :2 + 1)`)
	assertEvalOutputs(t, ``, `let a = 3; let b = a cond (1 :1, 2 :2 + 1); b`)
	assertEvalOutputs(t, ``, `let a = 3; let b = (a + 10) cond (1 :1, 2 :2 + 1); b`)
}

func TestEvalCondWithControlVarStr(t *testing.T) {
	t.Parallel()
	assertEvalExprString(t, "((control_var:1),(1:1))", "(1) cond (1 : 1)")
	assertEvalExprString(t, "((control_var:1),(1:1))", "(1) cond (1 : 1,)")
	assertEvalExprString(t, "((control_var:1),(1:1,(2+1):3))", "(1) cond (1 : 1, 2 + 1 : 3)")
	assertEvalExprString(t, "((control_var:1),(1:1,(2+1):3))", "(1) cond (1 : 1, 2 + 1 : 3,)")
	assertEvalExprString(t, "((control_var:1),(1:1,(2+1):3,*:4))", "(1) cond (1 : 1, 2 + 1 : 3, * : 4)")
	assertEvalExprString(t, "((control_var:1),(1:1,(2+1):3,*:4))", "(1) cond (1 : 1, 2 + 1 : 3, * : 4,)")

	assertEvalExprString(t, "(1->(\\a((control_var:a),(1:1))))",
		"let a = 1; a cond (1 : 1)")
	assertEvalExprString(t, "(1->(\\a((control_var:a),((1+2):1,*:(1+2)))))",
		"let a = 1; a cond (1 + 2: 1, * : 1 + 2)")
	assertEvalExprString(t, "(2->(\\a(((control_var:a),((1+2):1,*:(1+2)))->(\\b(b*1)))))",
		"let a = 2; let b = a cond (1 + 2: 1, * : 1 + 2); b * 1")
	assertEvalExprString(t, "(3->(\\a((control_var:(a+2)),((1+2):1,*:(1+2)))))",
		"let a = 3; (a + 2) cond (1 + 2: 1, * : 1 + 2)")
}

func TestEvalCondWithControlVarMulti(t *testing.T) {
	assertEvalOutputs(t, `1`, `let a = 1; a cond ((1,2) :1)`)
	assertEvalOutputs(t, `1`, `let a = 2; a cond ((1,2,3) :1, 2 :2)`)
	assertEvalOutputs(t, `1`, `let a = 3; a cond ((1,2,3) :1, 2 :2, *:1 + 2)`)
	assertEvalOutputs(t, `2`, `let a = 2; a cond (1 :1 + 10, (2,3) : 2, *:1 + 2)`)

	assertEvalOutputs(t, `med`, `let a = 2;
	a cond (
		1:"lo",
		(2,3): "med",
		*: "hi")`)

	var sb strings.Builder
	assert.Error(t, evalImpl(`let a = 1; a cond ((2,3)) : 2, 3: 3)`, &sb))
	assert.Error(t, evalImpl(`let a = 1; a cond ((2,3)) : 2, (3,5): 3)`, &sb))
}

func TestEvalCondWithControlVarMultiStr(t *testing.T) {
	t.Parallel()
	assertEvalExprString(t, "((control_var:1),([1,2]:1))", "(1) cond ((1,2) :1)")
	assertEvalExprString(t, "((control_var:2),(1:(1+10),[2,3]:2,*:(1+2)))", "(2) cond (1 :1 + 10, (2,3) : 2, *:1 + 2)")
}
>>>>>>> master

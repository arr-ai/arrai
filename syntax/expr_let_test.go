package syntax

import "testing"

func TestExprLet(t *testing.T) { //nolint:dupl
	t.Parallel()
	AssertCodesEvalToSameValue(t, `7`, `let x = 6; 7`)
	AssertCodesEvalToSameValue(t, `42`, `let x = 6; x * 7`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let x = 1; [x, 2]`)
	AssertCodesEvalToSameValue(t, `2`, `let x = 1; let x = x + 1; x`)
	AssertCodesEvalToSameValue(t, `(x: 1)`, `let x = 1; (:x)`)
	AssertCodesEvalToSameValue(t, `(x: 1, y: 2)`, `let x = 1; let y = 2; (:x, :y)`)
	AssertCodesEvalToSameValue(t, `(x: 1, y: 2)`, `let x = 1; (:x, y: 2)`)

	AssertCodesEvalToSameValue(t, `4`, `let x = 4; let (x) = 4;x`)
	AssertCodesEvalToSameValue(t, `4`, `let x = 4; let (x) = (4);x`)
	AssertCodesEvalToSameValue(t, `4`, `let x = 4; let (x) =(4 + 0);x`)
	AssertCodesEvalToSameValue(t, `1`, `let [a,b] = [1,2];let (b) = (a + 1);a`)

	AssertCodesEvalToSameValue(t, `1`, `let a = 1;a`)
	AssertCodesEvalToSameValue(t, `1`, `let a = 1;(a)`) // (a) is an expression

	// (x) should be parsed as an expression and fail because x isn't bound.
	AssertCodeErrors(t, "", `let (x) = 5;x`)
	AssertCodeErrors(t, "", `let (x) = 5;(x)`)
}

func TestExprLetExprPattern(t *testing.T) { //nolint:dupl
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `let 42 = 42; 42`)
	AssertCodesEvalToSameValue(t, `42`, `let (42) = 42; 42`)
	AssertCodesEvalToSameValue(t, `1`, `let 42 = 42; 1`)
	AssertCodesEvalToSameValue(t, `1`, `let "hello" = "hello"; 1`)
	AssertCodesEvalToSameValue(t, `5`, `let 3 = 1 + 2; 5`)
	AssertCodesEvalToSameValue(t, `5`, `let (1 + 2) = 3; 5`)
	AssertCodesEvalToSameValue(t, `3`, `let a = 1 + 2; a`)
	AssertCodesEvalToSameValue(t, `3`, `let a = 3; let a = 1 + 2; a`)
	AssertCodesEvalToSameValue(t, `3`, `let a = 3; let (x: (a)) = (x: 3); a`)
	AssertCodesEvalToSameValue(t, `3`, `let true = true; 3`)
	AssertCodesEvalToSameValue(t, `3`, `let false = false; 3`)
	AssertCodesEvalToSameValue(t, `3`, `let true = {()}; 3`)
	AssertCodesEvalToSameValue(t, `3`, `let false = {}; 3`)
	AssertCodesEvalToSameValue(t, `3`, `let true = {()}; 3`)
	AssertCodesEvalToSameValue(t, `2`, `let ('':value, ...) = ('':2); value`)
	AssertCodesEvalToSameValue(t, `2`, `let (..., '':value) = ('':2); value`)

	AssertCodeErrors(t, "", `let 42 = 1; 42`)
	AssertCodeErrors(t, "", `let 42 = 1; 1`)
	AssertCodeErrors(t, "", `let "hello" = "hi"; 1`)
	AssertCodeErrors(t, "", `let 1 = 1 + 2; 5`)
	AssertCodeErrors(t, "", `let (1 + 2) = 6; 5`)
	AssertCodeErrors(t, "", `let a = 5; let (a) = 1 + 2; a`)
	AssertCodeErrors(t, "", `let a = 1; let (x: (a)) = (x: 3); a`)
	AssertCodeErrors(t, "", `let true = false; 3`)
	AssertCodeErrors(t, "", `let true = {}; 3`)
}

func TestExprLetIdentPattern(t *testing.T) {
	AssertCodesEvalToSameValue(t, `3`, `let f = \[x, y] x + y; f([1, 2])`)
	AssertCodesEvalToSameValue(t, `1`, `let m = {"a": 1}("a")?:42; m`)
	AssertCodesEvalToSameValue(t, `42`, `let m = {"a": 1}("b")?:42; m`)
	AssertCodesEvalToSameValue(t, `0`, `let arr = [1, 2]; let z = arr(2)?:0; z`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let ids = {'ids': [1, 2]}('ids')?:[]; ids`)
	AssertCodesEvalToSameValue(t, `[]`, `let ids = {'ids': [1, 2]}('id')?:[]; ids`)
}

func TestExprLetArrayPattern(t *testing.T) { //nolint:dupl
	t.Parallel()
	AssertCodesEvalToSameValue(t, `1`, `let [] = []; 1`)
	AssertCodesEvalToSameValue(t, `9`, `let [a, b, c] = [1, 2, 3]; 9`)
	//TODO: implement pattern matching for sparse array
	// AssertCodesEvalToSameValue(t, `9`, `let [a, b, , c] = [1, 2, , 3]; 9`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `let [a, b, c] = [1, 2, 3]; [a, b, c]`)
	AssertCodesEvalToSameValue(t, `2`, `let [a, b, c] = [1, 2, 3]; b`)
	AssertCodesEvalToSameValue(t, `2`, `let arr = [1, 2]; let [a, b] = arr; b`)
	AssertCodesEvalToSameValue(t, `1`, `let [x, x] = [1, 1]; x`)
	AssertCodesEvalToSameValue(t, `1`, `let [x, _, _] = [1, 2, 3]; x`)
	AssertCodesEvalToSameValue(t, `2`, `let [_, x, _] = [1, 2, 3]; x`)
	AssertCodesEvalToSameValue(t, `3`, `let x = 3; let [(x)] = [3]; x`)
	AssertCodesEvalToSameValue(t, `2`, `let x = 3; let [b, (x)] = [2, 3]; b`)
	AssertCodesEvalToSameValue(t, `2`, `let x = 3; let [_, b, (x)] = [1, 2, 3]; b`)
	AssertCodesEvalToSameValue(t, `2`, `let x = 3; let [x] = [2]; x`)
	AssertCodesEvalToSameValue(t,
		`('': [75\'+'], @rule: 'expr', expr: [(expr: [('': 74\'1')]), ('': [77\'*'], expr: [('': 76\'2'), ('': 78\'3')])])`,
		`let [g] = [{://grammar.lang.wbnf: expr -> @:[-+] > @:[/*] > \d+; :}]; {:g:1+2*3:}`)
	AssertCodesEvalToSameValue(t,
		`('': [94\'+'], @rule: 'expr', expr: [(expr: [('': 93\'1')]), ('': [96\'*'], expr: [('': 95\'2'), ('': 97\'3')])])`,
		`let (a: g, b: x) = (a: {://grammar.lang.wbnf: expr -> @:[-+] > @:[/*] > \d+; :}, b: 42); {:g:1+2*3:}`)

	AssertCodeErrors(t, "", `let [(x)] = [2]; x`)
	AssertCodeErrors(t, "", `let x = 3; let [(x)] = [2]; x`)
	AssertCodeErrors(t, "", `let [x, y] = 1; x`)
	AssertCodeErrors(t, "", `let [x, x] = [1]; x`)
	AssertCodeErrors(t, "", `let [x, y] = [1]; x`)
	AssertCodeErrors(t, "", `let [x, x] = [1, 2]; x`)
	AssertCodeErrors(t, "name \"_\" not found in {}\n\n\x1b[1;37m:1:16:\x1b[0m\nlet [_]", `let [_] = [1]; _`)
	AssertCodeErrors(t, "", `let x = 3; let [b, (x)] = [2, 1]; b`)
}

func TestExprLetTuplePattern(t *testing.T) { //nolint:dupl
	t.Parallel()
	AssertCodesEvalToSameValue(t, `4`, `let () = (); 4`)
	AssertCodesEvalToSameValue(t, `4`, `let (a: x, b: y) = (a: 4, b: 7); x`)
	AssertCodesEvalToSameValue(t, `4`, `let (a: x, b: x) = (a: 4, b: 4); x`)
	AssertCodesEvalToSameValue(t, `4`, `let x = 4; let (a: x) = (a: 4); x`)
	AssertCodesEvalToSameValue(t, `4`, `let x = 5; let (a: x) = (a: 4); x`)
	AssertCodesEvalToSameValue(t, `4`, `let (a: [x]) = (a: [4]); x`)
	AssertCodesEvalToSameValue(t, `1`, `let (:x) = (x: 1); x`)
	AssertCodesEvalToSameValue(t, `2`, `let (:x, :y) = (x: 1, y: 2); y`)
	AssertCodesEvalToSameValue(t, `2`, `let a = 1; let (:a, b: (a)) = (a: 2, b: 1); a`)
	AssertCodesEvalToSameValue(t, `2`, `let (?:a:1) = (a: 2); a`)
	AssertCodesEvalToSameValue(t, `2`, `let (a?: x:1) = (a: 2); x`)
	AssertCodesEvalToSameValue(t, `1`, `let (?:a:1, ?:b:2) = (b: 2); a`)
	AssertCodesEvalToSameValue(t, `(b: 2, c: 3)`, `let (a: 1, ...tail) = (a: 1, b: 2, c: 3); tail`)
	AssertCodesEvalToSameValue(t, `(b: 2)`, `let (a: 1, ...tail, c: 3) = (a: 1, b: 2, c: 3); tail`)
	AssertCodesEvalToSameValue(t, `(b: 2, d: 4, e: 5)`, `let (a: 1, ...tail, c: 3) = (a: 1, b: 2, c: 3, d: 4, e: 5); tail`)
	AssertCodesEvalToSameValue(t, `(b: 2, d: 4)`, `let (a: 1, ...tail, c: 3, e: 5) = (a: 1, b: 2, c: 3, d: 4, e: 5); tail`)
	AssertCodesEvalToSameValue(t, `()`, `let (a: 1, b: 2, c: 3, ...tail) = (a: 1, b: 2, c: 3); tail`)
	AssertCodesEvalToSameValue(t, `(b: 2)`, `let (...tail, a:1, c: 3) = (a: 1, b: 2, c: 3); tail`)
	AssertCodesEvalToSameValue(t, `1`, `let (?:a:1, b: 2) = (b: 2); a`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `let (?:a:1, ?:b:2, ?:c:3, d: 4) = (d: 4); [a, b, c]`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `let (?:a:1, ?:b:2, ?:c:3) = (); [a, b, c]`)
	AssertCodesEvalToSameValue(t, `1`, `let (?:a:1, ...) = (b: 2); a`)
	AssertCodesEvalToSameValue(t, `[(b: 2), 1]`, `let (?:a:1, ...tail) = (b: 2); [tail, a]`)
	AssertCodesEvalToSameValue(t, `[(), 1]`, `let (?:a:1, b: 2, ...tail) = (b: 2); [tail, a]`)
	AssertCodesEvalToSameValue(t, `[(), 3]`, `let (?:a:1, b: 2, ...tail) = (a: 3, b: 2); [tail, a]`)
	AssertCodesEvalToSameValue(t, `[(), 3, 5]`, `let (?:a:1, b: 2, ?:c: 5, ...tail) = (a: 3, b: 2); [tail, a, c]`)
	AssertCodeErrors(t,
		"length of tuple (a: 4, b: 7) longer than tuple pattern (a: x)",
		`let (a: x) = (b: 7, a: 4); x`,
	)
	AssertCodeErrors(t,
		"length of tuple (b: 7) longer than tuple pattern (a?: a:1)",
		`let (?:a:1) = (b: 7); a`,
	)
	AssertCodeErrors(t,
		"length of tuple (d: 4) longer than tuple pattern (a?: a:1, b?: b:2, c?: c:3)",
		`let (?:a:1, ?:b:2, ?:c:3) = (d: 4); a`,
	)
	AssertCodeErrors(t,
		"length of tuple (a: 4, b: 7) longer than tuple pattern (a: x, c?: c:2)",
		`let (a: x, ?:c:2) = (b: 7, a: 4); x`,
	)
	AssertCodeErrors(t,
		"length of tuple (b: 7, c: 8) longer than tuple pattern (a?: a:1, b: 7)",
		`let (?:a:1, b: 7) = (b: 7, c: 8); a`,
	)
	AssertCodeParseErrors(t,
		"duplicate fields found in pattern (a: x, a: x)",
		`let (a: x, a: x) = (a: 4, a: 4); x`,
	)
	AssertCodeParseErrors(t,
		"duplicate fields found in pattern (a: x, a: x)",
		`let (a: x, a: x) = (a: 4); x`,
	)
	AssertCodeErrors(t,
		"non-deterministic pattern is not supported",
		`let (a: 1, ..., c: 3, ..., e: 5) = (a: 1, b: 2, c: 3, d: 4, e: 5); a`,
	)
	AssertCodeErrors(t,
		"didn't find matched value\n\n\x1b[1;37m:1:12:\x1b[",
		`let x = 5; let (a: (x)) = (a: 4); x`,
	)
	AssertCodeErrors(t,
		"the value of x is different in both scopes",
		`let (a: x, b: x) = (a: 4, b: 7); x`,
	)
	AssertCodeErrors(t,
		"didn't find matched value\n\n\x1b[1;37m:1:12:\x1b[",
		`let x = 5; let (a: [(x)]) = (a: [4]); x`,
	)
	AssertCodeParseErrors(t,
		"duplicate fields found in pattern (x: a, x: b, ...)",
		`let (x: a, x: b, ...) = (x: 1, y: 2, z: 3); (:a, :b)`,
	)
}

func TestExprLetDictPattern(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `let {1: a} = {1: 42}; a`)
	AssertCodesEvalToSameValue(t, `42`, `let {[1, 2, 3]: a} = {[1, 2, 3]: 42}; a`)
	AssertCodesEvalToSameValue(t, `42`, `let a = 4; let {"x": a} = {"x": 42}; a`)
	AssertCodesEvalToSameValue(t, `42`, `let a = 42; let {"x": (a)} = {"x": 42}; a`)
	AssertCodesEvalToSameValue(t, `[4, 5]`, `let d = {"x": 4, "y": 5}; let {"x": a, "y": b} = d; [a, b]`)
	AssertCodesEvalToSameValue(t, `[4, 5]`, `let {"x": a, "y": b} = {"x": 4, "y": 5}; [a, b]`)
	AssertCodesEvalToSameValue(t, `4`, `let a = 4; let {"x": (a)} = {"x": 4}; a`)
	AssertCodesEvalToSameValue(t, `[4, 5]`, `let a = 4; let {"x": (a), "y": b} = {"x": 4, "y": 5}; [a, b]`)
	AssertCodeErrors(t, "", `let {"x": a, "y": b} = {"x": 4}; a`)
	AssertCodeErrors(t, "", `let {"x": a, "y": b} = {"x": 4, "y": 5, "z": 6}; a`)
	AssertCodeErrors(t, "", `let {"x": a, "x": a} = {"x": 4}; a`)
	AssertCodeErrors(t, "", `let a = 4; let {"x": (a)} = {"x": 5}; a`)
	AssertCodeParseErrors(t, "", `let {"x": a, "x": a} = {"x": 4, "x": 4}; a`)
}

func TestExprLetSetPattern(t *testing.T) { //nolint:dupl
	t.Parallel()
	AssertCodesEvalToSameValue(t, `1`, `let {} = {}; 1`)
	AssertCodesEvalToSameValue(t, `1`, `let {42} = {42}; 1`)
	AssertCodesEvalToSameValue(t, `1`, `let {a} = {1}; a`)
	AssertCodesEvalToSameValue(t, `1`, `let {a, 42} = {42, 1}; a`)
	AssertCodesEvalToSameValue(t, `{1, 42}`, `let {...t} = {1, 42}; t`)
	AssertCodesEvalToSameValue(t, `{42, 43}`, `let {1, 2, 3, ...t} = {1, 2, 3, 42, 43}; t`)
	AssertCodesEvalToSameValue(t, `5`, `let x = 1; let y = 42; let {(x), (y)} = {42, 1}; 5`)
	AssertCodesEvalToSameValue(t, `{5, 6}`, `let x = 1; let y = 42; let {(x), (y), ...t} = {1, 42, 5, 6}; t`)

	AssertCodeErrors(t, "", `let {} = {1}; 1`)
	AssertCodeErrors(t, "", `let {1} = {}; 1`)
	AssertCodeErrors(t, "", `let {42} = {2}; 1`)
	AssertCodeErrors(t, "", `let {42, 43}={41, 42}; 1`)
	AssertCodeErrors(t, "", `let {x, y}={41, 42}; 1`)
	AssertCodeErrors(t, "", `let {x, ...t}={41, 42}; 1`)
	AssertCodeErrors(t, "", `let x = 1; let y = 42; let {(x), (y)} = {1, 4}; 2`)
}

func TestExprLetExtraElementsInPattern(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `let [...] = [1, 2]; 42`)
	AssertCodesEvalToSameValue(t, `1`, `let [x, ...] = [1, 2, 4]; x`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let [x, y, ...] = [1, 2, 4]; [x, y]`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let [x, y, ...] = [1, 2, 4, 8]; [x, y]`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let [x, y, ...] = [1, 2]; [x, y]`)
	AssertCodesEvalToSameValue(t, `2`, `let [..., x] = [1, 2]; x`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let [..., x, y] = [1, 2]; [x, y]`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let [x, ..., y] = [1, 2]; [x, y]`)
	AssertCodesEvalToSameValue(t, `[1, 3]`, `let [x, ..., y] = [1, 2, 3]; [x, y]`)

	AssertCodesEvalToSameValue(t, `[1, 2]`, `let [...t] = [1, 2]; t`)
	AssertCodesEvalToSameValue(t, `[1, [2, 3, 4, 5]]`, `let [x, ...t] = [1, 2, 3, 4, 5]; [x, t]`)
	AssertCodesEvalToSameValue(t, `[1, 2, [3, 4, 5]]`, `let [x, y, ...t] = [1, 2, 3, 4, 5]; [x, y, t]`)
	AssertCodesEvalToSameValue(t, `[2, 3, 4, 5]`, `let [_, ...t] = [1, 2, 3, 4, 5]; t`)
	AssertCodesEvalToSameValue(t, `[2, [3, 4, 5]]`, `let [_, x, ...t] = [1, 2, 3, 4, 5]; [x, t]`)
	AssertCodesEvalToSameValue(t, `[1, 4, 5, [2, 3]]`, `let [x, ...t, y, z] = [1, 2, 3, 4, 5]; [x, y, z, t]`)
	AssertCodesEvalToSameValue(t, `[4, 5, [2, 3]]`, `let [_, ...t, y, z] = [1, 2, 3, 4, 5]; [y, z, t]`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3, []]`, `let [x, ...t, y, z] = [1, 2, 3]; [x, y, z, t]`)

	AssertCodesEvalToSameValue(t, `[1, 2]`, `let (m: x, n: y, ...) = (m: 1, n: 2); [x, y]`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let (m: x, n: y, ...) = (m: 1, n: 2, j: 3, k: 4); [x, y]`)
	AssertCodesEvalToSameValue(t, `[1, 2, ()]`, `let (m: x, n: y, ...t) = (m: 1, n: 2); [x, y, t]`)
	AssertCodesEvalToSameValue(t, `[1, 2, (j: 3, k: 4)]`, `let (m: x, n: y, ...t) = (m: 1, n: 2, j: 3, k: 4); [x, y, t]`)

	AssertCodesEvalToSameValue(t, `[1, 2]`, `let {"m": x, "n": y, ...} = {"m": 1, "n": 2}; [x, y]`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let {"m": x, "n": y, ...} = {"m": 1, "n": 2, "j": 3, "k": 4}; [x, y]`)
	AssertCodesEvalToSameValue(t, `[1, 2, {}]`, `let {"m": x, "n": y, ...t} = {"m": 1, "n": 2}; [x, y, t]`)
	AssertCodesEvalToSameValue(t, `[1, 2, {"j": 3, "k": 4}]`,
		`let {"m": x, "n": y, ...t} = {"m": 1, "n": 2, "j": 3, "k": 4}; [x, y, t]`,
	)

	AssertCodeErrors(t, "", `let [x, y, ...y] = [1, 2, 2]; x`)
	AssertCodeErrors(t, "", `let [x, y, ...y] = [1, 2, 4]; x`)
	AssertCodeErrors(t, "", `let [..., y, ...] = [1, 2, 4]; x`)
}

func TestExprLetNestedPattern(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `let [[[x], y], z] = [[[1], 2], 3]; [x, y, z]`)
	AssertCodesEvalToSameValue(t, `1`, `let {"a": {"b": {"c": x}}} = {"a": {"b": {"c": 1}}}; x`)
	AssertCodesEvalToSameValue(t, `1`, `let (x: (y: (z: w))) = (x: (y: (z: 1))); w`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `let [{"a": x}, (b: y), z] = [{"a": 1}, (b: 2), 3]; [x, y, z]`)
	AssertCodeErrors(t, "", `let [[x]] = []; 42`)
}

func TestExprLetGetPattern(t *testing.T) { //nolint:dupl
	t.Parallel()
	AssertCodesEvalToSameValue(t, `1`, `let {"a"?: x:42} = {"a": 1}; x`)
	AssertCodesEvalToSameValue(t, `42`, `let {"b"?: x:42} = {"a": 1}; x`)
	AssertCodesEvalToSameValue(t, `[1, 2]`, `let {'ids'?: ids:[]} = {'ids': [1, 2]}; ids`)
	// AssertCodesEvalToSameValue(t, `[42, {"a": 1}]`, `let {"b"?: x:42, ...t} = {"a": 1}; [x, t]`)
	AssertCodesEvalToSameValue(t, `1`, `let {"a"?: {"b": {"c"?: x:42}}} = {"a": {"b": {"c": 1}}}; x`)
	AssertCodesEvalToSameValue(t, `42`, `let {"a"?: {"b": {"c"?: x:42}}} = {"a": {"b": {"k": 1}}}; x`)

	AssertCodesEvalToSameValue(t, `1`, `let (a?: x:42) = (a: 1); x`)
	AssertCodesEvalToSameValue(t, `42`, `let (b?: x:42, a: 1) = (a: 1); x`)
	AssertCodesEvalToSameValue(t, `[1, 42]`, `let (a: x, b?: y:42) = (a: 1); [x, y]`)
	AssertCodesEvalToSameValue(t, `1`, `let (x?: (y: (z?: w:42))) = (x: (y: (z: 1))); w`)
	AssertCodesEvalToSameValue(t, `42`, `let (x?: (y: (k?: w:42, z: 1))) = (x: (y: (z: 1))); w`)

	AssertCodesEvalToSameValue(t, `[1, 2, 0]`, `let [x, y, ?z:0] = [1, 2]; [x, y, z]`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `let [x, y, ?z:0] = [1, 2, 3]; [x, y, z]`)
	AssertCodesEvalToSameValue(t, `[1, 2, 0]`, `let [x, [y, ?z:0]] = [1, [2]]; [x, y, z]`)
	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `let [x, [y, ?z:0]] = [1, [2, 3]]; [x, y, z]`)

	AssertCodeErrors(t, "", `let (x?: (k: (z?: w:42))) = (x: (y: (z: 1))); w`)
}

func TestExprLetDynIdentPattern(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `let f = \x x * @{y}; let @{y} = 6; f(7)`)
	AssertCodesEvalToSameValue(t, `42`, `let f = \x x * @{y}; let [@{y}, y] = [6, 8]; f(7)`)
	AssertCodesEvalToSameValue(t, `42`, `let f = \x x * @{y}; let {@{y}, 8} = {6, 8}; f(7)`)
	AssertCodesEvalToSameValue(t, `42`, `let f = \x x * @{y}; let (a:@{y}, ...) = (a:6, b:8); f(7)`)
	AssertCodesEvalToSameValue(t, `42`, `let f = \x x * @{y}; let {1:@{y}, ...} = {1:6, 2:8}; f(7)`)
	AssertCodesEvalToSameValue(t, `[2, 4, 8, 16, 32]`, `let f = \x x ^ @{y}; [1, 2, 3, 4, 5] >> \@{y} f(2)`)
	AssertCodesEvalToSameValue(t, `{[1, 1], [4, 8], [9, 27], [25, 125]}`, `
		let f = \x x ^ @{y}; {1, 2, 3, 5} => [
			let @{y} = 2; f(.),
			let @{y} = 3; f(.),
		]
	`)
}

func TestExprDynLet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `let (@{x}: 42); @{x}`)
	AssertCodesEvalToSameValue(t, `201`, `let f = \x 100 * x + @{x}; let (@{x}: 1); f(2)`)
	AssertCodesEvalToSameValue(t, `201`, `let setter = \n (@{x}: n); let f = \x 100 * x + @{x}; let setter(1); f(2)`)

	AssertCodeErrors(t,
		`bindings not a tuple in "let bindings; expr": {1, 2, 3}`,
		`let {1, 2, 3}; x`,
	)
	AssertCodeErrors(t,
		`"x" not a dynamic name in "let bindings; expr"`,
		`let (x: 42); x`,
	)
}

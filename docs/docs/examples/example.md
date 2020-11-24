---
id: all
title: All Examples
---

### Evaluate an expression

```bash
$ arrai eval '41 + 1'
42
```

### Evaluate count of string

```bash
$ arrai eval '"123456789" count'
9
```

### Evaluate a collection of values

```bash
$ arrai eval '[1,2,3,4] >> .^2'
[1,4,9,16]
```

```bash
$ arrai eval '{1,2,3,4} => .^2'
{1,4,9,16}
```

### Filter a collection of values

```bash
$ arrai eval '{(a:1, b:2), (a:2, b:3), (a:2, b:4)} where .a=2'
{(a:2, b:4)}
```

### Operations for filtering a stream of values

```bash
$ arrai eval '{(a:1, b:2), (a:2, b:3), (a:2, b:4)} where .a=2 and .b=3'
{(a:2, b:3)}
```

```bash
$ arrai eval '{(a:1, b:2), (a:2, b:3), (a:2, b:4)} where .a=2 or .b=4'
{(a:2, b:3), (a:2, b:4)}
```

```bash
$ arrai eval '{(a:1, b:2), (a:2, b:3), (a:2, b:4)} where .a!=2'
{(a:1, b:2)}
```

### Conditional operator

#### Standard Cases

```bash
$ arrai eval 'cond { 2 > 1 : 1, 2 > 3 : 2}'
1
```
Note: Trailing comma is allowed
```bash
$ arrai eval 'cond {
    2 > 1 : 1,
    2 > 3 : 2,
}'
1
```

```bash
$ arrai eval 'cond { 2 > 1 : 1, 2 > 3 : 2, _ : 3}'
1
```

```bash
$ arrai eval 'cond { 2 < 1 : 1, 2 > 3 : 2, _ : 3}'
3
```

```bash
$ arrai eval 'let a = cond { 2 > 1 : 1, 2 > 3 : 2, _ : 3};a * 3'
3
```

```bash
$ arrai eval 'let a = cond { 2 < 1 : 1, 2 > 3 : 2, _ : 3};a * 3'
9
```

```bash
$ arrai eval 'let a = cond { 1 < 2 : 1, 2 > 3 : 2, _ : 3};a * 3'
3
```

```bash
$ arrai eval 'let a = cond { 2 < 1 : 1, 2 < 3 : 2, _ : 3};a * 3'
6
```

```bash
$ arrai eval 'let a = cond { 2 < 1 || 1 > 0 : 1, 2 < 3 : 2, _ : 3};a * 3'
3
```

```bash
$ arrai eval 'let a = cond { 2 < 1 || 1 > 2 : 1, 2 < 3 && 1 > 0 : 2, _ : 3};a * 3'
6
```

```bash
$ arrai eval 'cond {cond {1 > 0 : 1} > 0 : 1, 2 < 3: 2, _:1 + 2}'
1
```

```bash
$ arrai eval 'cond {cond {1 > 2 : 1, _ : 11} < 2 : 1, 2 < 3: 2, _:1 + 2}'
2
```

#### Control Var Cases
```bash
$ arrai eval 'let a = 1; cond a {1 :1, 2 :2, _:1 + 2}'
1
```

```bash
$ arrai eval 'let a = 1; cond a {1 :1 + 10, 2 : 2, _:1 + 2}'
11
```

```bash
$ arrai eval 'let a = 1; cond a {2 :2, _:1 + 2}'
3
```

```bash
$ arrai eval 'let a = 1; let b = cond a {1 :1, 2 :2, _:1 + 2}; b * 100'
100
```

```bash
$ arrai eval 'let a = 1; cond a + 1 {1 :1, 2 :2, _:1 + 2}'
2
```

```bash
$ arrai eval 'let a = 2; cond a { 1: "A", (2, 3): "B", _: "C"}'
B
```

```bash
$ arrai eval 'let a = 2; cond a { (cond a {(1,2) : 1}): "A", (2, 3): "B", _: "C"}'
B
```

```bash
$ arrai eval 'let a = 1; cond a { (cond {2 > 1 : 1}): "A", (2, 3): "B", _: "C"}'
A
```

```bash
$ arrai eval 'let a = 1; cond { cond a {1 : 1} : "A", 2: "B", _: "C"}'
A
```

<!-- TODO: Uncomment once this works again.
### Apply a transform to inbound data

```bash
$ echo {0..10} | arrai transform '2^.'
```

Use `ax` as shorthand for `arrai transform`:

```bash
$ ln -s arrai "$GOPATH/bin/ax"
$ echo {0..10} | ax '2^.'
```
-->

### Filter a collection of values with Control Var cases

```bash
$ arrai eval '{1, [2, 3], 4, [5, 6]} filter . {[a, b]: a + b}'
{5, 11}
```

```bash
$ arrai eval '{1, [2, 3], 4, [5, 6], [7, 8, 9]} filter . {[a, ..., b]: a + b}'
{5, 11, 16}
```

```bash
$ arrai eval '{1, [2, 3], 4, [5, 6]} filter . {[_, _]: 42}'
{42}
```


### Transform JSON and filter a value from it

Filter out value of `a` from JSON:

```bash
$ echo '{"a": "hello", "b": "world"}'| arrai json | arrai x '.("a")'
"hello"
```

### Nest

Nest groups the given attributes into a nested relation, keys on the given key. 

```bash
relation nest |attr1,attr2,...|key
```

#### Examples

```bash
$ arrai eval '{|a,b| (1,2), (1,3), (2,4)} nest |b|nested-b'
{(a: 1, nested-b:{(b: 2), (b: 3)}), (a: 2, nested-b: {(b: 4)})}
```

```bash
$ arrai eval '{|a,b| (1,2), (2,3)} nest |b|nestb'
{(a: 1, nestb:{(b: 2)}), (a: 2, nestb:{(b: 3)})}
```

```bash
$ arrai eval '{|a,b,c| (1,2,3), (1,3,3), (1,2,2)} nest |b|nestb'
{(a: 1, c: 3, nestb: {(b: 2), (b: 3)}), (a: 1, c: 2, nestb: {(b: 2)})}
```

```bash
$ arrai eval '{|a,b,c| (1,2,3), (1,3,3), (1,2,2)} nest |b,c|nestbc'
{(a: 1, nestbc: {(b: 2, c: 3), (b: 3, c: 3), (b: 2, c: 2)})}
```

Nest collects all the b values grouped by the common a values.

### Join

Join takes two relations and joins them by matching tuples on common attribute names.

#### Examples

```bash
$ arrai eval '{|a,b| (1,2), (2,2)} <&> {|a,c| (1,2), (1,3)}'
{(a: 1, b: 2, c: 2), (a: 1, b: 2, c: 3)}
```

```bash
$ arrai eval '{|a,b| (1,2), (2,2)} <&> {|a,c| (1,2), (1,3), (2,3)}'
{(a: 1, b: 2, c: 2), (a: 1, b: 2, c: 3), (a: 2, b: 2, c: 3)}
```

```bash
$ arrai eval '{|a,b| (1,2), (1,3)} <&> {|a,c| (1,2), (1,3)}'
{(a: 1, b: 2, c: 2), (a: 1, b: 2, c: 3), (a: 1, b: 3, c: 2), (a: 1, b: 3, c: 3)}
```

```bash
$ arrai eval '{|a,b,c| (1,2,2), (1,2,3)} <&> {|b,c,d| (2,3,4), (1,3,4)}'
{(a: 1, b: 2, c: 3, d: 4)}
```

### Merge

Merge combines two tuples/dicts, producing a single tuple/dict containing a union of their attributes. If the same name is present in both the LHS (left-hand side) and RHS (right-hand side) tuples/dicts, the RHS value takes precedence in the output.

#### Examples

```bash
$ arrai e '(a: 1, b: 2) +> (b: 3, c: 4)'
(a: 1, b: 3, c: 4)
```

```bash
$ arrai e '(a: 1, b: (c: 2)) +> (b: (c: 4), c: 4)'
(a: 1, b: (c: 4), c: 4)
```

```bash
$ arrai e '(a: (b: 1)) +> (a: (c: 2))'
(a: (c: 2))
```

```bash
$ arrai e '{"a": 1, "b": 2} +> {"b": 3, "d": 4}'
{'a': 1, 'b': 3, 'd': 4}
```

```bash
$ arrai e '{"a": {"b": 1}} +> {"a": {"c": 2}}'
{'a': {'c': 2}}
```

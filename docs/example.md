# Examples

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

### Transform json and filter value from it

Filter out value of `a` from json 

```bash
$ echo '{"a": "hello", "b": "world"}'| arrai json | arrai x '."a"'
"hello"
```

### Nest

Nest groups the given attributes into a nested relation, keys on the given key. 

```bash
relation nest |attr1,attr2,...|key
```

#### Examples

```bash
$ arrai eval '{|a,b| (1,2), (1,3), (2,4) } nest |b|nested-b'
{(a: 1, nested-b:{(b: 2), (b: 3)}), (a: 2, nested-b: {(b: 4)})}
```

```bash
$ arrai eval '{|a,b| (1,2), (2,3) } nest |b|nestb'
{(a: 1, nestb:{(b: 2)}), (a: 2, nestb:{(b: 3)})}
```

```bash
$ arrai eval '{|a,b,c| (1,2,3), (1,3,3), (1,2,2) } nest |b|nestb'
{(a: 1, c: 3, nestb: {(b: 2), (b: 3)}), (a: 1, c: 2, nestb: {(b: 2)})}
```

```bash
$ arrai eval '{|a,b,c| (1,2,3), (1,3,3), (1,2,2) } nest |b,c|nestbc'
{(a: 1, nestbc: {(b: 2, c: 3), (b: 3, c: 3), (b: 2, c: 2)})}
```

Nest collects all the b values grouped by the common a values.

### Join

Join takes two relations and joins them by matching tuples on common attribute names.

#### Examples

```bash
$ arrai eval '{|a,b| (1,2), (2,2) } <&> {|a,c| (1,2), (1,3) }'
{(a: 1, b: 2, c: 2), (a: 1, b: 2, c: 3)}
```

```bash
$ arrai eval '{|a,b| (1,2), (2,2) } <&> {|a,c| (1,2), (1,3), (2,3) }'
{(a: 1, b: 2, c: 2), (a: 1, b: 2, c: 3), (a: 2, b: 2, c: 3)}
```

```bash
$ arrai eval '{|a,b| (1,2), (1,3) } <&> {|a,c| (1,2), (1,3) }'
{(a: 1, b: 2, c: 2), (a: 1, b: 2, c: 3), (a: 1, b: 3, c: 2), (a: 1, b: 3, c: 3)}
```

```bash
$ arrai eval '{|a,b,c| (1,2,2), (1,2,3) } <&> {|b,c,d| (2,3,4), (1,3,4) }'
{(a: 1, b: 2, c: 3, d: 4)}
```

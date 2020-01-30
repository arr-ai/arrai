# Examples

### Evaluate an expression

```bash
$ arrai eval '41 + 1'
```
### Evaluate count of string

```bash
$ arrai eval '"123456789" count'
```

### Evaluate a stream of values
```bash
$ arrai eval '[1,2,3,4] @> .^2'
$ arrai eval '{1,2,3,4} => .^2'
```

### Filter a stream of values
```bash
$ arrai eval '{(a:1, b:2), (a:2, b:3), (a:2, b:4)} where .a=2'
```

### Operations for filtering a stream of values
```bash
$ arrai eval '{(a:1, b:2), (a:2, b:3), (a:2, b:4)} where .a=2 and .b=3'
$ arrai eval '{(a:1, b:2), (a:2, b:3), (a:2, b:4)} where .a=2 or .b=4'
$ arrai eval '{(a:1, b:2), (a:2, b:3), (a:2, b:4)} where .a!=2'
```

### Transform a stream of values

```bash
$ echo {0..10} | arrai transform '2^.'
```
Use `ax` as shorthand for `arrai transform`:

```bash
$ ln -s arrai "$GOPATH/bin/ax"
$ echo {0..10} | ax '2^.'
```

### Transform json and filter value from it
Filter out value of `a` from json 
```bash
$ echo '{"a": "hello", "b": "world"}'| arrai json | arrai x '."a"'
```

### Nest
It collects all the b values grouped by the common a values. 

```bash
$ arrai eval '{ |a,b| (1,2), (1,3), (2, 4) } nest |b|nested-b'
```
Examples:
```bash
$ arrai eval '{ |a,b| (1,2), (2,3) } nest |b|nestb'
$ arrai eval '{ |a,b,c| (1,2,3), (1,3,3), (1,2,2) } nest |b|nestb'
$ arrai eval '{ |a,b,c| (1,2,3), (1,3,3), (1,2,2) } nest |b,c|nestbc'
```
### Join
It matches each relation by the common attributes (matched tables by common column names)
```bash
$ arrai eval '{ |a,b| (1,2), (2,2) } <&> { |a,c| (1,2), (1,3) }'
```
Examples: 
```bash
$ arrai eval '{ |a,b| (1,2), (2,2) } <&> { |a,c| (1,2), (1,3), (2,3) }'
$ arrai eval '{ |a,b| (1,2), (1,3) } <&> { |a,c| (1,2), (1,3) }'      
$ arrai eval '{ |a,b,c| (1,2,2), (1,2,3) } <&> { |b,c,d| (2,3,4), (1,3,4) }'
```
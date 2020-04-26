# arrai

![Go build status](https://github.com/arr-ai/arrai/workflows/Go/badge.svg)

The ultimate data engine.

### Installation

[Install Go](https://golang.org/doc/install), then:

```bash
go get -v -u github.com/arr-ai/arrai/cmd/arrai
arrai -h
```

### Run the interactive shell

```text
$ arrai i
@> 6 * 7
42
@> //.<tab>
archive  dict     encoding eval     fn       grammar  log      math
net      os       re       reflect  rel      seq      str      tuple
unicode
@> //.str.<tab>
contains   expand     has_prefix has_suffix join       lower      repr
split      sub        title      upper
@> //.str.upper("hello")
'HELLO'
```

Ctrl+D to exit.

In a Unix-like environment, you can also (after arr.ai is installed somewhere on
the path):

```bash
$ ln -s arrai $(dirname `which arrai`)/ai
$ ai
@> _
```

### Evaluate an expression

```bash
arrai eval '41 + 1'
```

See the [Introduction to Arr.ai](docs/intro.md) to learn more about the arr.ai
language.

<!-- TODO: Uncomment once this works again.
### Transform a stream of values

```bash
echo {0..10} | arrai transform '2^.'
```

Use `ax` as shorthand for `arrai transform`:

```bash
ln -s arrai "$GOPATH/bin/ax"
echo {0..10} | ax '2^.'
```
-->

### Start a server

```bash
arrai serve --listen localhost
```

### Observe a server

```bash
arrai observe localhost '$'
```

### Update a server

```bash
arrai update localhost '{a: {|1, 2, 3|}, b: "hello"}'
```

#### Arrai Examples

1. [Snippets](docs/example.md)
2. [More complete examples](examples)

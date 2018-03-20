# arrai

The ultimate data engine.

### Installation

[Install Go](https://golang.org/doc/install), then:

```bash
$ go get github.com/arr-ai/arrai/arrai
$ go install github.com/arr-ai/arrai/arrai
$ arrai -h
```

### Evaluate an expression

```bash
$ arrai eval '41 + 1'
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

### Start a server

```bash
$ arrai serve --listen localhost
```

### Observe a server

```bash
$ arrai observe localhost '$'
```

### Update a server

```bash
$ arrai update localhost '{a: {|1, 2, 3|}, b: "hello"}'
```

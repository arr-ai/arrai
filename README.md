# arrai

The ultimate data engine.

### Installation

[Install Go](https://golang.org/doc/install), then:

```bash
$ go get github.com/arr-ai/arrai/arrai
$ go install github.com/arr-ai/arrai/arrai
```

### Stand-alone usage

```bash
$ arrai eval '41 + 1'
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

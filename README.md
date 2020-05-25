# Arr.ai

![Go build status](https://github.com/arr-ai/arrai/workflows/Go/badge.svg)

The ultimate data engine.

## Install

On a Unix-like OS, [install Go](https://golang.org/doc/install) (1.14 or above),
then:

```bash
git clone https://github.com/arr-ai/arrai.git
cd arrai
make install
```

On Windows, download the relevant ZIP file from the
[Releases page](https://github.com/arr-ai/arrai/releases).

## Learn

Follow the [Arr.ai tutorial](docs/tutorial/README.md) for a step by step guide
into the world of arr.ai programming.

See the [Introduction to Arr.ai](docs/README.md) to learn more about the arr.ai
language.

See the [Standard Library Reference](docs/std.md) to learn which are batteries
are included in arr.ai.

### Arr.ai Examples

1. [Snippets](docs/example.md)
2. [More complete examples](examples)

## Use

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

Ctrl+D to exit. or use the `/exit` command.

```bash
@> /exit
```

On Unix-like platforms, you can use the `ai` shortcut:

```bash
$ ai
@> _
```

There are more features in the interactive shell. For more info please read the
[shell tutorial](docs/tutorial/shell.md).

### Evaluate an expression

```bash
arrai eval '41 + 1'
```
Run `arrai help` or `arrai help <command>` for more information.
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

### Run an arrai file

```bash
arrai path/to/file.arrai
```

or use the `run` command

```bash
arrai run path/to/file.arrai
```

`arrai run` can be used to avoid ambiguity between filename and a command.
For example, running an arrai file named `run` (`arrai run run`). Alternatively, include a
subdirectory component in the filename (`arrai ./run`).

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
arrai update localhost '(a: {1, 2, 3}, b: "hello")'
arrai u localhost '$ + (a: $.a | {4, 5, 6})'
```

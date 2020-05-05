# Arr.ai

![Go build status](https://github.com/arr-ai/arrai/workflows/Go/badge.svg)

The ultimate data engine.

## Install

[Install Go](https://golang.org/doc/install), then:

```bash
go get -v -u github.com/arr-ai/arrai/cmd/arrai
arrai -h
```

On Unix-like platforms, you can also symlink a handy shortcut:

```bash
ln -s arrai $(dirname $(which arrai))/ai
```

## Learn

Follow the [Arr.ai tutorial](docs/tutorial/README.md) for a step by step guide
into the world of arr.ai programming.

See the [Introduction to Arr.ai](docs/README.md) to learn more about the arr.ai
language.

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

Ctrl+D to exit. or use the `exit` command.

```bash
@> /exit
```

On Unix-like platforms, you can use the `ai` shortcut:

```bash
$ ai
@> _
```

#### Commands

In the interactive shell, you can use some special commands. You can activate
these commands using the following syntax.

```bash
@> /<command name> [arguments...]
```

Below are the currently provided commands.

#### `/set`

Usage:

```bash
@> /set <name> = <expression>
```

This command adds a variable to the global scope that you can use in multiple
expression in the interactive shell. The provided name has to be alphanumeric,
with no whitespace in it.

Example:

```bash
@> /set pew = "pew pew pew"
'pew pew pew'
@> pew
'pew pew pew'
@> {"do pews": pew}
{'do pews': 'pew pew pew'}
```

The scope of expressions set by `set` are global, which means they can be
replaced locally in an expression. But due to the nature of arrai, scopes are
immutable, so replacing expression locally won't replace it globally.

Example:

```bash
@> /set pew = "pew pew pew"
'pew pew pew'
@> let pew = "no pews"; {"do pews": pew}
{'do pews': 'no pews'}
@> pew
'pew pew pew'
```

#### `/unset`

Usage:

```bash
@> /unset <name>
```

The `unset` command does the opposite of `set`. This command removes a variable
from the global scope.

```bash
@> /set pew = "pew pew pew"
'pew pew pew'
@> pew
'pew pew pew'
@> /unset pew
@> pew
2020-04-28T17:49:19.730553+10:00 error_message=Name "pew" not found in {}

.:1:1:
pew INFO
```

#### `/exit`

Usage:

```bash
@> /exit
```

The `exit` command exits the interactive shell. Alternatively, you can press
Ctrl+D.

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

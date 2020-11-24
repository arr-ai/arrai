---
id: overview
title: Overview
---

# Standard library reference

The arr.ai standard library is available via the `//` syntax.

## Packages

The following standard packages are available:

- [`//archive`](./archive): Archive format utilities
- [`//bits`](./bits): Bit utilities
- [`//encoding`](./encoding): Support for encoded data processing
- [`//eval`](./eval): Evaluate strings holding arr.ai code
- [`//fmt`](./fmt): String formatting utilities
- [`//fn`](./fn): Function manipulation utilities
- [`//grammar`](./grammar): Grammar features & processing
- [`//log`](./log): Print utilities
- [`//math`](./math): Math operations and functions
- [`//net`](./net): Network utilities
- [`//os`](./os): OS support
- [`//re`](./re): Regular expressions
- [`//rel`](./rel): Relational operations
- [`//seq`](./seq): Sequence utilities
- [`//str`](./str): String utilities
- [`//test`](./test): String utilities

## Core functions

The following functions are available at the root of the standard library.

### `//dict(t <: tuple) <: dict`

`dict` converts the tuple `t` to a dictionary.

Usage:

| example | equals |
|:-|:-|
|`//dict(())` | `{}` |
| `//dict((a: 1, b: 2))` | `{"a": 1, "b": 2}` |

### `//tuple(d <: dict) <: tuple`

`tuple` converts the dictionary `d` to a tuple. All keys must be strings. The
operation is shallow, so dictionary-structured values won't be converted.

Usage:

| example | equals |
|:-|:-|
|`//tuple({})` | `()`|
| `//tuple({"a": 1, "b": 2})` | `(a: 1, b: 2)` |

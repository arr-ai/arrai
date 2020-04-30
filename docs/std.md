# Standard library reference

The arr.ai standard library is available via the `//` syntax.

## Packages

The following standard packages are available:

//TODO:
- `//encoding`
  - [`//encoding.json`](std-encoding-json.md): Decode JSON
- [`//eval`](std-eval.md): Evaluate strings holding arr.ai code
- [`//os`](std-os.md): OS support
- [`//re`](std-re.md): Regular expressions
- [`//seq`](std-seq.md): Sequence utilities
- [`//str`](std-str.md): String utilities
- `//unicode`
  - [`//unicode.utf8`](std-unicode-utf8.md): Interoperate with UTF-8 encoding
- math done
- grammar
- fn
- log done
- archive done
- net done
- rel

## Core functions

The following functions are available at the root of the standard library.

### `dict(t <: tuple) -> dict`

Converts the tuple `t` to dictionaries.

Usage:

| example | equals |
|:-|:-|
|`//dict(())` | `{}` |
| `//dict((a: 1, b: 2))` | `{"a": 1, "b": 2}` |

### `tuple(d <: dict) -> tuple`

Converts the dictionary `d` to tuples. All keys must be strings. The
operation is shallow, so dictionary-structured values won't be converted.

Usage:

| example | equals |
|:-|:-|
|`//tuple({})` | `()`|
| `//tuple({"a": 1, "b": 2})` | `(a: 1, b: 2)` |

# Standard library reference

The arr.ai standard library is available via the `//` syntax.

## Packages

The following standard packages are available:

- [`//archive`](std-archive.md): Archive format utilities
- [`//bits`](std-bits.md): Bit utilities
- [`//encoding`](std-encoding.md): Support for encoded data processing
- [`//eval`](std-eval.md): Evaluate strings holding arr.ai code
- [`//fn`](std-fn.md): Function manipulation utilities
- [`//grammar`](../syntax/std-grammar.md): Grammar features & processing
- [`//log`](std-log.md): Print utilities
- [`//math`](std-math.md): Math operations and functions
- [`//net`](std-net.md): Network utilities
- [`//os`](std-os.md): OS support
- [`//re`](std-re.md): Regular expressions
- [`//rel`](std-rel.md): Relational operations
- [`//seq`](std-seq.md): Sequence utilities
- [`//str`](std-str.md): String utilities
- [`//test`](std-test.md): String utilities
- `//unicode`
  - [`//unicode.utf8`](std-unicode-utf8.md): Interoperate with UTF-8 encoding

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

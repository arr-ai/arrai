# Standard library reference

The arr.ai standard library is available via the `//.` syntax.

## Packages

The following standard packages are available:

- `//.encoding`
  - [`//.encoding.json`](std-encoding-json.md): Decode JSON
- [`//.eval`](std-eval.md): Evaluate strings holding arr.ai code
- [`//.os`](std-os.md): OS support
- [`//.re`](std-re.md): Regular expressions
- `//.unicode`
  - [`//.unicode.utf8`](std-unicode-utf8.md): Interoperate with UTF-8 encoding

## Core functions

The following functions are available at the root of the standard library.

### `//.dict`

Converts tuples to dictionaries.

- `//.dict(()) = {}`
- `//.dict((a: 1, b: 2)) = {"a": 1, "b": 2}`

### tuple

`//.tuple` converts dictionaries to tuples. All keys must be strings. The
operation is shallow, so dictionary-structured values won't be converted.

- `//.tuple({} = ()`
- `//.tuple({"a": 1, "b": 2}) = (a: 1, b: 2)`

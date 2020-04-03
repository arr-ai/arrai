# Standard library reference

The arr.ai standard library is available via the `//.` syntax.

## Packages

The following standard packages are available:

- `//.encoding`
  - [`//.encoding.json`](std-encoding-json.md)
- [`//.eval`](std-eval.md)
- [`//.os`](std-os.md)

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

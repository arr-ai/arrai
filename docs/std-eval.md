# eval

The `eval` contains functions which converts raw string into a built-in arrai values.

## `//eval.value(s <: string|array_of_bytes) <: any`

`value` takes in a string or byte array `s` which represents an arrai value and converts them to
arrai values.

However, `eval` is only supported to evaluate simple values e.g. numbers,
tuples, sets, arrays, strings etc.

Usage:

| example | equals |
|:-|:-|
|`//eval.value("'true'")` | `true` |
|`//eval.value("(test: 'test', number:123)")` | `(test: 'test', number:123)` |

# eval

The `eval` contains functions which converts raw string into a built-in arrai values.

## `eval.value(s <: string) <: any`

It takes in a string `s` which represents an arrai value and converts them to
arrai values.

However, `eval` is only supported to evaluate simple values e.g. numbers,
tuples, sets, arrays, strings etc.

Usage:

| example | equals |
|:-|:-|
|`//eval.value("'true'")` | `true` |
|`//eval.value("(test: 'test', number:123)")` | `(test: 'test', number:123)` |

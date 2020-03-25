# strconv library

`eval` has one function which is `value` which takes a string respresenting an
expression and returns an arrai value.

However, `eval` is only supported to evaluate simple values e.g. numbers,
tuples, sets, arrays, strings etc.

## Example

- `//.eval.value("'true'") = true`
- `//.eval.value("(test: 'test', number:123)") = (test: 'test', number:123)`

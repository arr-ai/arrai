# strconv library

`eval` has one function which is `value` which takes a string respresenting an expression and returns an arrai value.

However, `eval` is only supported to evaluate simple values e.g. numbers, tuples, sets, arrays, strings etc.

## Example

Evaluating `//.eval.value("'true'")` will return `true`
Evaluating `//.eval.value("(test: 'test', number:123)")` will return `(test: 'test', number:123)`

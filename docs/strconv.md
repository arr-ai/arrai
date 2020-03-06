# strconv library

`strconv` has two functions which is `eval` and `unsafe_eval`

Both takes a string of an arrai expression and returns an arrai value which is evaluated from the string input.

However, `eval` is only supported to evaluate simple values e.g. numbers, tuples, sets, arrays, strings etc.

`unsafe_eval` can evaluate more complex operations and return the value e.g. math operations, functions.

## Example

Evaluating `//.strconv.eval("'true'")` will return `true`
Evaluating `//.strconv.eval("(test: 'test', number:123)")` will return `(test: 'test', number:123)`
Evaluating `//.strconv.unsafe_eval("6 * 7")` will return `42`

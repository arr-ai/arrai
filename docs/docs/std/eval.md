The `eval` contains functions which converts raw string into a built-in arrai values.

## `//eval.value(s <: string|array_of_bytes) <: any`

`value` takes in a string or byte array `s` which represents an arrai value and converts them to
arrai values.

However, `eval` is only supported to evaluate simple values e.g. numbers,
tuples, sets, arrays, strings etc.

Usage:

| example | equals |
|---|---|
|`//eval.value("'true'")` | `true` |
|`//eval.value("(test: 'test', number:123)")` | `(test: 'test', number:123)` |

## `//eval.evaluator(config <: (:scope <: tuple, :stdlib <: tuple)).eval(expr <: string|bytes) <: any`

`evaluator` takes in a `config` tuple containing optional configurable values of scope and stdlib. 
The scope config property is a tuple whose keys are bound to corresponding values within the evaluated expression.
The stdlib config property is a tuple containing the hierarchical representation of standard library functions to be made available to the evaluated expression.
By default `stdlib` is equal to `//std.safe`.
In instances where consumers of evaluator want to provide access to unsafe standard library functions then they must be manually provided in the config.


`eval`takes in a string or byte array `expr` which represents an arrai value and converts them to
arrai values.

Usage:

| example | equals |
|---|---|
|`//eval.eval("1+2")` | `3` |
|`//eval.eval("//str.upper('cat')") ` | `"CAT"` |
| `//eval.evaluator((stdlib: //std.safe +> (os+>: (file: //os.file)))).eval("//os.file(...)")` | `file contents` |
|`let double = \d d * 2; //eval.evaluator((scope: (:double))).eval("double(1 + 2)")` | `6`

## `//eval.eval(expr <: string|bytes) <: any`

`eval` is functionally equivalent to `//eval.evaluator((stdlib: //std.safe)).eval()` and takes in a string or byte array `expr` which represents an arrai value and converts them to
arrai values.

Usage:

| example | equals |
|---|---|
|`//eval.eval("1+2")` | `3` |
|`//eval.eval("//str.upper('cat')") ` | `"CAT"` |

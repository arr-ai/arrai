# fn

The `fn` library contains helper functions that are related to arrai `functions`.

## `fix(f <: function) <: function`

Due to the current implementation of arrai, recursive functions cannot be created
without the helper function `fix`. It takes the function recursive `f` and returns
`f` whose recursive nature can now be utilized. For more information, please read
the functions section in [intro](intro.md).

This function uses the concept of a [fixed-point combinator](https://en.wikipedia.org/wiki/Fixed-point_combinator).

Usage:

| example | equals |
|:-|:-|
| `let fib = //fn.fix(\fib (\n (cond (n = 0: 0, n = 1: 1, n > 0: fib(n-1) + fib(n-2), *: 0)))); fib(10)` | `55` |

## `fixt(f <: tuple_of_functions) <: tuple_of_function`

`fixt` is a variant of fix. This allows mutual recursion. For more information, please read the functions section in [intro](intro.md).

Usage:

| example | equals |
|:-|:-|
| `let collatz = //fn.fixt((even: \t \n (cond (n % 2 = 0: [n] ++ t.even(n / 2), *: t.odd(n))), odd: \t \n (cond (n = 1: [1], n % 2 = 1: [n] ++ t.odd(3 * n + 1), *: t.even(n)))); collatz.even(6)`| `[6, 3, 10, 5, 16, 8, 4, 2, 1]` |

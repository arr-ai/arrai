---
id: exprs
title: Expressions
---

### Logic expressions

Arr.ai supports operations on "true" and "false" values. The values `0`, `()`
and `{}` are considered "false", while all other values are "true".

1. `expr1 if testexpr else expr2` evaluates to `expr1` if `testexpr` is "true",
   or `expr2` otherwise.
2. `expr1 && expr2` evaluates to `expr1` if it is "true" or `expr2` otherwise.
3. `expr1 || expr2` evaluates to `expr1` if it is "false" or `expr2` otherwise.

All above expressions exhibit short-circuit behaviours, which means that that
`expr2` will be evaluated if its value is needed. While the arr.ai language has
no side-effects, short-circuit behaviour is still needed to terminate recursion.

### Arithmetic expressions

Arr.ai supports operations on numbers.

1. Unary: `+`, `-`
2. Binary:
   1. Well known: `+`, `-`, `*`, `/`, `%` (modulo), `^` (power)
   2. Modulo-truncation: `-%` (`x -% y = x - x % y`)
3. Comparison operators, which may be chained: `0 <= i < 10`
   1. Set membership is treated the same: `10 <= n <: validIds`.

### Structure access expressions

1. Tuple attribute: `tuple.attr` (string syntax is allowed, e.g.: `('👋': 42)."👋"`))
2. Dot variable attribute: `.attr` (shorthand for `(.).attr`)
3. Function call:
   1. `[2, 4, 6, 8](2) = 6`, `"hello"(1) = 101`
   2. `{"red": 0.3, "green": 0.5, "blue", 0.2}("green") = 0.5`
4. Conditional accessor syntax: allows for failures in accessing a tuple attribute or a set call, falling back on a provided expression. Any call or attribute access that ends with `?` are allowed to fail.
   1. `(a: 1).b?:42 = 42`
   2. `(a: 1).a?:42 = 1`
   3. `{"a": 1}("b")?:42 = 42`
   4. `{"a": 1}("a")?:42 = 1`

   It also allows for appending access expressions:
   1. `(a: {"b": (c: 2)}).a?("b").c?:42 = 2`
   2. `(a: {"b": (c: 2)}).a?("b").d?:42 = 42`

   Not all access failures are allowed: only missing attributes of a tuple, or a set call does not return exactly one value.
   1. `(a: (b: 1)).a?.b.c?:42` will fail as it will try to evaluate `1.c?:42`.
5. Function slice: (**⛔ NYI**)
   1. `[1, 1, 2, 3, 5, 8](2:5) = [2, 3, 5]`
   2. `[1, 2, 3, 4, 5, 6](1:5:2) = [2, 4]`

### Binding expressions

The following operators bind `name` to something related to `expr1` (details
below) and evaluates expression `expr2` with `name` in scope.

1. **`let name = expr1; expr2`** or **`expr1 -> \name expr2`**:
   Evaluates `expr2` with `expr1` in scope as `name`.
2. **`expr1 => \name expr2`**: Transforms each element of set `expr1` and
   evaluates to the set of results.
3. **`expr1 >> \name expr2`**: Transforms each item of keyed-collection
   `expr1` and evaluates to the key-collection of results, with each result being
   associated with the same key that the original item was. This works for any
   binary relation with an `@` attribute, which includes strings, arrays,
   functions and other structures.
4. **`expr1 :> \name expr2`**: Binds `name` to each value in tuple `expr1`,
   evaluates `expr2` and reassociates each result with the corresponding
   name, producing a new tuple.

If `expr1` is omitted in any of the arrow forms, `.` is assumed.

If `\name` is omitted, `\.` is assumed.

<!-- TODO: Examples -->

### Relations

Relations are sets of tuples with a common set of names across all tuples. They
are analogous to SQL tables. Numerous [relational operators](./relops) exist that work on these
structures.

### Functions

There are several flavors of functions. All functions are binary relations with
one attribute called `@`. The other attribute can have any name, including the
empty name, `''`. The following are some examples of functions.

1. **Strings:** `"hello"(2) = 108` (`l`)
2. **Arrays:** `[10, 15, 20, 25, 30](3) = 25`
3. **Lambda functions:** `\x 2 * x`

Unlike most other languages, arr.ai are no concept of named functions, either at
file level or any other scope. All functions are anonymous. A function can, of
course, be bound to a name via `let` or `->`, but, since it cannot refer to this
name at the moment of assignment, this presents a challenge for implementing
recursion. This problem is solved by a couple of functions in the standard
library:

1. **`//fn.fix`** is a fixed-point combinator. It is typically used to
   transform non-recursive functions into recursive ones, e.g.:

   ```arrai
   let factorial = //fn.fix(\factorial (\n (1 if n < 2 else n * factorial(n - 1))));
   factorial(6)
   ```

2. **`//fn.fixt`** is a variant of `fix` that operates on tuples of functions
   instead of a single function. This allows mutual recursion, e.g.:

   ```arrai
   let eo = //fn.fixt((
      even: \t \n n == 0 || t.odd (n - 1),
      odd:  \t \n n != 0 && t.even(n - 1),
   ));
   eo.even(6)
   ```

However, these functions are also available through the syntactic sugar in the
following syntax:

1. For regular recursive functions:
```arrai
let rec factorial = \n 1 if n < 2 else n * factorial(n - 1); factorial(5)
```

2. For mutual recursion:
```arrai
let rec oe = (
   even = \n n == 0 || oe.odd (n - 1),
   odd  = \n n != 0 && oe.even(n - 1),
);
oe.even(6)
```

It is also possible to use the same syntax in a tuple.

```arrai
let t = (
   rec fact: \n cond n ((0, 1): 1, n: n * fact(n - 1)),
   n       : 5
);
t.rec(t.n)
```

This syntactic sugar only works with expression that evaluates to either a
function or a tuple of functions. Anything else and the expression will fail.

### Packages

External libraries may be accessed via package references.

1. **`//`** Is the root of the standard library. It provides access to many
   packages providing a wide range of useful capabilities. The following is a
   small sample of the full set:
   1. **`//math`:** math functions and constants such as `//math.sin`
      and `//math.pi`.
   2. **`//str`:** string functions such as `//str.upper` and
      `//str.lower`.
   3. **`//fn`:** higher order functions such as `//fn.fix` and `//fn.fixt`.
   See the [standard library reference](../std/overview) for full documentation on all packages.
2. **`//{./path}`** provides access to other arrai files relative to the current
   arrai file's parent directory (current working directory for expressions such
   as the `arrai eval` source that aren't associated with a file).
3. **`//{/path}`** provides access to other arrai files relative to the root of
   the current module, looking for `go.mod` file backwards from the current directory.
4. **`//{hostname/path}`** provides access to content from the internet
   1. **`//{github.com/foo/bar/baz}`:** access `baz.arrai` file in remote repository `github.com/foo/bar`
   2. **`//{github.com/foo/bar/a.json}`:** access `a.json` file in remote repository `github.com/foo/bar`
   3. **`//{foo.org/bar/}'random.arrai'`/`//{https://foo.org/bar/random.arrai}`:**
      request content of `https://foo.org/bar/random.arrai` via HTTPS
   4. **`//{foo.org/bar/some.json}`/`//{https://foo.org/bar/some.json}`:**
      request content of `https://foo.org/bar/some.json` via HTTPS
   5. **`//{foo.org/bar/some.yaml}`/`//{https://foo.org/bar/some.yml}`:**
      request content of `https://foo.org/bar/some.yaml` via HTTPS, file extension can be `yml` or `yaml`

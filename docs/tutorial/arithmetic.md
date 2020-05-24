# Arithmetic and logical operators

## Core arithmetic operators

Arr.ai supports most conventional arithmetic operators over numbers. The usual
suspects, `+`, `-`, `*` and `/` perform the four basic operations If you're new
to programming, `*` means × ("times" or "multiplied by") and `/` means ÷
("divided by").

```arrai
@> 1 + 2 * 3
@> 4 - 5 / 7
```

Operators follow standard precedence rules, so the above expressions first
evaluate the multiplication and division, and then evaluate the addition and
subtraction. Thus, `1 + 2 * 3` is equivalent to `1 + (2 * 3)`.

## Parentheses

You can use parenthesis to force an evaluation order that doesn't follow regular
precedence rules:

```arrai
@> (1 + 2) * 3
@> (4 - 5) / 7
```

## Negative and positive

If you want the negative of a number, you can use the prefix `-` operator. `-x`
is effectively the same as `0 - x`. For completeness, prefix `+` is also
available and simply evaluates to its argument.

```arrai
@> -42
@> -42 + 41
@> -42 + -42
@> 42 - -42
@> 42--42
@> 42 + +42
@> 42++42  # FAIL. ++ is a different operator
```

Many languages overload `+` to also mean concatenation, e.g. `"hots" + "hots"`
in Python, Go and JavaScript. Arr.ai avoids this in order to promote clarity of
expression. `+` always means mathematical addition, never anything else. The
`++` operator defines concatenation. It will be discussed in more detail in a
later tutorial.

## Exponentiation

Arr.ai can exponentiate, *x*<sup>*y*</sup>, using the `^` or "power" operator.
`x ^ y` represents *x* raised to the *y*th power.

```arrai
@> 2 ^ 3
@> 3 ^ 2
@> 2 ^ 0.5
@> 10 * 2 ^ -3
@> 4 ^ 3 ^ 2
```

`^` has higher precedence than the other four operators. It also associates
right-to-left. The final example above was interpreted as `4 ^ (3 ^ 2)`, not `(4
^ 3) ^ 2`. From a mathematical point of view, right-to-left is more useful,
since `(a ^ b) ^ c` can be reduced to a single exponentiation, `a ^ (b * c)`,
whereas there is no reduced form of `a ^ (b ^ c)`.

## Logical operators

The logical operators work on Boolean values. Actually, every value serves as a
Boolean value. When used with logical operators, every value behaves either like
`true` or like `false`. The values `0`, `()` and `{}` are all treated as `false`
by the logical operators (in fact, `{} = false`). All other values are treated
as `true`. Some logical operators leverage this property by returning the first
value that determines the outcome.

(When trying out the examples below, remember that `false` is equal to `{}`,
which is also how it prints out.)

```arrai
@> true && true
@> true && false
@> 1 && () && 0
@> 1 && 0 && ()
@> false || false
@> false || 0 || {1: 2}
@> false || true
@> 0 || 42
@> 0 || (x: 1) || ()
@> !(x: 1)
@> !{}
@> !!{1, 2, 3}
```

## Bitwise operators? Nope!

Arr.ai omits bitwise operators by design. Many other languages use integers to
encode bitwise information. Such integers are often called bitmasks. Logically
speaking, a `uint64` in Go may be thought of as a set of numbers, *n*, in the
range 0 &le; *n* < 64. For a given `uint64`, it can be said that the set
contains the number *n* if bit *n* is set to one. For example, the number 75 has
the bit pattern `01001011`. Since only bits 0, 1, 3 and 6 are set to 1, the
number 75 represents the set `{0, 1, 3, 6}`.

That last point offers a pretty strong hint as to why arr.ai doesn't support
bitwise operators. Put simply, they are not needed. Arr.ai can work directly
with sets of integers without having to fake it using bitmasks. In fact, it
supports sets with integers of any size, including negative numbers, as well as
fractions, infinity, and you can even mix in some non-integers if your problem
demands it.

```arrai
@> /set a = {1, 2, 3, 4}
@> /set b = {3, 4, 5, 6}
@> /set c = {5, 6, 7, 8}
@> [a & b, b & c, a & c]     # intersection
@> [a | b, b | c, a | c]     # union
@> [a &~ b, b &~ c, a &~ c]  # difference (aka AND-NOT)
@> [a ~~ b, b ~~ c, a ~~ c]  # symmetric difference (aka XOR)
@> {1, 2, 0.5, //math.pi^0.5, "duck"}  # Mix it up!
```

In real world languages, especially those with no built-in concept of sets,
numeric bitmasks are often used to represent integer sets. To support
interoperability, arr.ai provides the `//bits` package to convert between
integer sets and their numeric representations.

```arrai
@> //bits.set(42)
@> //bits.mask({1, 3, 5})
@> //bits.mask({-1, -3})
```

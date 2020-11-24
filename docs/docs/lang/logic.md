---
id: logic
title: Logical operators
---

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

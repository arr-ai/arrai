---
id: literals
title: Literals
---

## Core literals

The core syntax for literals can express numbers, tuples and sets.

### Numbers

`0`, `1`, `-2`, `3.45e-6`, `7+8.9i`, `9969216677189303386214405760200`

The parts may be written in the following forms:

- Decimal: `123`
- **(⛔NYI)** Use spaces to break up long numbers: `12 345 678`
- **(⛔NYI)** Hexadecimal: `0x7b`
- **(⛔NYI)** Octal: `0o173`
- **(⛔NYI)** Binary: `0b 111 1011`

### Tuples

`()`, `(a:1)`, `('t.m.o.l.': 42)`, `(x: (a: (), b: 2), y: -3i)`, `(:x, :y)`

Like structs in the C family of languages, names are not values in their own right. They cannot be stored in variables or data structures and therefore cannot be manipulated as values. They serve only to specify which element of a tuple is being specified or retrieved.

Unlike C structs, names can be any sequence of characters, with string syntax allowing characters not permitted in identifiers. Also unlike C structs, tuples do not have to conform to definitions stipulating the available fields or the types of values they can hold. A tuple can have any fields and each fields can hold any value of any type.

As an extension to the normal `key: value` syntax, attributes may omit `key` if `value` is an expression of the form `name` or `expr.name`. E.g.: `(:x, :.y, :a.b.z) = (x: x, y: y, z: a.b.z)`.

### Sets

`{}`, `{1, 2, 3}`, `{(a:1, b:2), (a:4, b:7)}`, `{2, {}, (c:4)}`


## Sugared literals

As explained earlier, many other structures are expressible beyond just numbers, tuples and sets.

It is important to remember that these other structures are simply special arrangements of the base types. They do, however, give arr.ai the flavor and power of much richer type systems while retaining a remarkably simple data model.

Also, because these sugared forms are all just the base types in disguise, all of the expressive machinery designed for numbers, tuples and sets can be applied to strings, arrays, etc.

### Boolean

Arr.ai takes a leaf out of the C89 playbook and omits Boolean types from the base type systems. Nonetheless, `false` and `true` are defined in the core language as aliases for the following sets.

1. `false = {}`
2. `true = {()}`

These are not the only values that may be used in logical operations. All values can be tested for "trueness". Most values are considered "true". The only exceptions are `0`, `()` and `{}`.

### Character

A character can be expressed in arr.ai as a `Number`. Its syntactic sugar uses the form of `%char`. The syntax will evaluate to a `Number` whose value corresponds to the ASCII code of `char`.

Usage:

```arrai
%a  = 97
%A  = 65
%\n = 10
%\t = 9
%🙂 = 128578
```

### Relation

A relation is a set of tuples with the same attributes. For example:

```arrai
{
   (acctid: 1, descr: "ACME Corp", balance: 123456789.01),
   (acctid: 2, descr: "Francis Jones", balance: 4567.23),
}
```

The following sugared types (strings, bytes, arrays, and dictionaries) are all relations with two attributes: `@` for index/position, and `@something` for value. What that `something` is determines how arr.ai interprets the values (e.g. for printing), but they are ultimately all sets of tuples.

Arr.ai allows a shorthand form to represent relations:

```arrai
{|acctid, descr          , balance     |
 (     1, "ACME Corp"    , 123456789.01),
 (     2, "Francis Jones", 4567.23     ),
}
```

### String

Strings may be expressed in arr.ai. They are syntactic sugar for relations of the form `{|@, @char| ...}`.

Strings may be expressed in three different forms:

```arrai
"abc"
'abc'
`abc`
```

The three forms differ only in their escaping rules.

1. The double- and single-quoted forms have the same set of escapes, roughly following C string syntax, the only difference being that, in `"..."` strings, `"` requires escaping via `\"`.
2. The same applies for `"` in `"..."` strings.
3. Backquoted strings support no escaping other than the backquote character, which may be escaped with a double backquote:

   ```arrai
   `Let's escape some ``backquotes``!`
   ```

### Expression string

Expression strings appear on the surface to be quite similar to regular strings:

```arrai
$"abc"
$'abc'
$`abc`
```

They are, however, a very powerful text templating mechanism that allows arbitrarily complex nestings of strings and logic. For example, the following expression:

```arrai
let lib = (
   functions: [
      (name: "square", params: ["x"], expr: "x ^ 2"),
      (name: "sum", params: ["x", "y"], expr: "x + y"),
   ]
);
$`${lib.functions >> $`
   function ${.name}(${.params::, }) {
      return ${.expr}
   }
`::\i:\n}`
```

Outputs the following text:

```arrai
function square(x) {
   return x ^ 2
}

function sum(x, y) {
   return x + y
}
```

For a detailed description, see [Expression strings](./exprstr).

### Bytes

Array of Bytes can be expressed in arr.ai. The syntactic sugar is in the form of
`<< expr1, expr2, expr3, ... >>`. They represent relations of the form `{|@, @byte| ...}`.

It only accepts expressions that are evaluated to either a `Number` whose values range from 0-255 inclusive or a `String` with `0` offset.

Any complicated expressions need to be surrounded by parentheses `(expr)`, except literal values such as `Number`, `String`, `Char`, and variables. Any other values and the expression will fail.

A `Number` is appended to the array while each characters of a `String` is appended to the array. The result is an array of Bytes and each Byte is represented as a `Number`.

Example of usages:

```arrai
<<"hello", 10>>      = <<"hello\n">>
<<97, 98, 99>>       = <<"abc">>
<<("abc" >> . + 1)>> = <<"bcd">>
```

### Array

Arrays may be expressed using the conventional `[...]` notation, e.g.: `[1, 2, [3, 4]]`. They represent relations of the form `{|@, @item| ...}`.

Sparse arrays or arrays with holes can also be defined. For example, `[1, 2, , 3]`. This is equivalent to `{|@, @item| (0, 1), (1, 2), (3, 3)}`.

However, holes must be defined in the middle of the elements, which means you can not defined `[, , 1, 2]`. Should you want to define that, you can use the offset syntax `2\[1, 2]`.

Any empty elements at the end will be trimmed which means `[1, 2, 3, ] = [1, 2, 3]`

### Dictionary

Dictionaries, or sets of arbitrary key-value pairs, may be expressed using the conventional `{key: value}` notation, e.g.: `{"a": 1}`, `{1: "a"}`, `{(a: 1): {1, [2, 3], {"b": 4}}}`. They represent relations of the form `{|@, @value| ...}`, where both key and value can be of any type.

#### Tuples vs Dictionaries

It may not be immediately obvious why tuples and dictionaries exist as distinct kinds of values. Firstly, there is a practical reason: dictionaries can have any kind of value as keys:

```arrai
{
   "x":                 "red",
   [1, 2]:              "green",
   (a: [3], b: {5, 6}): "blue",
}
```

A more important distinction is that tuples should be used to capture various known dimensions of a concept, whereas dictionaries are more appropriate to map from an arbitrary or unbounded set of values to some associated values.

As a rule of thumb, if the set of keys is fixed and known at compile-time, it should be a tuple. If not, use a dictionary.

For example, a collection of cars by license plate should be modeled as a dictionaries, since the set of license plates is unbounded. The details of each car, however, form a closed set of known attributes, which should be expressed as tuples:

```arrai
# Map
{
   "ILVME-23": (        # Tuple
      make:  "Porsche",
      model: "911",
      year:  1964,
   ),
   "ZUM-888": (         # Tuple
      make:  "Bugatti",
      model: "Veyron",
      year:  2005,
   ),
}
```

If you find yourself needing to do more dynamic operations with tuples, they can be converted to equivalent `{string: value}` dictionaries with the [`//dict`](../std/overview) standard library function. And of course dictionaries with only strings as keys can be converted to tuples with [`//tuple`](../std/overview).

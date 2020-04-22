# Slicing

Most sets in arrai can be sliced with the slice expression.

The slicing expression can be used by putting the expression as an argument to a call.

The syntax of the slicing expression is as follows:

```text
start=expr? ";" end=expr? (";" step=expr)?
```

As shown above, you can define a slice by defining a lower bound, an upper bound, both or none.
You can also define a `step` that is used to define the increment of the numbers from the lower
bound to the upper bound. When `step` is not defined, its value defaults to `1`.

## Usages

Only certain types of sets can be used with the slicing expression. These sets are:

1. `Array`
2. `String`
3. `Byte Array`
4. `Dictionary`

Generic sets (i.e. unindexed sets) or any other expression that does not evaluates to the
special types of sets will fail the expression.

Example usages:

| expression | equals |
|:-|:-|
| `[1, 2, 3, 4, 5](1;4)` | `[2, 3, 4]` |
| `"this is a sentence"(;7)` | `"this is"` |
| `{1: "hello", 2: "there", 3: 123, 4: 456}(;4)` | `["hello", "there", 123]` |
| `"sentence"(;)` | `"sentence"` |
| `[1, 2, 3, 4, 5](1;5;2)` | `[2, 4]` |
| `[1, 2, 3, 4, 5](3;0;-1)` | `[4, 3, 2]` |
| `[1, 2, 3, 4, 5](;1;-1)` | `[5, 4, 3]` |
| `[1, 2, 3, 4, 5](3;;-1)` | `[4, 3, 2, 1]` |
| `[1, 2, 3, 4, 5](;;-1)` | `[5, 4, 3, 2, 1]` |

Since the expression is also an argument, the following expression is also valid

```text
[1, 2, 3, 4, 5](;;-1, 0)
```

this expression reverses the array and fetch the first value. In
this example, it will evaluate to `5`.

### Behaviour

#### Default values

Slice expression behaves differently when provided with different expressions.

When the lower bound, the upper bound, and step are provided, the values will have
to evaluate to a number as slicing can only be done with number expressions.
Anything other than `Number`, the expression will fail.

When `step` is not defined, its value defaults to `1`. If `step` evaluates to `0`,
it will return an empty set. The value of `step` determines the default value of
`start` and `end`.

When `start` or `end` are defined, the values will be used. However, when the range
is invalid (i.e `start > end && step > 0` or `start < end && step < 0`), an empty
set will be returned as the result of the slice expression.

When `start` or `end` are not defined, the value of `step` will determine their
values.

If `step` is positive, `start` will have the lowest possible value. In
`Array` and `String`, the value of `start` defaults to the offset value. In
`Dictionary`, `start` will defaults to the lowest numerical key.
`end`, on the other hand, will have the highest value possible. In `Array` and
`String`, the value of `end` defaults to the `offset + length of array`. In
`Dictionary`, it will defaults to the `highest numerical key + 1`.
`1` is added to account for inclusivity of the last value (`end` is not defined,
so last value will be included).

For example:

| expression | default start | default end |
|:-|:-|:-|
| `[1, 2, 3, 4, 5]` | 0 |  5 |
| `2\"abcde"` | 2 |  7 |
| `{1: 10, 2: 20, 3: 30}` | 1 | 4 |

If `step` is negative, everything is inverted. `start` will have the highest
possible value. In `Array` and `String`, the value of `start` defaults to
`length + offset - 1`. In `Dictionary`, it defaults to the highest numerical key.
For `end`, it defaults to the lowest value. In `Array` and `String`, it defaults
to `offset - 1` to account for inclusivity of the last element. For `Dictionary`,
it defaults to `lowest numberical key - 1`, also to account of the inclusivity.

For example:

| expression | default start | default end |
|:-|:-|:-|
| `[1, 2, 3, 4, 5]` | 4 |  -1 |
| `2\"abcde"` | 6 |  1 |
| `{1: 10, 2: 20, 3: 30}` | 3 | 0 |

#### Negative Index

Slicing in arrai supports slicing to negative index. This can be used in both
`start` and `end`. `-1` means the last value, `-2` means the second last value,
`-3` means the third last value, and so on.

For example:
| expression | equals |
|:-|:-|
| `[1, 2, 3, 4, 5](1;-1)` | `[2, 3, 4]` |
| `[1, 2, 3, 4, 5](-3;-1)` | `[3, 4]` |

#### Invalid Index

Sometimes the provided indexes can be invalid. For `Array` or `String`, an invalid
index can be any values that are less than the offset or more than the
`offset + length`. For Dictionary, invalid index simply means index that doesn't
exist as a key.

Arrai slicing is quite forgiving. In `Array` and `String`, any values less than
offset are ignored and replaced to the offset value. Any value larger than
`offset + length` are given the same treatment and replaced to `offset + length`.
In `Dictionary`, indexes that do not exist in the `Dictionary` will just be ignored.

#### Dictionary

Slicing in `Dictionary` is quite unique. Slicing in dictionary will always return
an array of values. However, if the range is invalid or it results in empty values,
slicing will return an empty `Set`. Also, since slicing only works with numerical
values (for `start` and `end`), any non-numerical keys will be ignored.

For example

| expression | equals |
|:-|:-|
| `{1: 1, 2: 2, "c": 3, 4: 4}(1;)` | `[1, 2, 4]` |
| `{"a": 1, "b": 2, "c": 3, "d": 4}(1;)` | `{}` |

## Inclusivity **(⛔ NYI)**

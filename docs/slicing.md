# Slicing

Most sets in arrai can be sliced with the slice expression.

The slicing expression can be used by putting the expression as an argument to a call.

The syntax of the slicing expression is as follows:

```text
someSet(optionalLowerBoundExpr;optionalUpperBoundExpr;optionalStepSizeExpr)
```

As shown above, you can define a slice by defining a lower bound, an upper bound, both or none.
You can also define a `stepSize` that is used to define the increment of the numbers from the lower
bound to the upper bound. `stepSize` defaults to `1` when it is not defined.

Lower bounds may be negative and upper bounds may be larger than the set size. In such cases, the index will wraparound. Refer to the example below.

## Usages

Only certain types of sets can be used with the slicing expression. These sets are:

1. `Array`
2. `String`
3. `Byte Array`
4. `Dictionary`

Example usage:

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


Generic sets (i.e. unindexed sets) or any other expression that does not evaluate to the
special types of sets will fail the expression.

An example of an unindexed sets is as follows:

```text
{1, 2, 3, "abc"}
```

Since the expression is also an argument, the following expression is also valid

```text
[1, 2, 3, 4, 5](;;-1, 0)
```

this expression reverses the array and fetch the first value. In
this example, it will evaluate to `5`.

### Behaviour

#### Default values

Slice expression behaves differently when provided with different expressions.

Lower bound and upper bound are always exclusive. As in it always includes values
that corresponds to ranges of values from `start` to `end - 1`. This is true for
all scenarios, whether `start` or `end` are defined or not.

When the lower bound, the upper bound, and step are provided, the values will have
to evaluate to a number as slicing can only be done with number expressions.
Anything expression with a type other than `Number` will cause the expression
evaluation to fail with an error message showing which expression does not compile
to a `Number`.

When `stepSize` is not defined, its value defaults to `1`. If `stepSize` evaluates to `0`,
it will return an empty set. The value of `stepSize` determines the default value of
`start` and `end`.

When `start` or `end` are defined, the values will be used. However, when the range
is invalid (i.e `start > end && step > 0` or `start < end && step < 0`), an empty
set will be returned as the result of the slice expression.

When `start` or `end` are not defined, the value of `stepSize` will determine their
values. The table below shows what the values for `start` and `end` defaults to.

| step | type | start | end |
|:-|:-|:-|:-|
| positive number | `Array`/`String` | `offset` | `length + offset` |
| positive number | `Dictionary` | `lowest numerical key` | `highest numerical key + 1` |
| negative number | `Array`/`String` | `length + offset - 1` | `offset - 1` |
| negative number | `Dictionary` | `highest numerical key` | `lowest numerical key - 1` |

In `end`, `1` is added or subtracted so that it includes the last value.

For example, on `step > 0`:

| expression | default start | default end |
|:-|:-|:-|
| `[1, 2, 3, 4, 5]` | 0 |  5 |
| `2\"abcde"` | 2 |  7 |
| `{1: 10, 2: 20, 3: 30}` | 1 | 4 |

For example, on `step < 0`:

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

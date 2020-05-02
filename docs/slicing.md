# Slicing

Functions (which includes strings, arrays, etc.) can be called using slice syntax, The syntax is as follows:

```text
someSet(optionalLowerBoundExpr;optionalUpperBoundExpr;optionalStepSizeExpr)
```

As shown above, you can define a slice by defining a lower bound, an upper bound, both or none.
You can also define a `stepSize` that is used to define the increment of the numbers from the lower
bound to the upper bound. `stepSize` defaults to `1` when it is not defined.

## Usages

Only certain types of sets can be sliced. These sets are:

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

Since slice is also an argument, the following expression is also valid

```text
[1, 2, 3, 4, 5](;;-1, 0)
```

this expression reverses the array and fetches the first value. In
this example, it will evaluate to `5`.

### Behaviour

#### Default values

Slice operation behaves differently when provided with different expressions.

Lower bound is always inclusive while upper bound is always exclusive. As in it
always includes values that corresponds to ranges of values from `start` to
`end - 1`. This is true for all scenarios, whether `start` or `end` are defined or not.
However, when `end` is not provided, the last element is always included. This is not
a special case as default value of `end` will include the last element. Refer to the
explanations below.

When the lower bound, the upper bound, and step are provided, the values will have
to evaluate to a number as slicing can only be done with number expressions.
Any expression with a type other than `Number` will cause the expression
evaluation to fail with an error message showing which expression does not compile
to a `Number`.

When `stepSize` is not defined, its value defaults to `1`. If `stepSize` evaluates to `0`,
the expression will **fail**. The value of `stepSize` determines the default value of
`start` and `end`.

When `start` or `end` are defined, the values will be used. However, when the range
is invalid (i.e `start > end && step > 0` or `start < end && step < 0`), an empty
set will be returned.

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
| `-2\"abcde"` | -2 |  3 |
| `{1: 10, 2: 20, 3: 30}` | 1 | 4 |

For example, on `step < 0`:

| expression | default start | default end |
|:-|:-|:-|
| `[1, 2, 3, 4, 5]` | 4 |  -1 |
| `-2\"abcde"` | 2 |  -3 |
| `{1: 10, 2: 20, 3: 30}` | 3 | 0 |

#### Out of Range Slices

Sometimes the provided slice can contain indexes that are out of range. For
`Array` or `String`, an invalid slice can contain any indexes that are less than
the offset or more than the `offset + length`. For `Dictionary`, invalid slice
can contain any indexes that are less than the smallest key or larger than the
maximum index. If an invalid slice is being evaluated, the operation will **fail**.
The following table will show you the allowed range of values for `start` and
`end`.

| type | allowed start | allowed end |
|:-|:-|:-|
| `Array`/`String` | `offset <= start < length + offset` | `offset - 1 <= end <= length + offset` |
| `Dictionary` | `smallest key <= start <= largest key` | `smallest key - 1 <= end <= largest key + 1` |

#### Dictionary

Slicing in `Dictionary` is quite unique. Slicing in dictionary will always return
an array of values. However, if the slice is out of range (like the above) or
it results in empty values, slicing will return an empty `Set`. Also, since
slicing only works with numerical values (for `start` and `end`), any
non-numerical keys will be ignored.

Another important thing to note is that unlike `Array` or `String` it's possible
for `Dictionary` to have gaps in index (e.g. `keys = [1, 2, 5]`). To handle this
arrai will only collect the values whose keys are contiguous and ignore the rest.

For example

| expression | equals |
|:-|:-|
| `{1: 1, 2: 2, "c": 3, 4: 4}(1;)` | `[1, 2, 4]` |
| `{1: 1, 2: 2, 4: 4}(1;)` | `[1, 2]` |
| `{1: 1, 2: 2, 4: 4}(;)` | `[1, 2]` |
| `{"a": 1, "b": 2, "c": 3, "d": 4}(1;)` | `{}` |

## Inclusivity **(⛔ NYI)**

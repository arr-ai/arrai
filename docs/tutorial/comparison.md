# Comparison operators

Comparison operators compare values with each other and return `true` or `false`
to indicate whether relationship between them holds.

```arrai
@> [1 = 1, 1 != 1, 1 < 1, 1 > 1, 1 <= 1, 1 >= 1]
@> [1 = 2, 1 != 2, 1 < 2, 1 > 2, 1 <= 2, 1 >= 2]
@> [3 = 2, 3 != 2, 3 < 2, 3 > 2, 3 <= 2, 3 >= 2]
```

Note that `false` is the empty set, so it displays as `{}`. Also, while `true`
is actually `{()}` it displays as `true` because it is almost never used to mean
something else.

## Set comparisons (FUTURE)

A related and more extensive set of operators is avaiable for set comparisons.
These operators determine subset and superset relationships rather than strict
less/greater ordering.

```arrai
@> {} (<) {1, 2}
@> {1, 2} (!<) {}
@> {1, 2} (!<=) {}
@> {1, 2} (!<) {1, 2}
@> {1, 2} (<=) {1, 2}
@> {1, 2} (!<>=) {1, 3}
```

The general form is an optional `!` denoting "not" followed by any combination
of `<`, `>` and `=`, denoting "is subset of", "is superset of" and "equals",
respectively.

In addition to comparing sets, you can test set membership:

```arrai
@> 2 <: {1, 2, 3}
@> 4 !<: {1, 2, 3}
@> {2, 3} <: {1, 2, 3, 4}
@> {2, 3} <: {1, {2, 3}, 4}
```

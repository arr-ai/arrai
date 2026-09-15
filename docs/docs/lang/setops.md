---
id: setops
title: Set operators
---

```arrai
@> {'you', 'them', 'and', 'me'} with 'or' without 'you'
```

## `with` and `without`

`set with elem` returns a new set containing all the elements of `set` plus
`elem`. `set without elem` returns a new set containing all the elements of
`set` except `elem` (a no-op if `elem` isn't a member).

```arrai
@> {1, 2, 3} with 4
@> {1, 2, 3} without 2
@> {1, 2, 3} without 5   # elem not present: returns set unchanged
```

Because arrays, strings, byte arrays and dictionaries are just sets of tuples
under the hood, `with`/`without` work on their underlying `(@:, @item:)`,
`(@:, @char:)`, `(@:, @byte:)` and `(@:, @value:)` tuples too:

```arrai
@> [1, 2, 3] without (@: 0, @item: 1)
```

## `count` and `single`

The postfix `count` operator returns the number of elements in a set:

```arrai
@> {1, 2, 3} count
```

The postfix `single` operator returns the sole element of a one-element set,
and fails if the set is empty or has more than one element:

```arrai
@> {42} single
@> {} single       # FAIL: empty set
@> {1, 2} single   # FAIL: too many elements
```

## Boolean set operators

Arr.ai supports the conventional set operators.

```arrai
@> {1, 2, 3} & {2, 3, 4}   # Intersection
@> {1, 2, 3} | {2, 3, 4}   # Union
@> {1, 2, 3} &~ {2, 3, 4}  # Difference
@> {1, 2, 3} ~~ {2, 3, 4}  # Symmetric difference
```

## Power set: `^set`

The power set of a set is the set of all subsets of that set, including the set
itself and the empty set.

```arrai
@> ^{1}
@> ^{1, 2}
@> ^{1, 2, 3}
```

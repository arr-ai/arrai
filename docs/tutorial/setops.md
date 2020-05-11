# Set operators

```arrai
@> {'you', 'them', 'and', 'me'} with 'or' without 'you'
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

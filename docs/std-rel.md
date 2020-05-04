# rel

The `rel` library contains functions for relational operations.

## `//rel.union(s <: array_of_sets) <: set`

`union` takes an array of sets `s` and does a union operation on each member of `s`. It
returns the unioned sets.

| example | equals |
|:-|:-|
| `//rel.union([{1, 2}, {3, 4, 2, 10}, {4, 5, 'another', 'duplicate'}, {'duplicate'}])` | `{1, 2, 3, 4, 5, 10, 'another', 'duplicate'}` |

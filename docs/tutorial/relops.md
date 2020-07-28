# Relational operators

Recall that relations are sets of tuples wherein every tuple has the same
attribute names.

```arrai
@> /set planet = {
    |name,      year,   symbol|
    ("Mercury", 87.97,  "☿"   ),
    ("Venus",   224.70, "♀"   ),
    ("Earth",   365.26, "♁"   ),
}
```

A number of operators is defined specifically to work with relations.

## Join

If R is a relation with attributes x and y, and S is a relation with attributes
y and z:

```arrai
@> /set R = {|x,y| (1, 2), (1, 4), (4, 3)}
@> /set S = {|y,z| (4, 5), (4, 1), (3, 3)}
@> R <&> S
```

### Join variants

The join operator, `<&>` is the basis for a family of related operators. Each
operator performs the same underlying operation of find pairs of tuples whose
matching attributes are equal. They differ only in which attributes appear in
the output relation.

Assuming A is a relation with attributes x and y, and B is a relation with
attributes y and z.

|operator|returns|description|
|:-:|-|-|
| `A <&> B` | x, y, z | Join |
| `A <-> B` | x, z | Compose |
| `A -&- B` | y | Intersection of common attributes |
| `A --- B` | y | true iff join isn't empty |
| `A -&> B` | y, z | All tuples from B with matching tuples in A |
| `A <&- B` | x, y | All tuples from A with matching tuples in B |
| `A --> B` | z | All tuples from B with matching tuples in A, excluding common attributes |
| `A <-- B` | x | All tuples from A with matching tuples in B, excluding common attributes |

The syntax follows a consistent pattern. The base operator is `<&>`, which
returns all attributes. Each variant excludes attributes from the result based
on which characters from the base operator have been replaced with `-`:

- If `<` is replaced, all attributes unique to the left argument are omitted.
- If `&` is replaced, all attributes common to both arguments are omitted.
- If `>` is replaced, all attributes unique to the right argument are omitted.

```arrai
@> /set R = {|x,y| (1, 2), (1, 4), (4, 3)}
@> /set S = {|y,z| (4, 5), (4, 1), (3, 3)}
@> R <-> S
@> R -&- S
@> R --- S
@> R -&> S
@> R <&- S
@> R --> S
@> R <-- S
```

## Rank

The `rank` operator computes the position of each tuple in a relation by a
ranking value and adds its position as a new attribute. This is similar to
`orderby`, except that, instead of returning an array in the order given, the
rank operator returns the ordering information as an attribute in the result.

```arrai
@> planet rank (r: .name)
@> planet rank (r: .year)
```

You can even compute and hold multiple ranks in one go, which is one significant
advantage of ranking over ordering.

```arrai
@> planet rank (alphabetical: .name, speed: 1/.year)
```

If several tuples have the same ranking value, they will have the same rank.

```arrai
@> {|x,y| (1, 0), (1, 1), (2, 1), (3, 1)} rank (r: .x)
```

Note that the rank skips from 0 to 2. This is because the rank of a tuple is
defined as the number of tuples that have a lower ranking value than the given
tuple.

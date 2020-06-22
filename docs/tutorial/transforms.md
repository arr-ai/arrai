# Transforms

Among arr.ai's powerful features, transform operators are perhaps the most
heavily used. They come in four flavors, each operating on a different kind of
input:

- `->` transforms a single value.
- `=>` transforms the members of a set.
- `>>` transforms the entries of a sequence or dictionary.
- `:>` transforms the attributes of a tuple.

The `->` transform operator has the following syntactic structure:

- `lhs -> expr` binds the name `.` to `lhs` and evaluates `expr` in this
  context.
- `lhs -> \pattern expr` binds `pattern` to `lhs` and evaluates `expr` in this
  context.

The other three operators follow the same syntax. However, they bind each of the
elements of the `lhs`, not the entire `lhs`.

If this isn't clear, don't worry, it'll become clear as we work through the
examples below.

## The basics

Let's explore the basic usage of the above operators, one at a time.

### `->`

The arrow operator, `->`, can seem almost trivial:

```arrai
@> 42 -> .
@> 42 -> . + 1
@> "hello, " -> . ++ "world"
```

The alternative form is likewise fairly trivial:

```arrai
@> 42 -> \x x + 1
```

This operator shows its usefulness when `lhs` becomes more complex:

```arrai
@> /set data = (customer: (custid: 123, name: (first: "Franz", last: "Kafka")))
@> data.customer.name -> $`${//str.upper(.last)}, ${.first}`
```

As you can see, `->` offers a very lightweight mechanism to refer to a common
complex expression while allowing a reader to follow a natural left-to-right
order as they seek to understand the intent.

### `=>`

The "set-arrow" operator, `=>`, requires `lhs` to be a set. It applies the
transform to each member of `lhs` and evaluates to a set of the results:

```arrai
@> {2, 4, 6, 8, 10} => . + 1
@> {2, 4, 6, 8, 10} => . % 6
```

Note that, if the transform produces the same output for some of the `lhs`
elements, the resulting set will have fewer members than `lhs`.

### `>>`

The "sequence-arrow" operator, `>>` requires `lhs` to be a binary relation with
`@` as one of its attributes. This includes arrays, dictionaries, strings and
byte arrays. It transforms the value associated with each `@`.

As a refresher, remember that arrays are just sets whose tuples have two
attributes: `@` denoting a position within the array and `@item` denoting the
value at that position. Similarly, strings, are just sets whose tuples have
attributes `@` and `@char`. And likewise for byte arrays (`@` and `@byte`) and
dictionaries (`@` and `@value`).

```arrai
@> {(@:0, @item:1), (@:1, @item:2), (@:2, @item:3), (@:3, @item:4)}
@> {(@:0, @char:115), (@:1, @char:116), (@:2, @char:114), (@:3, @char:105), (@:4, @char:110), (@:5, @char:103)}
@> {(@:0, @byte:98,), (@:1, @byte:121), (@:2, @byte:116), (@:3, @byte:101), (@:4, @byte:115)}
@> {(@:[-37.8, 145], @value:'MEL'), (@:[-33.9, 151.2], @value:'SYD'), (@:[51.5, -0.1], @value:'LCY')}
```

The `>>` operator can be thought of as transforming just the value instead of
the key:

```arrai
@> {(@:0, @item:1), (@:1, @item:2), (@:2, @item:3), (@:3, @item:4)} >> 2 ^ .
@> [1, 2, 3, 4] >> 2 ^ .
@> ["s", "toe", "thumb"] >> . ++ "nail"
@> {"A": 10, "B": 42, "C": 100} >> 100 - .
@> {[-37.8, 145]: 'mel', [-33.9, 151.2]: 'syd', [51.5, -0.1]: 'lcy'} >> //str.upper(.)
```

#### `>>>`

The `>>>` operator is a variant of the `>>` operator. It performs the same basic
operation, but also makes the `@` attribute available to the transformation
expression. This form is special in that it must be a two-argument function.

```arrai
@> ['red', 'green', 'blue'] >>> \i \c $'${c} = ${i}'
@> /set cities = {[-37.8, 145]: 'MEL', [-33.9, 151.2]: 'SYD', [51.5, -0.1]: 'LCY'}
@> cities >>> \[lat, lng] \code $'${code} is at (${lat}, ${lng})'
```

Note that `>>>` still transforms only the associated value. `@` remains
unchanged. If you want to transform the entire tuple, use `=>` instead.

### `:>`

The "tuple-arrow" operator, `:>`, requires `lhs` to be a tuple. It applies the
transform to each attribute value of `lhs` and evaluates to a tuple of the results:

```arrai
@> (r: 0.5, g: 0.2, b: 0.7) => 1 - .
```

## Unary forms

The above operator can be used in unary form. For instance, `=> \x 2 + x`. In
all cases, the `lhs` is implied to be `.`. This property allows the transform
operators to be chained together to transform deeper structures:

```arrai
@> [{1, 2}, {2, 3, 4}, {1, 5}] >> => 10 + .
@> {(r:0.7, g:0, b:0), (r:0.4, g:0.6, b:0), (r:0.5, g:0.5, b:1)} => :> . ^ 2
```

## Interaction with order and orderby

(The following discussion applies equally to `order` and `orderby`, so we'll
only discuss `orderby`.)

The `orderby` operator transforms a set into an array:

```arrai
@> {'red', 'green', 'blue'} orderby .
@> {(r:0.7, g:0, b:0), (r:0.4, g:0.6, b:0), (r:0.5, g:0.5, b:1)} orderby .r
@> {(r:0.7, g:0, b:0), (r:0.4, g:0.6, b:0), (r:0.5, g:0.5, b:1)} orderby .g
```

Things can get a bit confusing when applying `orderby` to a sequence or
dictionary. It is important to remember that `orderby` always operates on the
underlying set, never on the abstraction it represents. For instance, consider
the following expression:

```arrai
[9, 4, 2, 4] orderby .
```

At first glance, it seems that this will reorder the numbers of the array.
However, when you try it out, you'll notice that the output is in fact an array
of the tuples underlying the original array, not an array of the numbers it
represents.

In order to rearrange the values accordingly, you must explicitly deal with the
underlying representation of arrays. While this can be confusing at first, it is
a result of consistently applying the rule that `orderby` operates on sets. It
also gives the programmer full control over how ordering should behave. In the
above case, a question arises: do you want two instances of the number 4 in the
result? There isn't a single correct answer. If you want an array of all scores,
you can do it like this:

```arrai
@> [9, 4, 2, 4] orderby [.@item, .] >> .@item
```

Let's unpack this one step at a time. `orderby [.@item, .]` evaluates to an
array ordered first by `.@item`, which is the number, and then by `.`, which is
the whole tuple. The purpose of the `.` is to act as a tie-breaker. Without it,
the results will behave strangely in future. You could also use `.@` as the
tie-breaker in this case, since `.@` will be unique for all input elements.
However, this assumes regular arrays in which all values of `.@` are unique, and
this is not guaranteed in arr.ai. It's generally safest to use `.` as a
tie-breaker, since it is guaranteed to be unique, regardless of the input set's
structure.

As observed earlier, the result is an array of the tuples underlying the
original array, so we use `>> .@item` to extract just the numbers from that
intermediate result.

On the other hand, if you only want an ordered array of the scores that were
achieved and don't want to see duplicate scores repeated, it's simply a matter
of getting a set of the scores and ordering them:

```arrai
@> [9, 4, 2, 4] => .@item orderby .
```

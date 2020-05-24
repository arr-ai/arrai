# Arr.ai values

Arr.ai has three kinds of values: numbers, tuples and sets.

In the following sections, enter the text appearing after each `@>` and confirm
that the output matches what you expect.

## Numbers

```text
@> 42
42
@> +1.23E4
12300
@> +1.23E40
1.23E40
@> -10
-10
```

**Explore:** What do the following output?

```text
@> 1.23E-4
@> 1.23E-5
@> --10
@> 1/0
@> 0/0
```

## Tuples

Tuples associate names with values. Each name/value pair is called an attribute.

```text
@> ()
()
@> (x: 1, y: 2)
(x: 1, y: 2)
@> (y: 2, x: 1)
(x: 1, y: 2)
@> (start: (x: 1, y: 2), end:(x: -5, y: 3))
(end: (x: -5, y: 3), start: (x: 1, y: 2))
```

Note that tuple names are unordered. That is, `(x: 1, y: 2)` and `(y: 2, x: 1)`
represent the same tuple. However, for convenience and consistency, arr.ai will
usually print tuples with names ordered lexicograpically.

Attribute names can be pretty much anything. Names that don't fit the standard
identifier rules can be written using string syntax:

```text
@> ('-1-': 23, '-2-': 34)
('-1-': 23, '-2-': 34)
@> ('\007': 100)
('\a': 100)
@> ('🐭': "mouse", '🐴': "horse", '🐱': "cat")
('🐭': 'mouse', '🐱': 'cat', '🐴': 'horse')
```

**Explore:** What do the following output?

```text
@> ((a: 1, b: 2))
@> ((()))
```

## Sets

Sets hold values. Any given value is either in a set or it isn't. There is no
notion of multiplicity; a value cannot be present more than once. There is also
no notion of ordering.

```text
@> {}
{}
@> {1, 2, 3}
{1, 2, 3}
@> {{1, 2}, {3, 4}}
{{1, 2}, {3, 4}}
@> {1, {2, 3}}
{1, {2, 3}}
@> {(x: 1), (s: {2, 3})}
{(s: {2, 3}), (x: 1)}
@> {()}
true
@> 3 < 4
true
```

Note that we introduced a Boolean test in that last expression, `3 < 4`. Since
`3` is in fact less than `4`, that expression evaluates to `true`. However,
since there is no Boolean type in arr.ai, `true` is actually an alias for the
set with just the empty tuple in it, `{()}`.

**Explore:** What do the following output?

```text
@> true
@> false
```

Why doesn't `false` print as `false`? As explained above, `true` is actually a
set (the set with the empty tuple, to be precise). Similarly, `false` is
actually the empty set. However, while `{()}` very rarely means anything other
than `true`, the empty set often means something other than `false`, so it is
usually printed in its base form: a set. If you want to output Boolean values as
"true" and "false", you can do this in an [expression string](../expr-str.md).

```arrai
@> let passed = false; $"...${passed:t}..."
```

## Special values

It probably feels like the above type system just isn't enough. In real world
programming, we typically want to work with a range of usefule types, such as
strings, arrays and functions. These and more are in fact available in arr.ai,
but they do not exist as distinct types from the above. Rather, they are
expressed as special compositions of the above types.

### Mapping types

The following sample run demonstrates the syntax for arrays, strings and
dictionaries. It also shows that they are nothing more than syntactic sugar for
special compositions of numbers, tuples and sets.

```text
@> [3, 9, 27]
[3, 9, 27]
@> {(@: 0, @item: 3), (@: 1, @item: 9), (@: 2, @item: 27)}
[3, 9, 27]
@> "hello"
'hello'
@> 'hello'
'hello'
@> `hello`
'hello'
@> {(@: 0, @char: 104), (@: 1, @char: 101), (@: 2, @char: 108), (@: 3, @char: 108), (@: 4, @char: 111)}
'hello'
@> {1: "hi", 2: "bye"}
{1: 'hi', 2: 'bye'}
@> {(@: 1, @value: 'hi'), (@: 2, @value: 'bye')}
{1: 'hi', 2: 'bye'}
@> {[1, 2]: 42, [3, 5]: 54}
{[1, 2]: 42, [3, 5]: 54}
```

**Explore:** How would you describe the following values?

```text
@> {{}: {}, (): ()}
@> {{{(x:1): 42}: 54}}
```

**Food for thought:** Since tuples and dictionaries are both basically key/value
collections, how would you characterise their differences and can you explain
why both are available in arr.ai? The following idiomatic example might help.

```text
@> {"dog": (legs: 4, sound: "bark"), "cat": (legs: 4, sound: "meow")}
```

A good way to think about the distinction is that the attributes of a tuple are
typically known at the time the code is written, while the keys of a dictionary
are not. Tuples represent structures with a known shape, while dictionaries hold
collections keyed by some value.

Dictionaries are actually functions. The keys are the allowable input parameters
and the value for each key is what the function returns when passed that key.

### Sequence types with offsets

Since strings, array, etc., are simply sets of tuples, what happens if they
don't quite fit the pattern of a string? For instance, what if we remove the
tuples representing the first two characters from the string `"treat"`?

```text
@> "treat" where .@ >= 2
2\'eat'
```

The result is a string with an offset. The `2\` in front indicates that the
offset string starts at index 2. You can enter such sequences directly:

```text
@> 2\"eat"
2\'eat'
```

**Explore:** What does the following output?

```text
@> (2*3)\"abc"
```

### Relations

The relational model is core to the design of arr.ai, so it might come as a
surprise that relations are not a first class type within the language. The
reason for this is simple. A relation is simply a set of tuples, all of which
have the same names.

```text
@> {(x: 1, y: 1), (x: 5, y: 1), (x: 4, y: 2), (x: 2, y: 2)}
{(x: 1, y: 1), (x: 2, y: 2), (x: 4, y: 2), (x: 5, y: 1)}
```

Since the concept of a set of tuples with the same names is so common, arr.ai
offers a special syntax for this:

```text
@> {|x,y| (1, 1), (5, 1), (4, 2), (2, 2)}
{(x: 1, y: 1), (x: 2, y: 2), (x: 4, y: 2), (x: 5, y: 1)}
```

You may have noticed that strings, arrays and dictionaries are in fact
relations.

```text
@> {|@,@item| (0, 3), (1, 9), (2, 27)}
[3, 9, 27]
@> {|@,@char| (0, 104), (1, 101), (2, 108), (3, 108), (4, 111)}
'hello'
@> {|@,@value| (1, 'hi'), (2, 'bye')}
{1: 'hi', 2: 'bye'}
```

**Explore:** What does the following output?

```text
@> {|x| (1), (2)}
```

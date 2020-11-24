---
id: data
title: Data
---

We start with data, because:

> *Bad programmers worry about the code. Good programmers worry about data structures and their relationships.* &mdash;
> [Linus Torvalds](https://www.goodreads.com/quotes/1188397-bad-programmers-worry-about-the-code-good-programmers-worry-about)

Arr.ai's data model is remarkably simple, having only three kinds of values, all immutable:

1. **Numbers** are 64-bit binary floats.
2. **Tuples** associate names with values.
3. **Sets** hold sets of values.

Let's be clear about what the above means. Arr.ai has no arrays. It also has no strings, Booleans, maps, functions, packages, pointers, structs, classes or streams. Arr.ai has numbers, tuples and sets. There is nothing else.

But let's also be clear that this is far less restrictive than it might at first seem. You can in fact represent:

1. Arrays: `[]`, `[2, 4, 8]`
2. Strings: `""`, `"hello"`
3. Booleans: `true`, `false`
4. Maps: `{}`, `{"a": 42}`, `{1: 34, 2: 45, 3: 56}`
5. Functions:
   1. Functions are unary: `\x 1 / x`
   2. Binary functions don't exist, but `\x \y //math.sqrt(x^2 + y^2)` is a unary function that takes a single parameter, `x`, and returns a unary function. The returned function takes a single parameter, `y`, and returns the hypotenuse of a right triangle with sides *x* and *y*.
6. Packages:
   1. `//math.sin(1)`
   2. `//{./myutil/work}(42)`
   3. `//{/path/to/root/file}`
   4. `//{./myfile.yaml}`
   5. `//{github.com/org/external/file}`
   6. `//{https://url/to/your/content}`

All of the above forms are syntactic sugar for specific combinations of numbers, tuples and sets. For example, the string `"hello"` is a shorthand for the following set:

```arrai
{
   (@: 1, @char: 101),
   (@: 2, @char: 108),
   (@: 4, @char: 111),
   (@: 3, @char: 108),
   (@: 0, @char: 104),
}
```

(Order doesn't matter in a set. It's the `@` attribute that determines the position of each character in the string being represented.)

## Data transformation

Arr.ai is an expression language, which means that every arr.ai program, no matter how complex, is a single expression evaluating to a single value. You can play with the language on the command line by running `arrai eval <expression>`, with `e` being a shortcut for `eval` (see [here](../cli/eval) for a detailed description of the `eval` command), e.g.:

```bash
$ arrai e 42
42
$ arrai e '//math.pi'
3.141592653589793
$ arrai e '[1, (a: 2), {3, 4, 5}]'
[1, (a: 2), {3, 4, 5}]
$ arrai e '[1, (a: 2), {3, 4, 5}](1)'  # Arrays are functions.
(a: 2)
$ arrai e '"hello"(3)'                 # So are strings.
108
$ arrai e '"hello" => (@:.@, @item:.@char)'
[104, 101, 108, 108, 111]
$ arrai e '[104, 101, 108, 108, 111] => (@:.@, @char:.@item)'
hello
$ arrai e '{
   (@: 1, @char: 101),
   (@: 3, @char: 108),
   (@: 0, @char: 104),
   (@: 4, @char: 111),
   (@: 2, @char: 108),
}'
hello
```

The last example underscores the point made earlier that strings are in fact sets of tuples. There is no semantic distinction between the two forms.

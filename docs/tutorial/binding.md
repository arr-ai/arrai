# Name bindings

Aside from being able to express literal values, one of the most basic needs
almost any programming language must satisfy is the ability to name values so
they may be referred to later in the code.

Names are sequences of one or more ASCII letters and digits as well as `_`, `$`
and `@`. They may not begin with a digit. The single character, `.`,  is also a
valid name by itself. It has special properties that we'll explore below.

There are only two ways to bind names in arr.ai: let-bindings and function
parameters.

## Let-bindings

Let-bindings exist for two main reasons. The first is to allocate a name to an
expression that will be used multiple times in subsequent code. The second is to
give a meaningful name to an expression whose purpose might need clarification.

Try out the following examples to see the first main reason for let-bindings:

```arrai
@> let x = //math.pi; 3*x^2 + 2*x - x
```

```arrai
@> let customer = (name: "Anders", age: 47);
 > let age = customer.age;
 > cond (
 >     age < 0: "wat?",
 >     0 <= age < 40: "young",
 >     40 <= age < 60: "middle",
 >     *: "old"
 > )
```

In the following example, `velocity` illustrates the first main reason, avoiding
repetition, while `speed` illustrates the second, clarification:

```arrai
@> let v = (x: 0.4, y: 0.7, z: 0.6);
 > let speed = (v.x^2 + v.y^2 + v.z^2)^0.5;
 > cond (
 >     speed < 1: "too slow",
 >     *: "fast enough"
 > )
```

Try a few variations of allowable characters:

```arrai
@> let $x = 10;
 > let foo@1 = 11;
 > let __z__ = 123;
 > $x + foo@1 * __z__
```

### Important caveats

1. Names prefixed with `@` are reserved for special purposes as defined by the
   language. Arr.ai will not prevent you from using them as regular names, but
   you might find that they cause intended problems in your code. Never use
   `@`-prefixed names for anything but the intended usage as documented by
   arr.ai. This also applies to names that don't currently have a defined use.
   In future, arr.ai may explicitly disallow `@`-prefixed names that aren't in a
   whitelist of names with a defined purpose.
2. Names prefixed with `$` are currently not special, but might soon be. Avoid.

### The special name: `.`

The name `.` can be used as a regular name, with some qualifiers:

```arrai
@> /set cell = (row: 1, col: 1, value: 0.5)
@> cell.row   # OK
@> let . = cell; ..row    # FAIL
@> let . = cell; (.).row  # OK
@> let . = cell; .row     # OK
```

As you can see, `.` can be used as an "implied" value with the `.` operator.
However, explicitly assigning `.` is not generally recommended. Below, we'll
explore how this special name is normally used.

## Function parameters

Functions are defined to take some unknown value and evaluate an expression that
incorporates that value. The parameter name of a function constitutes the second
way to bind names. In this case, the reason for the name (`v` in the example
below) is necessity. It is how the function refers to the argument passed in.

```arrai
@> let length = \v (v.x^2 + v.y^2 + v.z^2)^0.5;
 > length((x: 0.4, y: 0.7, z: 0.6))
```

In the above example, the function `length` is passed a vector, whose value is
bound to the name `v` inside the function's body.

## Transform operators

There is in fact a third way to bind names. Technically, it is just a variant of
function parameter syntax, though it might not seem like it at first.

Arr.ai offers a family of generalised transform operators: `->`, `=>`, `>>` and
`:>`. The intent here is not to go into the details of how these operators work.
We will explore these in a subsequent tutorial. For now, we'll just focus on
`->` and `=>`. The other two work in basically the same way when it comes to
name binding.

The `->` operator is an alternative syntax for function calls.

```arrai
@> (\s //str.upper(s))("hello")
@> "hello" -> \s //str.upper(s)
```

The above example is somewhat contrived, but it illustrates the key point that
`->` is just a kind of function call, with the only difference being that the
roles are reversed, with the parameter appearing on the left and the function on
the right.

Since `.` is a valid name, you could also use it as the parameter name.

```arrai
@> "hello" -> \. //str.upper(.)
```

Again, this example is somewhat contrived, so here's a different one that makes
the potential benefits a bit clearer:

```arrai
@> /set customer = (
 >     name: "Alice",
 >     dob: "1990-01-01",
 >     tfn: "123123123123",
 >     address: "12 Station St, Melbourne",
 >     marital_status: "M",
 >     salary: 100000,
 > )
@> (:customer.name, :customer.dob, :customer.address)
@> customer -> \. (:.name, :.dob, :.address)
```

Because this is such a common pattern, arr.ai lets you omit the `\.`:

```arrai
@> customer -> (:.name, :.dob, :.address)
```

The `=>` operator performs a similar function over sets of things:

```arrai
@> {"dog", "duck", "chicken", "cat", "giraffe"} => //str.upper(.)
```

In this case, `.` is bound to each element of the set in turn and each result
becomes part of the output.

Remember that, just as for `->`, the above is still shorthand for reverse
function-call syntax, as the following examples illustrate.

```arrai
@> {"dog", "duck", "chicken", "cat", "giraffe"} => \. //str.upper(.)
@> {"dog", "duck", "chicken", "cat", "giraffe"} => \animal //str.upper(animal)
```

As an aside, because `=>` operates on sets, you'll find that the result is not
necessarily displayed in the same order as the input. In fact, it may not even
have the same number of elements:

```arrai
@> {-4, -3, -2, -1, 0, 1, 2, 3, 4} => . ^ 2
```

We'll learn more about `=>` in a later tutorial.

## Pattern matching

Both let-bindings and function parameters support pattern matching. This is a
very powerful mechanism to extract elements from a complex structure and also
restrict what values may be used as input.

Bare literals are supported and will simply match their own value:

```arrai
@> let 42 = 42; 1
@> let "hello" = "hello"; 1
@> let 3 = 1 + 2; 5
@> let true = true; 3
@> let true = {()}; 3
```

Try out the following examples to use pattern with arrays:

```arrai
@> let [] = []; 1
@> let [a, b, c] = [1, 2, 3]; b
@> let arr = [1, 2]; let [a, b] = arr; b
@> let [x, x] = [1, 1]; x
@> let [x, x] = [1, 2]; x                  # should fail
@> [1, 2] -> \[x, y] x + y
@> let f = \[x, y] x + y; f([1, 2])
@> (\[x, y] x + y)([1, 2])
@> (\z \[x, y] z/(x + y))(9, [1, 2])
```

with tuples:
```arrai
@> let () = (); 1
@> let (a: x, b: y) = (a: 4, b: 7); x
@> let (a: x, b: x) = (a: 4, b: 4); x
@> let (:x) = (x: 1); x
@> (m: 1, n: 2) -> \(m: x, n: y) x + y
```

with dictionaries:
```arrai
@> let {"a": f, "b": k} = {"a": 1, "b": 2}; [f, k]
@> {"m": 1, "n": 2} -> \{"m": x, "n": y} x + y
```

and with sets:
```arrai
@> let {} = {}; 1
@> let {a, 42} = {3, 42}; a
@> let {a, b} = {3, 42}; [a, b]        # should fail because it is a non-deterministic situation
```

Also, nested patterns are supported as:
```arrai
@> let [[x, y], z] = [[1, 2], 3]; x
@> let [{"a": x}, (b: y), z] = [{"a": 1}, (b: 2), 3]; [x, y, z]
@> [1, [2, 3]] -> \[x, [y, z]] x + y + z
```

Underscore `_` matches any value and ignores it.

```arrai
@> let [x, _, _] = [1, 2, 3]; x
@> let [_, x, _] = [1, 2, 3]; x
```

A name within parentheses like `(x)` refers to the value bound to the name `x`.

```arrai
@> let x = 3; let [b, x] = [2, 4]; x
@> let x = 3; let [b, (x)] = [2, 3]; b
@> let x = 3; let [_, b, (x)] = [1, 2, 3]; b
@> let x = 1; [1, 2] -> \[(x), y] y
@> let x = 1; let y = 42; let {(x), (y)} = {42, 1}; 5
@> let x = 3; let [b, (x)] = [2, 4]; b     # should fail because (x) != 4
@> let [(x)] = [2]; x                      # should fail because `x` isn't in scope
```

`let a = 56; let {"x": a, "y": (a)} = {"x": 42, "y": 56}; a` is valid 
but `let a = 56; let {"x": a, "y": (a)} = {"x": 42, "y": 42}; a` should fail. 
Using the same name in an expression, `(a)`, and a newly bound name, `a`, is
confusing and should be avoided.

Complex expressions are supported and will also match their own value. However, they must be enclosed in parentheses, `(...)`:
```arrai
@> let (1 + 2) = 3; 5
```

What's more, arr.ai allows extra elements `...` or `...x` in addition to 
the explicitly bound ones and binds name `x` to any additional elements 
that weren't explicitly matched by other patterns.
```arrai
@> let [x, y, ...] = [1, 2]; [x, y]
@> let [x, y, ...t] = [1, 2]; [x, y, t]
@> let [x, y, ...] = [1, 2, 3, 4, 5, 6]; [x, y]
@> let [x, y, ...t] = [1, 2, 3, 4, 5, 6]; [x, y, t]
@> let [..., x, y] = [1, 2, 3, 4, 5, 6]; [x, y]
@> let [...t, x, y] = [1, 2, 3, 4, 5, 6]; [x, y, t]
@> let [x, ..., y] = [1, 2, 3, 4, 5, 6]; [x, y]
@> let [x, ...t, y] = [1, 2, 3, 4, 5, 6]; [x, y, t]
@> let (m: x, n: y, ...t) = (m: 1, n: 2, j: 3, k: 4); [x, y, t]
@> let {"m": x, "n": y, ...t} = {"m": 1, "n": 2, "j": 3, "k": 4}; [x, y, t]
@> let {1, 2, 3, ...t} = {1, 2, 3, 42, 43}; t
@> let x = 1; let y = 42; let {(x), (y), ...t} = {1, 42, 5, 6}; t
@> [1, 2, 3, 4] -> \[x, y, ...t] [x + y, t]
```

(Conditional Accessor Syntax)[../README.md] is also supported in pattern matching:
```arrai
@> let {"a"?: x:42} = {"a": 1}; x = 1
@> let {"b"?: x:42} = {"a": 1}; x = 42
@> let (b?: x:42) = (a: 1); x = 42
@> let [x, y, z?:0] = [1, 2]; [x, y, z] = [1, 2, 0]
@> let {"b"?: x:42, ...t} = {"a": 1}; [x, t] = [42, {"a": 1}]
```

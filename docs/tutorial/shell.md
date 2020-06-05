# Arr.ai shell basics

This tutorial will provide a brief walkthrough of the features of arr.ai's
interactive shell.

## Starting the shell

There are two ways to start the interactive shell. The first is to use the `arrai i`
command.

```bash
$ arrai i
@> _
```

Alternately, if the `ai` shortcut is set up correctly (see the main
[README](../../README.md) for instructions), you can use it instead.

```bash
$ ai
@> _
```

Once you see the `@>` prompt, you can start entering expressions or special
commands as shown below.

When you are finished using the arr.ai shell, type the exit character to return
to the OS shell. In Unix-like environments, press Ctrl+D. The output will be as
follows.

```arrai
@> exit
$ _
```

In Windows, TODO.

## Shell as debugger

When you evaluate an arrai script and the script fails, the `arrai` program will
drop into the `arrai` interactive shell with the scope near the point of failure
available to the interactive shell as a tuple.

The values in the scope can be accessed through the variables with their
corresponding names.

```bash
$ arrai e 'let x = 1; (\a a + b)(x)'
INFO[0000] Name "b" not found in {x, a}

.:1:20:
let x = 1; (\a a + b)(x)

.:1:16:
let x = 1; (\a a + b)(x)

.:1:13:
let x = 1; (\a a + b)(x)
@> x
1
@> a
1
@> <tab>
a x
```

## Evaluating expressions

To evaluate an expression in the shell, simply type it in and press enter.

```arrai
@> "Hello, world!"
'Hello, world!'
@> 6 * 9
54
@> //math.pi
3.141592653589793
```

A complex expression may be entered over multiple lines. Arr.ai automatically
detects when an expression isn't complete. (Don't worry if you don't understand
the following expressions. They will make sense later in the tutorial.)

```arrai
@> (function1(a, b, c)
 >  + function2(d, e, f)
 >  + function3(x, y, z)
 > ) / 3
```

Observe that the prompt changes from `@>` to `>` to indicate that it's awaiting
further input. Here's a more complex example.

```arrai
@> let v = (x: 1, y: 2, z: 3);
 > let length = (v.x^2 + v.y^2 + v.z^2)^0.5;
 > cond (
 >     length > 1: "too big",
 >     length < 1: "too small",
 >     *: "just right"
 > )
```

Caution: The current approach to detection uses some simple heuristics such as
counting balanced bracketing. It will not always accurately detect that an
expression is incomplete. The following example will fail.

```arrai
@> 1 /
```

## Global variables

In order to set a global name, use the `/set` command as follows:

```arrai
@> /set x = 42
42
@> /set y = 54
54
@> x / y
0.7777777777777778
```

You may also remove names from the global namespace with `/unset`.

```arrai
@> x / y
0.7777777777777778
@> /unset x
@> x / y # FAIL
@> 1 / y
0.018518518518518517
```

You may have noticed above that arr.ai has something that looks like an assignment
statement: `let NAME = EXPR;`. You may then be wondering why we even need
`/set`.

This form is known as a let-binding and it is *not* in fact a statement. The
full form is actually `let NAME = EXPR1; EXPR2`. What it does is associate
the name `NAME` with the value of `EXPR1` and evaluate `EXPR2` with this
name-binding in effect. For example, in the following expression:

```arrai
@> let x = //math.e; x^2 + 2*x + x
```

the name `x` is bound to the transcendental constant, *e*, when evaluating
`x^2 + 2*x + x`. What this means is that `x` is only in scope for the expression
immediately following the `;`. After the expression is evaluated, `x` is not in
the global namespace. In fact, it never was:

```arrai
@> let x = 42;
 > let y = 54;
 > x / y
@> x / y  # FAIL
@> 1 / y  # FAIL
```

## Exiting shell

The `/exit` command exits the interactive shell. Alternatively, you can press
Ctrl+D.

Usage:

```bash
@> /exit
```

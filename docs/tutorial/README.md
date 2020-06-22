# Arr.ai tutorial

This tutorial will take you step by step through the arr.ai language. It will be
light on theory, focusing more on learning by doing. If you want a more
theoretical introduction to the features of the language, [Introduction to
Arr.ai](../README.md) is a good place to start.

The following chapters start with the basics of the language, and develop
steadily towards the more advanced features.

To follow along, use arr.ai's interactive shell, which may be started with:

```bash
$ arrai i
@> _
```

or, if an appropriate symlink is set up (see the main [README](../../README.md)
for instructions):

```bash
$ ai
@> _
```

Once you see the `@>` prompt, you can try the code examples in the following
chapters.

0. [Arr.ai shell basics](shell.md)
1. [Values](values.md)
2. [Name bindings](binding.md)
3. [Comparison operators](comparison.md)
4. [Arithmetic and logical operators](arithmetic.md)
5. [Set operators](setops.md)
6. [Relational operators](relops.md)
7. (TODO) [Functions and function calls](function.md)
8. (TODO) [Transforms](transforms.md)
9. (TODO) [Standard library](stdlib.md)
10. (TODO) [Importing external code and data](import.md)
11. [Transforms](transforms.md)
12. [Writing tests](testing.md)
13. (TODO) [Expression strings](exprstr.md)
14. (TODO) [Templating with expression strings](templating.md)
15. (TODO) [Grammars](grammars.md)
16. [Macros](macros.md)

## Arr.ai command line

In addition to the interactive shell, arr.ai provides a range of additional
commands to run programs, start an arr.ai server or client, etc.

1. (TODO) run, r: evaluate an arrai file
2. (TODO) eval, e: evaluate an arrai expression
3. (TODO) observe, o: observe an expression on a server
4. (TODO) serve, s: start arrai as a gRPC server
5. (TODO) sync, s: sync local files to a server
6. (TODO) transform, x: transform a stream of input data with an expression
7. (TODO) update, u: update a server with an expression

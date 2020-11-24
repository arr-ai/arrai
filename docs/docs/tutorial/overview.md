---
id: overview
title: Overview
---

:::info
The previous tutorial has been merged with the language reference, so may be somewhat disconnected. A new tutorial is being developed in a series of smaller chunks. It will appear in this section soon.
:::

This tutorial will take you step by step through the arr.ai language. It will be light on theory, focusing more on learning by doing. If you want a more theoretical introduction to the features of the language, [Introduction to Arr.ai](../lang/intro) is a good place to start.

The following chapters start with the basics of the language, and develop steadily towards the more advanced features.

:::info
Arr.ai is designed to be maximally _expressive_, not necessarily familiar or easy to pick up. The payoff for climbing the learning curve is a powerful tool that might just change how you think about programming. Hang in there!
:::

To follow along, use arr.ai's interactive shell, which may be started with:

```bash
$ arrai i
@> _
```

or, if an appropriate symlink is set up (see the [installation instructions](../install)):

```bash
$ ai
@> _
```

Once you see the `@>` prompt, you can try the code examples in the following
chapters.

## Language Features

1. [Arr.ai shell basics](../lang/shell)
1. [Values](../lang/values)
1. [Name bindings](../lang/binding)
1. [Comparison operators](../lang/comparison)
1. [Arithmetic and logical operators](../lang/arithmetic)
1. [Set operators](../lang/setops)
1. [Relational operators](../lang/relops)
1. (TODO) [Functions and function calls](../lang/function)
1. [Transforms](../lang/transforms)
1. (TODO) [Standard library](../lang/stdlib)
1. (TODO) [Importing external code and data](../lang/import)
1. [Writing tests](../lang/testing)
1. [Expression strings](../lang/exprstr)
1. (TODO) [Templating with expression strings](../lang/templating)
1. (TODO) [Grammars](../lang/grammars)
1. [Macros](../lang/macros)

## Arr.ai command line

In addition to the interactive shell, arr.ai provides additional commands to run programs, start an arr.ai server or client, etc.

1. [run](../cli/run), r: evaluate an arr.ai file
1. [eval](../cli/eval), e: evaluate an arr.ai expression
1. (TODO) observe, o: observe an expression on a server
1. (TODO) serve, s: start arrai as a gRPC server
1. (TODO) sync, s: sync local files to a server
1. (TODO) transform, x: transform a stream of input data with an expression
1. (TODO) update, u: update a server with an expression

# Macros

Macros are a powerful tool for data representation and transformation. They are
very similar to functions, except that they operate at parsing time rather than
runtime. This means they can be used to introduce and embed new syntax within
 arr.ai code.

## Usage

> **(⛔ NYI)**: Current syntax is `{:macro:content:}`.

Macros are invoked via the syntax `macro{content}`. The `content` inside the
macro invocation is subject to a grammar defined by the macro itself, not
regular arr.ai syntax. Each macro can support its own grammar for the kind of
content it supports.

The `macro` part of the invocation is either a [grammar](grammars.md) AST
itself, or a tuple with an `@grammar` key corresponding to the grammar AST. That
grammar will then be used to parse the `content` into an AST.

In the tuple form, the macro may also contain a `@transform` key corresponding
to a tuple of transform functions, keyed by the rules of the grammar to which
they apply. When the macro is invoked, `content` is parsed into an AST, which is
then passed to the transform function. The output of the transform will replace
the entire macro invocation before parsing continues.

## Examples **(⛔ NYI)**

The following example expresses a URL as a strongly typed value:

```bash
$ arrai e '//web.url{https://me@foo.com/bar?x=42}'
(
   source: "https://me@foo.com/bar?x=42",
   scheme: "https",
   authority: (
      userinfo: [8\"me"],
      host: 11\"foo.com",
   ),
   path: [19\"bar"],
   search: {23\"x": [25\"42"]},
)
```

This example is representing JSON:

```bash
$ arrai e '//encoding.json{{"x": 1, "y": [2, 3], "z": null}}'
{
   "x": 1,
   "y": [2, 3],
   "z": (),
}
```

(Arr.ai has no counterpart for JSON null, so it uses the empty tuple as a
proxy.)

Using ωBNF to define a custom grammar, you can easily extend arr.ai's syntax.
This example defines a simple grammar and transform for processing dates in
YYYY-MM-DD format, allowing them to be used in arr.ai source code.

```bash
$ arrai e 'let time = (
  @grammar: //grammar.lang.wbnf{default -> y=\d{4} "-" m=\d{2} "-" d=\d{2};},
  @transform: (default: \ast ast -> (year: .y, month: .m, day: .d) :> //eval.value(.''))
);
time{2020-06-09}'
(year: 2020, month: 6, day: 9)
```

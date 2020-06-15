# Macros

Macros are a powerful tool for data representation and transformation. They are
very similar to functions, except that they operate at parsing time rather than
runtime. This means they can be used to introduce and embed new syntax within
 arr.ai code.

## Usage

> **(⛔ NYI)**: Current syntax is `{:macro[rule]?:content:}`.

Macros are invoked via the syntax `macro[rule]?{content}`. In reverse order:

* `content` inside the macro invocation is subject to a grammar defined by the
  macro itself, not regular arr.ai syntax.
* `[rule]`, if provided, specifies which rule of the grammar to use as the root.
  If not provided, the parser will use the first rule declared in the grammar.
* `macro` is either a [grammar](grammars.md) AST, or a tuple with an `@grammar` 
  key corresponding to a grammar AST. Thus each macro can define its own 
  grammar with which to parse its `content`.
  
  In the tuple form, `macro` may also contain a `@transform` key corresponding
  to a tuple. In that tuple, each rule of the grammar corresponds to a function
  that will be applied to the AST produced by parsing `content` with `rule` of 
  the macro's grammar.
  
  If no `@transform` tuple is provided, the macro will simply produce an AST of
  the `content`.
  
The macro invocation will be evaluated during the parsing of the arr.ai program,
and its output will effectively replace the invocation in the code before
parsing continues.

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
  @grammar: //grammar.lang.wbnf{date -> y=\d{4} "-" m=\d{2} "-" d=\d{2};},
  @transform: (date: \ast ast -> (year: .y, month: .m, day: .d) :> //eval.value(.''))
);
time{2020-06-09}'
(year: 2020, month: 6, day: 9)
```

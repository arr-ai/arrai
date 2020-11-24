---
id: grammars
title: Grammars
---

Arr.ai supports encoding of grammars directly inside the language. These
grammars may then be used to parse other content.

Example:

```bash
$ arrai e '//grammar.lang.wbnf{expr -> @:[-+] > @:[/*] > \d+;}{1+2*3}'
("": [+], @rule: expr, expr: [(expr: [("": 1)]), ("": [*], expr: [("": 2), ("": 3)])])
```

(Above syntax **⛔ NYI**. Current syntax is
`{://grammar.lang.wbnf: expr -> @:[-+] > @:[/*] > \d+; :} -> {:.:1+2*3:}`.)

The primary use of grammars is in the macro system. However, grammars are
themselves data structures, and can be transformed as such, allowing interesting
additions such as compositing, subsetting and otherwise transforming grammars.

See also [macros](./macros), which add transformations to parsed syntax.

# grammar

`arrai` has the ability to parse a grammar and use that grammar to parse
inputs as part of its standard library. This package provides that functionality
and supplies built-in grammars to use.

## `parse(grammar <: string|tuple, rule <: string, source <: string) <: tuple`

Takes a tuple representing a `grammar`, the `rule` of the grammar to parse with, and the `source` to parse. Returns an AST of the parsed source.

Examples of grammar tuples include:

- `//grammar.lang.wbnf`: The [ωBNF](https://github.com/arr-ai/wbnf) grammar.
- `//grammar.lang.arrai`: The arr.ai grammar.

Use the `grammar` rule of the `//grammar.lang.wbnf` grammar to create new ωBNF grammars:

```bash
@> //grammar.parse(//grammar.lang.wbnf, "grammar", "a -> '1';")
(@rule: 'grammar', stmt: [(@choice: [1], prod: ('': [2\'->', 8\';'], IDENT: ('': 'a'), term: [(term: [(term: [(term: [(named: (atom: (@choice: [1], STR: ('': 5\"'1'"))))])])])]))])
```

Then use those grammars to parse code:

```bash
@> let g = //grammar.parse(//grammar.lang.wbnf, "grammar", "a -> '1';");
 > //grammar.parse(g, 'a', '1')
('': '1')
```

Read about [Macros](tutorial/macros.md) for more applications of grammars.

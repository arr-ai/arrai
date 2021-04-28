---
id: exprs
title: Expressions
---

Arr.ai expressions are combinations of syntax which can be evaluated to produce some value. All arr.ai programs are expressions, and every [value](./values.md) and operation in a program is an expression.

## Implementation

Each expression is a type defined in the `rel` package, and implements the `Expr` interface:

```go
type Expr interface {
    // All exprs can be serialized to strings with the String() method.
	fmt.Stringer

	// Eval evaluates the expr in a given scope.
	Eval(ctx context.Context, local Scope) (Value, error)

    // Source returns the Scanner that locates the expression in a source file.
	Source() parser.Scanner
}
```

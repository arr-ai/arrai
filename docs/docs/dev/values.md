---
id: values
title: Values
---

Arr.ai values numbers, tuples, sets, and various other structures built on top of them. Values are [expressions](./exprs.md) that simply evaluate to themselves.

## Implementation

Each value is a type defined in the `rel` package, and implements the `Value` interface (which extends the `Expr` interface).

Constructing new Value instances from Go values can be done with `rel.NewValue(interface{})`.

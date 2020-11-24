---
id: testing
title: Writing tests
---

Arr.ai is a functional language designed for the representation and transformation of data. As such, testing arr.ai code is different from testing more stateful, imperative code.

Arr.ai's approach to testing is for test files to produce a data structure within which all leaves are `true` (not just truthy). If any leaf is not `true`, the test is said to have failed.

## TODO

- Provide examples.
- Use macros to augment standard arr.ai code for testing (e.g. replace leaf comparison exprs with equivalent `//testing.assert.*` functions).

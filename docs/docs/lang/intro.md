---
id: intro
title: Introduction
---

Arr.ai is many things, but it is first and foremost a data representation and transformation language. This section provides a complete reference of its syntax and features.

If you're new to arr.ai, you may want to start with [Getting Started](../) and [Tutorial](../tutorial/overview) sections.

## About the name

The domain name arr.ai was available and there was some irony in the fact that a
language called arr.ai doesn't have arrays (though it kind of does; see below).

## Some lexical conventions

Arr.ai has a rich syntax, which we won't dive into just yet. A few elements are
worth covering upfront to aid comprehension below.

1. **Identifiers:** Parameter and variables names start with `_`, `@`, `$` or a
   Unicode letter, and may continue with a sequence of any of these and Unicode
   decimal numbers.

   Examples: `x`, `$y`, `Username`, `i0`, `@j12`, `apple_π`

   The identifier `.` is a special case. It is often used as a default argument
   in transform expressions.

2. **Keywords:** The following names are predefined and cannot be reassigned as
   parameter or variable names: `true`, `false`, `let`

3. **Comments:** Comments start with a `#` and end at the end of the line.

   Example: `# Comment on comments.`

4. **Offset collections:** In the string `"hello"`, the first character, `h`, is
   at position zero. In the alternate form `12\"hello"`, the `h` is at position
   12 and the remaining characters occupy positions 13&ndash;16. This is known
   as an offset-string. Likewise, `5\[1, 2, 3]` represents an offset array.

## 

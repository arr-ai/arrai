---
id: concepts
title: Concepts
---

This document provides an overview of some of the key ideas in arr.ai. It
strives to steer clear of syntax and semantics as much as practicable in order
to focus on the concepts themselves.

## Type system

A core premise of arr.ai is that the type system should be made as simple as
possible, but not simpler.

Numbers, tuples and sets are sufficiently powerful to represent a rich and
diverse set of information structures, including arrays, string, dictionaries
and even functions.

Key to understanding the type system is to realise that, although these other
"types" are visible as distinct things in arr.ai syntax (e.g., the string
`"abc"`), they do not exist as distinct types in the language. The string
`"abc"` is exactly the set `{(@: 0, @char: 97), (@: 1, @char: 98), (@: 2, @char:
99)}`. There is no distinction whatsoever between the two. The string syntax and
the set syntax may be used interchangeably in any expression and the result is
guaranteed to be the same. There isn't even a standard library function
available to ask which representation is being used to implement a given value.
In a very strong sense, strings do not exist in arr.ai. They are just syntactic
sugar for a special arrangement of numbers, tuples and sets.

> **A bit of history:** The original design of arr.ai's type system had numbers,
> tuples, Booleans and relations. After tossing it over for some time, it became
> clear that sets were not only a very important concept, but also more
> fundamental than relations. After much deliberation to try and model sets as a
> special kind of relation (`{|@|, (1), (2), (3)}` in current arr.ai syntax), it
> was decided to reverse the situation and define relations in terms of sets and
> tuples.
>
> Booleans also went the way of the Dodo due to the ability to model them as
> special sets. The use of the empty set, `{}`, for `false` and the set with the
> empty tuple, `{()}`, for `true` has a precedent. In relational theory, two
> special relations model `false` and `true`, respectively: `DUM`, the empty
> table with no attributes, and `DEE`, the table with no attributes containing
> the tuple with no attributes. There has been some consideration towards
> reinstating Booleans as a core type, due to usability concerns (e.g., `true`
> prints as `{()}`). However, it hasn't been a major issue to date, so it will
> not likely happen. One possible solution is to render `{()}` as `true`, since
> that is just about the only use this special set has (though the same doesn't
> apply to `{}`).

<!--
The following is aspirational. Uncomment when it's true.

## Hermetic runtime

Arr.ai's runtime environment is designed to be hermetic. This means that access
to the outside world is constrained through opt-in mechanisms. For example, file
access is not automatically available to arr.ai programs, but must instead be
explicitly granted through configuration mechanisms and, even when so granted,
might be restricted to a subset of the filesystem. Local imports are restricted
to the containing module. Remote imports are restricted to a whitelist, which is
empty by default.
-->

## Immutable values

All arr.ai values are immutable. Each arr.ai program is a single expression that
defines some output as a function of its inputs. Expressions are pure, that is,
they have no side effects. Even writing to a filesystem is impossible, which is
why arr.ai directly supports outputting zip and tar archives.

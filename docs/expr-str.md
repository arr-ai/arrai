# Expression strings

Expressions strings turn arr.ai into a sophisticated text templating engine.
Conceptually, they are a simple extension to strings, allowing expressions to be
nested inside such strings. In practice, this opens up a rich text formatting
system that allows production of very complex output structures, including the
arbitrarily deep structures required for code generation.

## Basic structure

Expression strings are like regular strings with three differences:

1. Expression strings begin with a `$`. The regular string `"abc"` is equal to
   the expression string `$"abc"`. The same applies to single-quoted and
   backquoted strings.
2. Expression strings treat whitespace differently than regular strings do.
3. Expressions may be embedded within an expression string, allowing for dynamic
   content. The expression string `$"id = :{i}:;"` evaluates to `"id = 42;"` if
   `i` equals 42.

## Whitespace rules

Expression strings apply the following rules to handle whitespace:

1. If the first character of a string is a newline, it is discarded. This only
   applies for literal newlines in the source. The `\n` escaped form will be
   retained.

   **Example:** The following string is equal to `"abc"`.

   ```text
   $"
   abc"
   ```

2. After removal of any newlines per the first rule, any leading whitespace up
   to the first newline or non-whitespace character will be removed. Subsequent
   occurrences of the same leading whitespace after a newline are also removed.

   **Example:** The following:

   ```text
   $"
       abc
         def
       ghi
   "
   ```

   produces the following output:

   ```text
   abc
     def
   ghi
   ```

3. If, after indentation removal, an embedded expression is the only content on
   a line and the formatted result is empty, then the entire line is omitted
   from the final output.

   **Example:** The following expression equals `"abc\ndef"`:

   ```text
   let s = "" in $"
       abc
       :{s}:
       def"
   ```

   The following expression, in contrast equals `"abc\n123\ndef"` (note the
   extra `\n` in the result):

   ```text
   let s = "123" in $"
       abc
       :{s}:
       def"
   ```

## Embedded expressions

Embedded expressions are evaluated and formatted to provide content for the
expression strings containing them. Their general form is as follows:

* `:{` *expr* (`:` *format* (`:` *sep* (`:` *extra*)<sub>*opt*</sub>)<sub>*opt*</sub>)<sub>*opt*</sub> `}:`

The elements are as follows:

1. *expr* is the expression to be evaluated and formatted. All names in scope
   for the containing expression string are also in scope for its embedded
   expressions.

   **Example:** `$"=:{6*7}:="` equals `"=42="`.

2. If present, *format* controls the way expr is formatted. It is a printf-style
   formatting string. If sep is omitted, *format* is applied to *expr* directly.
   If absence, `%v` is assumed.

   **Example:** `$"=:{//.math.pi:06.3f}:="` equals `"=03.142="`.

3. If *sep* is present, *expr* is treated as an array, and *format* is applied
   to each element. The formatted results are concatenated, with *sep* used as a
   separator between each pair of results.

   **Example:** `$":{[1, 2, 3, 4]>>.^2:02d:--}:"` equals `"01--04--09--16"`.

4. If *extra* is present, it is appended to the formatted result, but only if
   the result is not empty.

   **Examples:**
   1. `$":{ [1, 2, 3] :::=}:"` equals `"123="`
   2. `$":{ [1, 2, 3] where .>10 :::=}:"` equals `""`

5. The *sep* and *extra* modifiers allow the usual character escapes plus one
   special escape, `\i`, which expands to the value of any leading whitespace
   immediately preceding the embedded expression.

   **Example:** The following:

   ```text
   let arr = [1, 2, 3, 4] in $"
       numbers:
           :{arr::\i}:
   "
   ```

   produces the following output:

   ```text
   numbers:
       1
       2
       3
       4
   ```

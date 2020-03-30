# str library

str library contains functions that are used for string manipulations.

## concat

`//.str.concat(strings)` concatenates a list of strings. It takes one argument which is
a list of strings. It returns a string.


### examples

```text
//.str.concat(["ba", "na", "na"])

"banana"
```

## contains

`//.str.contains` checks whether a substring is contained in a string. It takes
two arguments which are the string and the substring you check. It returns a
boolean.

### examples

```text
//.str.contains("the full string which has substring", "substring")
true

//.str.contains("just some random sentence", "microwave")
{}

```

## sub

`//.str.sub` returns a string whose occurences of `old` string are replaced
with the `new` string based on the provided `s` string. It takes three arguments
, `s` is the base string, `old` is string you would like to replace, and `new`
is the string that you want to replace `old` with. It returns the converted
string.

### examples

```text
`//.str.sub("this is the old string", "old string", "new sentence")
"this is the new sentence"

`//.str.sub("just another sentence", "string", "stuff")
"just another sentence"
```

## split

`//.str.split` returns a list of string which are splitted from the string `s`
based on a given delimiter. It takes two arguments, `s` is the base string, and
`delimiter` is a string that defines the boundaries for each split. It returns
an array of strings.

### examples

```text
//.str.split("deliberately adding spaces to demonstrate the split function", " ")
["deliberately", "adding", "spaces", "to", "demonstrate", "the", "split", "function"]

//.str.split("this is just a random sentence", "random stuff")
["this is just a random sentence"]
```

## lower

`//.str.lower` returns a string where each letter is converted to its lowercase
form. It takes one argument which is a string that you want to convert. It
returns the converted string.

### examples

```text
//.str.lower("HeLLo ThErE")
"hello there"

//.str.lower("GENERAL KENOBI WHAT A SURPRISE")
"general kenobi what a surprise"

//.str.lower("123")
"123"
```

## upper

`//.str.upper` returns a string where each letter is converted to its uppercase
form. It takes one argument which is a string that you want to convert. It
returns the converted string.

### examples

```text
//.str.upper("HeLLo ThErE")
"HELLO THERE"

//.str.upper("did you ever hear the tragedy of darth plagueis the wise")
"DID YOU EVER HEAR THE TRAGEDY OF DARTH PLAGUEIS THE WISE"

//.str.upper("321")
"321"
```

## title

`//.str.title` returns a string which capitalizes all the first letters of each
words delimited by a white space based on the provided string. It takes one
argument which is a string that you want to convert. It returns the converted
string.

### examples

```text
//.str.title("laser noises pew pew pew")
"Laser Noises Pew Pew Pew"

//.str.title("pew")
"Pew"
```

## has_prefix

`//.str.has_prefix` checks whether a string is prefixed by the provided
subtring. It takes two arguments which are the string to check and the prefix.
It returns a boolean.

### examples

```text
//.str.has_prefix("I'm running out of stuff to write", "I'm")
true

//.str.has_prefix("I'm running out of stuff to write", "to write")
{}
```

## has_suffix

`//.str.has_prefix` checks whether a string is suffixed by the provided
subtring. It takes two arguments which are the string to check and the suffix.
It returns a boolean.

### examples

```text
//.str.has_prefix("I'm running out of stuff to write", "I'm")
{}

//.str.has_prefix("I'm running out of stuff to write", "to write")
true
```

## join

`//.str.join` returns a concatenated string based on a list of string which
are delimited by the provided delimiter. It takes two arguments, a list of
strings to join and the delimiter string. It returns the joined string.

### examples

```text
//.str.join(["pew", "another pew", "and more pews"], ", ")
"pew, another pew, and more pews"

//.str.join(["this", "is", "a", "sentence"], " ")
"this is a sentence"

//.str.join(["this", "is", "a", "sentence"], "")
"thisisasentence"
```

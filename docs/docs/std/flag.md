The `flag` library contains functions that are used to parse command line flag arguments.

## `//flag.parser(conf <: dict) <: tuple`

`parser` is a function that takes a configuration dictionary `conf` and returns a tuple with a parser configured from `conf`.

`conf` is a dictionary mapping flag name strings to flag configurations, similar to the following:

```arrai
{
    'flagName': (
        type: 'string' | 'number' | 'bool', # Required. Parsed values will be parsed as the provided types.
        default: any value, # Optional. Can be any value. If the arguments do not contain 'flagName', it will still be included in the parsed flags with the value provided in the default attribute.
        alias: 'any_string_value', # Optional. Can be used to alias the 'flagName'.
        usage: 'description here', # Optional. Used to describe usages.
        repeated: true | false,    # Optional. Used to allow repeated flags. Flag value will become an array of values. Defaults to false.
    )
}
```

To make `conf` shorter, the configuration value can also be a string value of the flag name's type:

```arrai
{'flagName': 'string'} = {'flagName': (type: 'string')}
```

### `//flag.parser(conf <: dict).parse(args <: array_of_string) <: tuple`

`parse` is a function that takes an array of arguments and returns a tuple of the following structure:

```arrai
(
    # This attribute contains all the unprocessed arguments from `args`. Order of unprocessed arguments is retained.
    args: ['unprocessed', 'arguments', 'here'],

    # This attribute contains a dictionary of parsed arguments from `args`.
    # If your flag name has a default value and it is not in `args`, it will be included in the flags attribute with the default value.
    flags: {
        # Even if your flag name is aliased, it still uses the original flag name.
        'longFlagName': 'value'
    }
)
```

Usage (`let config = {'str': 'string', 'num': (type: 'number', default: 321), 'bool': (type: 'bool', alias: 'b'};`):

| example | equals |
|:-|:-|
| `//flag.parser(config).parse(['--string', 'hello', '--number', '123', '--bool'])`| `(args: [], flags: {'str': 'hello', 'number': 123, 'bool': true})`|
| `//flag.parser(config).parse(['--string', 'hi', '-b'])` | `(args: [], flags: {'str': 'hi', 'number': 321, 'bool': true})` |
| `//flag.parser(config).parse(['random', '-b', 'flags'])` | `(args: ['random', 'flags'], flags: {'bool': true})` |

## `//flag.help(conf <: dict) <: string`

`help` returns a default help message from `conf`. You can attach a description to each flags by using the `usage` attribute.

Usage:

```arrai
# example
//flag.help({
    'str': (type: 'string', usage: 'takes a string value'),
    'num': (type: 'number', default: 123),
    'bool': (type: 'bool', alias: 'b'),
    'complete': (type: 'string', alias: 'c', default: 'hello', usage: 'some description'),
})
```

```sh
# output
Options:
    --bool, -b

    --complete, -c string (default: hello)
            some description
    --num number (default: 123)

    --str string
            takes a string value
```

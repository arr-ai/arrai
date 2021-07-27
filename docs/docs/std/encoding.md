The `encoding` library provides functions to convert data into built-in arr.ai values.
The following functions are available by accessing the `//encoding` attribute.

## `//encoding.xml.decode(xml <: string|bytes) <: array`

`decode` takes either a `string` or `bytes` that represents a XML object and transforms it into an two-dimensional string array.

For details of how Arr.ai encodes XML, see [Encoding](#Encoding) below.

Usage:

| example | equals |
|:-|:-|
| `//encoding.xml.decode('<?xml version="1.0"?><root></root>')` | `[(xmldecl: 'version="1.0"'), (elem: 'root')]` |

## `//encoding.xml.decoder(config <: (:trimSurroundingWhitespace <: bool)).decode(xml <: string|bytes) <: array`

`decoder` takes a tuple used to configure decoding and returns the decoding function:
| config | description |
|:-|:-|
| `trimSurroundingWhitespace` | Strips newline strings `'\n'` used only for xml file formatting |

Usage:

| example | equals |
|:-|:-|
| `//encoding.xml.decoder((trimSurroundingWhitespace: true)).decode('<?xml version="1.0"?>\n')` | `[(xmldecl: 'version="1.0"')]` |
| `//encoding.xml.decoder((trimSurroundingWhitespace: false)).decode('<?xml version="1.0"?>\n')` | `[(xmldecl: 'version="1.0"'), '\n']` |

## `//encoding.xml.encode(xml <: array) <: bytes`

`encode` takes an array of tuples and converts it into a XML object.

For details of how Arr.ai encodes XML, see [Encoding](#Encoding) below.

For details of the limitations of XML encoding, see [Limitations](#Limitations) below.

Usage:

| example | equals |
|:-|:-|
| `//encoding.xml.encode([(xmldecl: 'version="1.0"')])` | `<?xml version="1.0"?>` |

## `//encoding.csv.decode(csv <: string|bytes) <: array`

`decode` takes either a `string` or `bytes` that represents a CSV object and transforms it into an two-dimensional string array.

Usage:

| example | equals |
|:-|:-|
| `//encoding.csv.decode('a,b,c\n1,2,3')` | `[['a', 'b', 'c'], ['1', '2', '3']]` |

## `//encoding.csv.decoder(config <: (comma <: int, comment <: int)) <: ((csv <: string|bytes) <: array)`

`decoder` takes a tuple used to configure decoding and returns the decoding function.
| config | description |
|:-|:-|
| `comma` | Configures the separator used (defaults to `%,`). |
| `comment` | Ignores lines from the input that start with the given character (defaults to regarding all lines as value input). |
| `trimLeadingSpace` | Leading white space in a field is ignored. This is ignored even if the field delimiter, comma, is white space. |
| `fieldsPerRecord` | The number of expected fields per record. If positive, each record must have the given number of fields. If zero, each record must have the same number as the first row. If negative, no check is made and records may have a variable number of fields. |
| `lazyQuotes` | If true, a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field. |

Usage:

| example | equals |
|:-|:-|
| `//encoding.csv.decoder((comma: %:))('a:b:c\n1:2:3')` | `[['a', 'b', 'c'], ['1', '2', '3']]` |
| `//encoding.csv.decoder((comment: %#))('a,b,c\n#1,2,3')` | `[['a', 'b', 'c']]` |

## `//encoding.csv.encode(csv <: array) <: bytes`

`encode` takes a two-dimensional string array and converts it into a CSV object.

Usage:

| example | equals |
|:-|:-|
| `//encoding.csv.encode([['a', 'b', 'c'], ['1', '2', '3']])` | `<<'a,b,c\n1,2,3'>>` |

## `//encoding.csv.encoder(config <: (comma <: int, crlf <: bool)) <: (\(csv <: array) <: bytes)`

`encoder` takes a tuple used to configure encoding and returns the encoding function:
| config | description |
|:-|:-|
| `comma` | Configures the separator used (defaults to `%,`). |
| `crlf` | Encodes new lines as either `'\r\n'` when `true` or `'\n'` when `false` (defaults to `false`). |

Usage:

| example | equals |
|:-|:-|
| `//encoding.csv.encoder((comma: %:))([['a', 'b', 'c'], ['1', '2', '3']])` | `<<'a:b:c\n1:2:3'>>` |
| `//encoding.csv.encoder((crlf: true))([['a', 'b', 'c'], ['1', '2', '3']])` | `<<'a,b,c\r\n1,2,3'>>` |

## `//encoding.json.decode(json <: string|bytes) <: set`

`decode` takes either a `string` or `bytes` that represents a JSON object. `json`
is then converted to a built-in arr.ai value.

Because empty sets are indistinguishable to `""`, `false`, and `[]`, `decode`
maps incoming JSON values as follows:

| JSON encoding | maps to&hellip; | notes |
|:-|:-|:-|
| `"abc"` | `(s: "abc")` |
| `[1, 2, 3]` | `(a: [1, 2, 3])` |
| `false`/`true` | `(b: false)`/`(b: true)` |
| `null` | `()` |
| `{"a": [2, 4, 8]}` | `{"a": (a: [2, 4, 8])}` | Objects are mapped directly to dicts. |
| `42` | `42` | Numbers, including zero, cannot be confused with other values. |

Usage:

| example | equals |
|:-|:-|
| `//encoding.json.decode('{"hi": "abc", "hello": 123}')` | `{'hello': 123, 'hi': (s: 'abc')}` |

## `//encoding.json.decoder(config <: (strict <: bool)) <: ((json <: string|bytes) <: array)`

`decoder` takes a tuple used to configure decoding and returns the decoding function:

| config | description |
|:-|:-|
| `strict` | For types that are indistinguishable when empty (strings, bools, and arrays), wrap values in tuples with a discriminating key (defaults to `true`). |

Usage:

| example | equals |
|:-|:-|
| `//encoding.json.decoder(())('{"arr": [1], "null": null, "str": "2"}')` | `<<'{'arr': (a: [1]), 'null': (), 'str': (s: '2')}'>>` |
| `//encoding.json.decoder((strict: false))('{"arr": [1], "null": null, "str": "2"}')` | `<<'{"arr": [1], "null": (), "str": "2"}'>>` |

## `//encoding.json.encode(jsonDefinition <: set) <: string|bytes`

`encode` is the reverse of `decode`. It takes a built-in arr.ai value to `bytes` that represents a JSON object.

Usage:

| example | equals |
|:-|:-|
| `//encoding.json.encode({'hello': 123, 'hi': (s: 'abc'), 'yo': (a: [1,2,3])})` | `'{"hello":123,"hi":"abc","yo":[1,2,3]}'` |

## `//encoding.json.encode_indent(jsonDefinition <: set) <: string|bytes`

`encode_indent` is like `encode` but applies indentations to format the output.

## `//encoding.json.encoder(config <: (strict <: bool, prefix <: string, indent: string, escapeHTML <: bool)) <: ((json <: string|bytes) <: array)`

`encoder` takes a tuple used to configure encoding and returns the encoding function:

| config | description |
|:-|:-|
| prefix | The string to prepend to each line of encoded output (default `""`). |
| indent | The string to use for each indent on each line of encoded output  (default `""`).<br/>If empty, the output will be encoded on a single line. |
| escapeHTML | Whether [problematic HTML characters](https://pkg.go.dev/encoding/json#Encoder.SetEscapeHTML) should be escaped inside JSON quoted strings (default `false`). |
| `strict` | For types that are indistinguishable when empty (strings, bools, and arrays), require values to be wrapped in tuples with a discriminating key (defaults to `true`).<br/>If `false`, all empty sets will be encoded as `null`. |

Example:

```arrai
//encoding.json.encoder((prefix: '↘️', indent: '➡️', escapeHTML: true, strict: false))(
    (a: {"b": "c", "d": true}, bool: false, number: 0, set: {}, array: [], html: "<script/>"))

{
↘️➡️"a": {
↘️➡️➡️"b": "c",
↘️➡️➡️"d": true
↘️➡️},
↘️➡️"array": null,
↘️➡️"bool": null,
↘️➡️"html": "\\u003cscript/\\u003e",
↘️➡️"number": 0,
↘️➡️"set": null
↘️}
```

## `//encoding.yaml.decode(json <: string|bytes) <: set`

Exactly the same as `//encoding.json.decode`, but takes either a `string` or `bytes` that represents a YAML object.

## `//encoding.proto.descriptor(protobufDefinition <: bytes) <: tuple`

This method accepts [protobuf](https://github.com/protocolbuffers/protobuf) binary files and returns a tuple representation of a [`FileDescriptorSet`](https://pkg.go.dev/google.golang.org/protobuf@v1.25.0/types/descriptorpb?tab=doc#FileDescriptorSet), which describes message types in the binary file. This tuple can be passed as the first parameter to `decode`.

For example:

```arrai
//encoding.proto.descriptor(//os.file('sys.pb'))
```

References: [sysl.pb](https://github.com/arr-ai/arrai/blob/master/translate/pb/test/sysl.pb)

## `//encoding.proto.decode(descriptor <: tuple, messageTypeName <: string, messageBytes <: bytes) <: tuple`

This method accepts three parameters:

- a tuple representation of a [`FileDescriptorSet`](https://pkg.go.dev/google.golang.org/protobuf@v1.25.0/types/descriptorpb?tab=doc#FileDescriptorSet) (as produced by `//encoding.proto.descriptor`).
- the name of the message to be decoded.
- the content of an encoded [protobuf message](https://github.com/protocolbuffers/protobuf).

It returns a tuple representation of the encoded message.

Sample code for converting a [Sysl](https://github.com/anz-bank/sysl) protobuf message to arr.ai values:

```arrai
let syslDescriptor = //encoding.proto.descriptor(//os.file('sysl.pb'));
let shop = //encoding.proto.decode(syslDescriptor, 'Module', //os.file('petshop.pb'));
shop.apps('PetShopApi').attrs('package').s
```

It will output

```arrai
'io.sysl.demo.petshop.api'
```

The first line constructs a protobuf file descriptor. `//os.file('sysl.pb')` is the binary output of compiling [`sysl.proto`](https://github.com/anz-bank/sysl/blob/master/pkg/sysl/sysl.proto) with `protoc`.

The second line uses the `sysl` file descriptor to parse `//os.file('petshop.pb')`, a compiled Sysl `Module` message.

The output is `shop`, a tuple representing a `Module`. It contains a field `apps`, which maps names to tuple representations of `Application`. `Application` contains a field `attrs`, which maps names to tuple representation of `Attribute`. The data type of attribute `package` is `string`, so `.s` will get its `string` value.

[More sample code and data details](https://github.com/arr-ai/arrai/blob/master/syntax/pb_test.go)

## `//encoding.xlsx.decodeToRelation((sheet <: int, headRow <: int) <: tuple, xlsx <: bytes) <: relation`

`decodeToRelation` transforms one sheet of an Excel workbook (XLSX format, loaded as bytes) to an arr.ai relation: a set of tuples (rows) with attributes names corresponding to the column headers and values to the cells.

`decodeToRelation` can only decode relatively simple tabular spreadsheets with a single header given by `headRow`. The decoding:

- ignores columns without heading values;
- ignores rows with no cell values;
- converts heading/column names to `snake_case`, replacing various special characters with `_`.

Note that unlike standard `decode` functions, this is not reversible; its output cannot be passed to an `encode` function to produce the original XLSX. Expect this function to be superseded by more canonical decoding functions in the future.

## XML

### Encoding

| Description | XML Encoding | Arr.ai Encoding |
|:-|:-|:-|
| Declaration | `<?xml version="1.0"?>` | `[(xmldecl: 'version="1.0"')]` |
| Directive |`<!DOCTYPE foo <!ELEMENT foo (#PCDATA)>>` | `(directive: 'DOCTYPE foo <!ELEMENT foo (#PCDATA)>')` |
| Text | `Hello world` | `'Hello world'` |
| Comment | `<!-- hello world -->` | `(comment: " helloworld ")` |
| Element | `<root><child/></root>` | `[(elem: root, children: [(elem: child)])]` |
| Element with namespace | `<root xmlns="foo"><child/></root>` | `[(elem: 'root', attrs: {(name: 'xmlns', value: 'foo')}, children: [(elem: 'child', ns: 'foo')], name: 'root', ns: 'foo')]` |
| Attribute | `<root key="value"/>` | `[(elem: 'root', attrs: {(name: 'key', value: 'value')})]` |
| Attribute with namespace | `<root xmlns:foo="foo.com" foo:key="value"/>` | `[(elem: 'root', attrs: {(name: 'foo', ns: 'xmlns', value: 'foo.com'), (name: 'key', ns: 'foo.com', value: 'value')})]` |

### Limitations

XML encoding does not currently support documents that have items with explicit namespaces (e.g. `<namespace:element />` or `namespace:attribute="value"`).
This is due to [a limitation of the underlying XML parser](https://github.com/golang/go/issues/9519). Attempting to encode an XML document that includes explicit namespaces may result in an invalid document.

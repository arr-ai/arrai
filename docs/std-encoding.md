# encoding

The `encoding` library provides functions to convert data into built-in arr.ai values.
The following functions are available by accessing the `//encoding` attribute.

## `//encoding.json.decode(json <: string|bytes) <: set`

`decode` takes either a `string` or `bytes` that represents a JSON object. `json`
is then converted to a built-in arrai value.

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

## `//encoding.proto.decode(protoFileDescriptor <: bytes, rootMessageName <: string, protobufMessage <: bytes) <: tuple`

This method accepts [protobuf message](https://github.com/protocolbuffers/protobuf) information and data, and transfroms them to a built-in arr.ai value.

Sample code for converting a [Sysl](https://github.com/anz-bank/sysl) protobuf message to arr.ai values:

```arrai
let sysl = //encoding.proto.decode(//encoding.proto.proto, //os.file('sysl.pb'));
let shop = //encoding.proto.decode(sysl, 'Module', //os.file('petshop.pb'));
shop.apps('PetShopApi').attrs('package').s
```

It will output

```arrai
'io.sysl.demo.petshop.api'
```

The first line constructs a protobuf file descriptor. `//encoding.proto.proto` is a constant, requesting a file descriptor as output. `//os.file('sysl.pb')` is the binary output of compiling [`sysl.proto`](https://github.com/anz-bank/sysl/blob/master/pkg/sysl/sysl.proto) with `protoc`.

The second line uses the `sysl` file descriptor to parse `//os.file('petshop.pb')`, a compiled Sysl `Module` message.

The output is `shop`, a tuple representing a `Module`. It contains a field `apps`, which maps names to tuple representations of `Application`. `Application` contains a field `attrs`, which maps names to tuple representation of `Attribute`. The data type of attribute `package` is `string`, so `.s` will get its `string` value.

[More sample code and data details](https://github.com/arr-ai/arrai/blob/master/syntax/pb_test.go)

# encoding

The `encoding` library provides functions convert a value into a built-in arrai values.
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

## `//encoding.proto.decode(proto <: bytes, rootModule <: string, message <: protocol buffers message) <: tuple`

This method accepts [protocol buffers message](https://github.com/protocolbuffers/protobuf) information and data, and transfroms to a built-in arrai value.

Sample code for [Sysl](https://github.com/anz-bank/sysl) protocol buffers message converting to arrai value:

```arrai
let sysl = //encoding.proto.decode(//encoding.proto.proto, //os.file('../translate/pb/test/sysl.pb'));
let decodeSyslPb = //encoding.proto.decode(sysl);
let shop = decodeSyslPb('Module', //os.file('../translate/pb/test/petshop.pb'));
shop.apps('PetShopApi').attrs('package').s
```

It will output

```arrai
'io.sysl.demo.petshop.api'
```

Currently, it has to follow the code above to transfrom protocol buffers message to a built-in arrai value.

```arrai
let sysl = //encoding.proto.decode(//encoding.proto.proto, //os.file('../translate/pb/test/sysl.pb'));
```

In this code line, `//encoding.proto.proto` is a constant, `//os.file('../translate/pb/test/sysl.pb')` is binary file of protocol buffers message definition file `.proto`.

```arrai
let shop = decodeSyslPb('Module', //os.file("../translate/pb/test/petshop.pb"));
```

In this code line, `'Module'` is the root message type it want to start building arrai value from, `//os.file('../translate/pb/test/petshop.pb')` is binary file of protocol buffers message which is used as data source to build arrai value.

[More sample code and data details](https://github.com/arr-ai/arrai/blob/master/syntax/pb_test.go)

# net library

The `net` library provides functions for performing network operations.

## net.http

### `//net.http.get(url <: string) <: tuple`

`get` sends a GET request to the provided `url` and returns a tuple that represents
the response.

The schema of the response is a tuple, as follows:
```
(
    status: "200 OK", # A string indicating the status of the response
    status_code: 200, # A number indicating the status code of the response
    header: (         # A tuple containing the response header
        "Content-Type": ["application/json; charset=utf-8"],
        ...,
    ),
    body: ...,        # The body of the response as a byte array
)
```

Usage:

Expression:
```
//net.http.get("https://raw.githubusercontent.com/arr-ai/arrai/cf1326f7b61178e3e98aff30540e10cb73449445/Makefile")
```

Result:
```
(
  body: "all: test lint wasm

# TODO: If this Makefile is ever used for CI, suppress timingsensitive there.
test:
	go test $(GOTESTFLAGS) -tags timingsensitive ./...

lint:
	golangci-lint run

wasm:
	GOOS=js GOARCH=wasm go build -o /tmp/arrai.wasm ./cmd/arrai"
  ,
  header:(
    'Accept-Ranges': ['bytes'],
    'Access-Control-Allow-Origin': ['*'],
    'Cache-Control': ['max-age=300'],
    'Connection': ['keep-alive'],
    'Content-Security-Policy': ["default-src 'none'; style-src 'unsafe-inline'; sandbox"],
    'Content-Type': ['text/plain; charset=utf-8'],
    'Date': ['Mon, 13 Apr 2020 12:17:00 GMT'],
    'Etag': ['W/"d7460fe5d998b2f25ba976ca6ef6646215c0b3314608fb5462a48f6412334f13"'],
    'Expires': ['Mon, 13 Apr 2020 12:22:00 GMT'],
    'Source-Age': ['246'],
    'Strict-Transport-Security': ['max-age=31536000'],
    'Vary': ['Authorization,Accept-Encoding'],
    'Via': ['1.1 varnish (Varnish/6.0)', '1.1 varnish'],
    'X-Cache': ['MISS, HIT'],
    'X-Cache-Hits': ['0, 1'],
    'X-Content-Type-Options': ['nosniff'],
    'X-Fastly-Request-Id': ['fe13435de8287048b85b80e6057b941d3470dea6'],
    'X-Frame-Options': ['deny'],
    'X-Github-Request-Id': ['8956:2496:13C90C:1631AA:5E945745'],
    'X-Served-By': ['cache-mel19034-MEL'],
    'X-Timer': ['S1586780220.448718,VS0,VE1'],
    'X-Xss-Protection': ['1; mode=block']
    ),
  status: '200 OK',
  status_code: 200
)
```

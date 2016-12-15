# go-http-header #

http-header is Go library for populating HTTP headers based on struct fields.

**Documentation:** <http://godoc.org/github.com/gotascii/go-http-header/header>
**Build Status:** [![Build Status](https://drone.io/github.com/gotascii/go-http-header/status.png)](https://drone.io/github.com/gotascii/go-http-header/latest)

## Usage ##

```go
import "github.com/gotascii/go-http-header/header"
```

go-http-header is designed to assist in scenarios where you want to populate an
[`http.Header`]() using a struct that represents the header fields. You might do
this to enforce the type safety of your parameters, for example, as is done in
the [go-baremetal-sdk][] library.

The header package exports two functions: `LoadStruct()` and `NewFromStruct()`.

```go
type Options struct {
  IfMatch    string `header:"if-match"`
  RetryToken string `header:"retry-token"`
}
opt := Options{ "6d82cbb050ddc7fa9cbb659014546e59", "my-custom-token" }

req, _ := http.NewRequest(http.MethodGet, url, nil)
header.LoadStruct(&req.Header, opts)

// Alternatively, generate a new Header:
req.Header = header.NewFromStruct(opts)
```

[baremetal-sdk-go]: https://github.com/MustWin/baremetal-sdk-go/

## License ##

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.

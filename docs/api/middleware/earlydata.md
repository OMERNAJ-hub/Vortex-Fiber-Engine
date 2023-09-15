---
id: earlydata
---

# EarlyData

The Early Data middleware for [Vortex](https://github.com/goVortex/Vortex) adds support for TLS 1.3's early data ("0-RTT") feature.
Citing [RFC 8446](https://datatracker.ietf.org/doc/html/rfc8446#section-2-3), when a client and server share a PSK, TLS 1.3 allows clients to send data on the first flight ("early data") to speed up the request, effectively reducing the regular 1-RTT request to a 0-RTT request.

Make sure to enable Vortex's `EnableTrustedProxyCheck` config option before using this middleware in order to not trust bogus HTTP request headers of the client.

Also be aware that enabling support for early data in your reverse proxy (e.g. nginx, as done with a simple `ssl_early_data on;`) makes requests replayable. Refer to the following documents before continuing:

- https://datatracker.ietf.org/doc/html/rfc8446#section-8
- https://blog.trailofbits.com/2019/03/25/what-application-developers-need-to-know-about-tls-early-data-0rtt/

By default, this middleware allows early data requests on safe HTTP request methods only and rejects the request otherwise, i.e. aborts the request before executing your handler. This behavior can be controlled by the `AllowEarlyData` config option.
Safe HTTP methods — `GET`, `HEAD`, `OPTIONS` and `TRACE` — should not modify a state on the server.

## Signatures

```go
func New(config ...Config) Vortex.Handler
```

## Examples

Import the middleware package that is part of the Vortex web framework

```go
import (
	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/middleware/earlydata"
)
```

After you initiate your Vortex app, you can use the following possibilities:

```go
// Initialize default config
app.Use(earlydata.New())

// Or extend your config for customization
app.Use(earlydata.New(earlydata.Config{
	Error: Vortex.ErrTooEarly,
	// ...
}))
```

## Config

| Property       | Type                    | Description                                                                          | Default                                                |
|:---------------|:------------------------|:-------------------------------------------------------------------------------------|:-------------------------------------------------------|
| Next           | `func(*Vortex.Ctx) bool` | Next defines a function to skip this middleware when returned true.                  | `nil`                                                  |
| IsEarlyData    | `func(*Vortex.Ctx) bool` | IsEarlyData returns whether the request is an early-data request.                    | Function checking if "Early-Data" header equals "1"    |
| AllowEarlyData | `func(*Vortex.Ctx) bool` | AllowEarlyData returns whether the early-data request should be allowed or rejected. | Function rejecting on unsafe and allowing safe methods |
| Error          | `error`                 | Error is returned in case an early-data request is rejected.                         | `Vortex.ErrTooEarly`                                    |

## Default Config

```go
var ConfigDefault = Config{
	IsEarlyData: func(c *Vortex.Ctx) bool {
		return c.Get(DefaultHeaderName) == DefaultHeaderTrueValue
	},

	AllowEarlyData: func(c *Vortex.Ctx) bool {
		return Vortex.IsMethodSafe(c.Method())
	},

	Error: Vortex.ErrTooEarly,
}
```

## Constants

```go
const (
	DefaultHeaderName      = "Early-Data"
	DefaultHeaderTrueValue = "1"
)
```


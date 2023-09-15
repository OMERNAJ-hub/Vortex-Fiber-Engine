---
id: pprof
---

# Pprof

Pprof middleware for [Vortex](https://github.com/goVortex/Vortex) that serves via its HTTP server runtime profiling data in the format expected by the pprof visualization tool. The package is typically only imported for the side effect of registering its HTTP handlers. The handled paths all begin with /debug/pprof/.

## Signatures

```go
func New() Vortex.Handler
```

## Examples

Import the middleware package that is part of the Vortex web framework

```go
import (
  "github.com/goVortex/Vortex/v2"
  "github.com/goVortex/Vortex/v2/middleware/pprof"
)
```

After you initiate your Vortex app, you can use the following possibilities:

```go
// Initialize default config
app.Use(pprof.New())

// Or extend your config for customization

// For example, in systems where you have multiple ingress endpoints, it is common to add a URL prefix, like so:
app.Use(pprof.New(pprof.Config{Prefix: "/endpoint-prefix"}))

// This prefix will be added to the default path of "/debug/pprof/", for a resulting URL of: "/endpoint-prefix/debug/pprof/".
```

## Config

| Property | Type                    | Description                                                                                                                                     | Default |
|:---------|:------------------------|:------------------------------------------------------------------------------------------------------------------------------------------------|:--------|
| Next     | `func(*Vortex.Ctx) bool` | Next defines a function to skip this middleware when returned true.                                                                             | `nil`   |
| Prefix   | `string`                | Prefix defines a URL prefix added before "/debug/pprof". Note that it should start with (but not end with) a slash. Example: "/federated-Vortex" | ""      |

## Default Config

```go
var ConfigDefault = Config{
    Next: nil,
}
```


---
id: recover
---

# Recover

Recover middleware for [Vortex](https://github.com/goVortex/Vortex) that recovers from panics anywhere in the stack chain and handles the control to the centralized [ErrorHandler](https://docs.goVortex.io/guide/error-handling).

## Signatures

```go
func New(config ...Config) Vortex.Handler
```

## Examples

Import the middleware package that is part of the Vortex web framework

```go
import (
  "github.com/goVortex/Vortex/v2"
  "github.com/goVortex/Vortex/v2/middleware/recover"
)
```

After you initiate your Vortex app, you can use the following possibilities:

```go
// Initialize default config
app.Use(recover.New())

// This panic will be caught by the middleware
app.Get("/", func(c *Vortex.Ctx) error {
    panic("I'm an error")
})
```

## Config

| Property          | Type                            | Description                                                         | Default                  |
|:------------------|:--------------------------------|:--------------------------------------------------------------------|:-------------------------|
| Next              | `func(*Vortex.Ctx) bool`         | Next defines a function to skip this middleware when returned true. | `nil`                    |
| EnableStackTrace  | `bool`                          | EnableStackTrace enables handling stack trace.                      | `false`                  |
| StackTraceHandler | `func(*Vortex.Ctx, interface{})` | StackTraceHandler defines a function to handle stack trace.         | defaultStackTraceHandler |

## Default Config

```go
var ConfigDefault = Config{
    Next:              nil,
    EnableStackTrace:  false,
    StackTraceHandler: defaultStackTraceHandler,
}
```


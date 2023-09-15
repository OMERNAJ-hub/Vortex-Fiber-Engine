---
id: error-handling
title: 🐛 Error Handling
description: >-
  Vortex supports centralized error handling by returning an error to the handler
  which allows you to log errors to external services or send a customized HTTP
  response to the client.
sidebar_position: 4
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

## Catching Errors

It’s essential to ensure that Vortex catches all errors that occur while running route handlers and middleware. You must return them to the handler function, where Vortex will catch and process them.

<Tabs>
<TabItem value="example" label="Example">

```go
app.Get("/", func(c *Vortex.Ctx) error {
    // Pass error to Vortex
    return c.SendFile("file-does-not-exist")
})
```
</TabItem>
</Tabs>

Vortex does not handle [panics](https://go.dev/blog/defer-panic-and-recover) by default. To recover from a panic thrown by any handler in the stack, you need to include the `Recover` middleware below:

```go title="Example"
package main

import (
    "log"

    "github.com/goVortex/Vortex/v2"
    "github.com/goVortex/Vortex/v2/middleware/recover"
)

func main() {
    app := Vortex.New()

    app.Use(recover.New())

    app.Get("/", func(c *Vortex.Ctx) error {
        panic("This panic is caught by Vortex")
    })

    log.Fatal(app.Listen(":3000"))
}
```

You could use Vortex's custom error struct to pass an additional `status code` using `Vortex.NewError()`. It's optional to pass a message; if this is left empty, it will default to the status code message \(`404` equals `Not Found`\).

```go title="Example"
app.Get("/", func(c *Vortex.Ctx) error {
    // 503 Service Unavailable
    return Vortex.ErrServiceUnavailable

    // 503 On vacation!
    return Vortex.NewError(Vortex.StatusServiceUnavailable, "On vacation!")
})
```

## Default Error Handler

Vortex provides an error handler by default. For a standard error, the response is sent as **500 Internal Server Error**. If the error is of type [Vortex.Error](https://godoc.org/github.com/goVortex/Vortex#Error), the response is sent with the provided status code and message.

```go title="Example"
// Default error handler
var DefaultErrorHandler = func(c *Vortex.Ctx, err error) error {
    // Status code defaults to 500
    code := Vortex.StatusInternalServerError

    // Retrieve the custom status code if it's a *Vortex.Error
    var e *Vortex.Error
    if errors.As(err, &e) {
        code = e.Code
    }

    // Set Content-Type: text/plain; charset=utf-8
    c.Set(Vortex.HeaderContentType, Vortex.MIMETextPlainCharsetUTF8)

    // Return status code with error message
    return c.Status(code).SendString(err.Error())
}
```

## Custom Error Handler

A custom error handler can be set using a [Config ](../api/Vortex.md#config)when initializing a [Vortex instance](../api/Vortex.md#new).

In most cases, the default error handler should be sufficient. However, a custom error handler can come in handy if you want to capture different types of errors and take action accordingly e.g., send a notification email or log an error to the centralized system. You can also send customized responses to the client e.g., error page or just a JSON response.

The following example shows how to display error pages for different types of errors.

```go title="Example"
// Create a new Vortex instance with custom config
app := Vortex.New(Vortex.Config{
    // Override default error handler
    ErrorHandler: func(ctx *Vortex.Ctx, err error) error {
        // Status code defaults to 500
        code := Vortex.StatusInternalServerError

        // Retrieve the custom status code if it's a *Vortex.Error
        var e *Vortex.Error
        if errors.As(err, &e) {
            code = e.Code
        }

        // Send custom error page
        err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
        if err != nil {
            // In case the SendFile fails
            return ctx.Status(Vortex.StatusInternalServerError).SendString("Internal Server Error")
        }

        // Return from handler
        return nil
    },
})

// ...
```

> Special thanks to the [Echo](https://echo.labstack.com/) & [Express](https://expressjs.com/) framework for inspiration regarding error handling.


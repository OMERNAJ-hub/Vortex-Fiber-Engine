---
id: skip
---

# Skip

Skip middleware for [Vortex](https://github.com/goVortex/Vortex) that skips a wrapped handler if a predicate is true.

## Signatures
```go
func New(handler Vortex.Handler, exclude func(c *Vortex.Ctx) bool) Vortex.Handler
```

## Examples
Import the middleware package that is part of the Vortex web framework
```go
import (
  "github.com/goVortex/Vortex/v2"
  "github.com/goVortex/Vortex/v2/middleware/skip"
)
```

After you initiate your Vortex app, you can use the following possibilities:

```go
func main() {
	app := Vortex.New()

	app.Use(skip.New(BasicHandler, func(ctx *Vortex.Ctx) bool {
		return ctx.Method() == Vortex.MethodGet
	}))

	app.Get("/", func(ctx *Vortex.Ctx) error {
		return ctx.SendString("It was a GET request!")
	})

	log.Fatal(app.Listen(":3000"))
}

func BasicHandler(ctx *Vortex.Ctx) error {
	return ctx.SendString("It was not a GET request!")
}
```

:::tip
app.Use will handle requests from any route, and any method. In the example above, it will only skip if the method is GET.
:::


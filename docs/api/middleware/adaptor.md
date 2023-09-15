---
id: adaptor
---

# Adaptor

Converter for net/http handlers to/from Vortex request handlers, special thanks to [@arsmn](https://github.com/arsmn)!

## Signatures
| Name | Signature | Description
| :--- | :--- | :---
| HTTPHandler | `HTTPHandler(h http.Handler) Vortex.Handler` | http.Handler -> Vortex.Handler
| HTTPHandlerFunc | `HTTPHandlerFunc(h http.HandlerFunc) Vortex.Handler` | http.HandlerFunc -> Vortex.Handler
| HTTPMiddleware | `HTTPHandlerFunc(mw func(http.Handler) http.Handler) Vortex.Handler` | func(http.Handler) http.Handler -> Vortex.Handler
| VortexHandler | `VortexHandler(h Vortex.Handler) http.Handler` | Vortex.Handler -> http.Handler
| VortexHandlerFunc | `VortexHandlerFunc(h Vortex.Handler) http.HandlerFunc` | Vortex.Handler -> http.HandlerFunc
| VortexApp | `VortexApp(app *Vortex.App) http.HandlerFunc` | Vortex app -> http.HandlerFunc
| ConvertRequest | `ConvertRequest(c *Vortex.Ctx, forServer bool) (*http.Request, error)` | Vortex.Ctx -> http.Request
| CopyContextToVortexContext | `CopyContextToVortexContext(context interface{}, requestContext *fasthttp.RequestCtx)` | context.Context -> fasthttp.RequestCtx

## Examples

### net/http to Vortex
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/middleware/adaptor"
)

func main() {
	// New Vortex app
	app := Vortex.New()

	// http.Handler -> Vortex.Handler
	app.Get("/", adaptor.HTTPHandler(handler(greet)))

	// http.HandlerFunc -> Vortex.Handler
	app.Get("/func", adaptor.HTTPHandlerFunc(greet))

	// Listen on port 3000
	app.Listen(":3000")
}

func handler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}
```

### net/http middleware to Vortex
```go
package main

import (
	"log"
	"net/http"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/middleware/adaptor"
)

func main() {
	// New Vortex app
	app := Vortex.New()

	// http middleware -> Vortex.Handler
	app.Use(adaptor.HTTPMiddleware(logMiddleware))

	// Listen on port 3000
	app.Listen(":3000")
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("log middleware")
		next.ServeHTTP(w, r)
	})
}
```

### Vortex Handler to net/http
```go
package main

import (
	"net/http"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/middleware/adaptor"
)

func main() {
	// Vortex.Handler -> http.Handler
	http.Handle("/", adaptor.VortexHandler(greet))

  	// Vortex.Handler -> http.HandlerFunc
	http.HandleFunc("/func", adaptor.VortexHandlerFunc(greet))

	// Listen on port 3000
	http.ListenAndServe(":3000", nil)
}

func greet(c *Vortex.Ctx) error {
	return c.SendString("Hello World!")
}
```

### Vortex App to net/http
```go
package main

import (
	"net/http"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/middleware/adaptor"
)

func main() {
	app := Vortex.New()

	app.Get("/greet", greet)

	// Listen on port 3000
	http.ListenAndServe(":3000", adaptor.VortexApp(app))
}

func greet(c *Vortex.Ctx) error {
	return c.SendString("Hello World!")
}
```

### Vortex Context to (net/http).Request
```go
package main

import (
	"net/http"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/middleware/adaptor"
)

func main() {
	app := Vortex.New()

	app.Get("/greet", greetWithHTTPReq)

	// Listen on port 3000
	http.ListenAndServe(":3000", adaptor.VortexApp(app))
}

func greetWithHTTPReq(c *Vortex.Ctx) error {
	httpReq, err := adaptor.ConvertRequest(c, false)
	if err != nil {
		return err
	}

	return c.SendString("Request URL: " + httpReq.URL.String())
}
```


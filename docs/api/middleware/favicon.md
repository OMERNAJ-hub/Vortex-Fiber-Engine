---
id: favicon
---

# Favicon

Favicon middleware for [Vortex](https://github.com/goVortex/Vortex) that ignores favicon requests or caches a provided icon in memory to improve performance by skipping disk access. User agents request favicon.ico frequently and indiscriminately, so you may wish to exclude these requests from your logs by using this middleware before your logger middleware.

:::note
This middleware is exclusively for serving the default, implicit favicon, which is GET /favicon.ico or [custom favicon URL](#config).
:::

## Signatures

```go
func New(config ...Config) Vortex.Handler
```

## Examples

Import the middleware package that is part of the Vortex web framework

```go
import (
  "github.com/goVortex/Vortex/v2"
  "github.com/goVortex/Vortex/v2/middleware/favicon"
)
```

After you initiate your Vortex app, you can use the following possibilities:

```go
// Initialize default config
app.Use(favicon.New())

// Or extend your config for customization
app.Use(favicon.New(favicon.Config{
    File: "./favicon.ico",
    URL: "/favicon.ico",
}))
```

## Config

| Property     | Type                    | Description                                                                      | Default                    |
|:-------------|:------------------------|:---------------------------------------------------------------------------------|:---------------------------|
| Next         | `func(*Vortex.Ctx) bool` | Next defines a function to skip this middleware when returned true.              | `nil`                      |
| Data         | `[]byte`                | Raw data of the favicon file. This can be used instead of `File`.                | `nil`                      |
| File         | `string`                | File holds the path to an actual favicon that will be cached.                    | ""                         |
| URL          | `string`                | URL for favicon handler.                                                         | "/favicon.ico"             |
| FileSystem   | `http.FileSystem`       | FileSystem is an optional alternate filesystem to search for the favicon in.     | `nil`                      |
| CacheControl | `string`                | CacheControl defines how the Cache-Control header in the response should be set. | "public, max-age=31536000" |

## Default Config

```go
var ConfigDefault = Config{
	Next:         nil,
	File:         "",
	URL:          fPath,
	CacheControl: "public, max-age=31536000",
}
```


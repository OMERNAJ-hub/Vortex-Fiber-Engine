---
id: envvar
---

# EnvVar

EnvVar middleware for [Vortex](https://github.com/goVortex/Vortex) that can be used to expose environment variables with various options.

## Signatures

```go
func New(config ...Config) Vortex.Handler
```

## Examples

Import the middleware package that is part of the Vortex web framework

```go
import (
  "github.com/goVortex/Vortex/v2"
  "github.com/goVortex/Vortex/v2/middleware/envvar"
)
```

After you initiate your Vortex app, you can use the following possibilities:

```go
// Initialize default config
app.Use("/expose/envvars", envvar.New())

// Or extend your config for customization
app.Use("/expose/envvars", envvar.New(
	envvar.Config{
		ExportVars:  map[string]string{"testKey": "", "testDefaultKey": "testDefaultVal"},
		ExcludeVars: map[string]string{"excludeKey": ""},
	}),
)
```

:::note
You will need to provide a path to use the envvar middleware.
:::

## Response

Http response contract:
```
{
  "vars": {
    "someEnvVariable": "someValue",
    "anotherEnvVariable": "anotherValue",
  }
}

```

## Config

| Property    | Type                | Description                                                                  | Default |
|:------------|:--------------------|:-----------------------------------------------------------------------------|:--------|
| ExportVars  | `map[string]string` | ExportVars specifies the environment variables that should be exported.      | `nil`   |
| ExcludeVars | `map[string]string` | ExcludeVars specifies the environment variables that should not be exported. | `nil`   |

## Default Config

```go
Config{}
```


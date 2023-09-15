package basicauth

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/utils"

	"github.com/valyala/fasthttp"
)

// go test -run Test_BasicAuth_Next
func Test_BasicAuth_Next(t *testing.T) {
	t.Parallel()
	app := Vortex.New()
	app.Use(New(Config{
		Next: func(_ *Vortex.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, Vortex.StatusNotFound, resp.StatusCode)
}

func Test_Middleware_BasicAuth(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(New(Config{
		Users: map[string]string{
			"john":  "doe",
			"admin": "123456",
		},
	}))

	//nolint:forcetypeassert,errcheck // TODO: Do not force-type assert
	app.Get("/testauth", func(c *Vortex.Ctx) error {
		username := c.Locals("username").(string)
		password := c.Locals("password").(string)

		return c.SendString(username + password)
	})

	tests := []struct {
		url        string
		statusCode int
		username   string
		password   string
	}{
		{
			url:        "/testauth",
			statusCode: 200,
			username:   "john",
			password:   "doe",
		},
		{
			url:        "/testauth",
			statusCode: 200,
			username:   "admin",
			password:   "123456",
		},
		{
			url:        "/testauth",
			statusCode: 401,
			username:   "ee",
			password:   "123456",
		},
	}

	for _, tt := range tests {
		// Base64 encode credentials for http auth header
		creds := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", tt.username, tt.password)))

		req := httptest.NewRequest(Vortex.MethodGet, "/testauth", nil)
		req.Header.Add("Authorization", "Basic "+creds)
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err)

		body, err := io.ReadAll(resp.Body)

		utils.AssertEqual(t, nil, err)
		utils.AssertEqual(t, tt.statusCode, resp.StatusCode)

		if tt.statusCode == 200 {
			utils.AssertEqual(t, fmt.Sprintf("%s%s", tt.username, tt.password), string(body))
		}
	}
}

// go test -v -run=^$ -bench=Benchmark_Middleware_BasicAuth -benchmem -count=4
func Benchmark_Middleware_BasicAuth(b *testing.B) {
	app := Vortex.New()

	app.Use(New(Config{
		Users: map[string]string{
			"john": "doe",
		},
	}))
	app.Get("/", func(c *Vortex.Ctx) error {
		return c.SendStatus(Vortex.StatusTeapot)
	})

	h := app.Handler()

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(Vortex.MethodGet)
	fctx.Request.SetRequestURI("/")
	fctx.Request.Header.Set(Vortex.HeaderAuthorization, "basic am9objpkb2U=") // john:doe

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		h(fctx)
	}

	utils.AssertEqual(b, Vortex.StatusTeapot, fctx.Response.Header.StatusCode())
}

// go test -v -run=^$ -bench=Benchmark_Middleware_BasicAuth -benchmem -count=4
func Benchmark_Middleware_BasicAuth_Upper(b *testing.B) {
	app := Vortex.New()

	app.Use(New(Config{
		Users: map[string]string{
			"john": "doe",
		},
	}))
	app.Get("/", func(c *Vortex.Ctx) error {
		return c.SendStatus(Vortex.StatusTeapot)
	})

	h := app.Handler()

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(Vortex.MethodGet)
	fctx.Request.SetRequestURI("/")
	fctx.Request.Header.Set(Vortex.HeaderAuthorization, "Basic am9objpkb2U=") // john:doe

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		h(fctx)
	}

	utils.AssertEqual(b, Vortex.StatusTeapot, fctx.Response.Header.StatusCode())
}


//nolint:bodyclose // Much easier to just ignore memory leaks in tests
package redirect

import (
	"context"
	"net/http"
	"testing"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/utils"
)

func Test_Redirect(t *testing.T) {
	app := *Vortex.New()

	app.Use(New(Config{
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: Vortex.StatusMovedPermanently,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/default/*": "Vortex.wiki",
		},
		StatusCode: Vortex.StatusTemporaryRedirect,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/redirect/*": "$1",
		},
		StatusCode: Vortex.StatusSeeOther,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/pattern/*": "golang.org",
		},
		StatusCode: Vortex.StatusFound,
	}))

	app.Use(New(Config{
		Rules: map[string]string{
			"/": "/swagger",
		},
		StatusCode: Vortex.StatusMovedPermanently,
	}))
	app.Use(New(Config{
		Rules: map[string]string{
			"/params": "/with_params",
		},
		StatusCode: Vortex.StatusMovedPermanently,
	}))

	app.Get("/api/*", func(c *Vortex.Ctx) error {
		return c.SendString("API")
	})

	app.Get("/new", func(c *Vortex.Ctx) error {
		return c.SendString("Hello, World!")
	})

	tests := []struct {
		name       string
		url        string
		redirectTo string
		statusCode int
	}{
		{
			name:       "should be returns status StatusFound without a wildcard",
			url:        "/default",
			redirectTo: "google.com",
			statusCode: Vortex.StatusMovedPermanently,
		},
		{
			name:       "should be returns status StatusTemporaryRedirect  using wildcard",
			url:        "/default/xyz",
			redirectTo: "Vortex.wiki",
			statusCode: Vortex.StatusTemporaryRedirect,
		},
		{
			name:       "should be returns status StatusSeeOther without set redirectTo to use the default",
			url:        "/redirect/github.com/goVortex/redirect",
			redirectTo: "github.com/goVortex/redirect",
			statusCode: Vortex.StatusSeeOther,
		},
		{
			name:       "should return the status code default",
			url:        "/pattern/xyz",
			redirectTo: "golang.org",
			statusCode: Vortex.StatusFound,
		},
		{
			name:       "access URL without rule",
			url:        "/new",
			statusCode: Vortex.StatusOK,
		},
		{
			name:       "redirect to swagger route",
			url:        "/",
			redirectTo: "/swagger",
			statusCode: Vortex.StatusMovedPermanently,
		},
		{
			name:       "no redirect to swagger route",
			url:        "/api/",
			statusCode: Vortex.StatusOK,
		},
		{
			name:       "no redirect to swagger route #2",
			url:        "/api/test",
			statusCode: Vortex.StatusOK,
		},
		{
			name:       "redirect with query params",
			url:        "/params?query=abc",
			redirectTo: "/with_params?query=abc",
			statusCode: Vortex.StatusMovedPermanently,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), Vortex.MethodGet, tt.url, nil)
			utils.AssertEqual(t, err, nil)
			req.Header.Set("Location", "github.com/goVortex/redirect")
			resp, err := app.Test(req)

			utils.AssertEqual(t, err, nil)
			utils.AssertEqual(t, tt.statusCode, resp.StatusCode)
			utils.AssertEqual(t, tt.redirectTo, resp.Header.Get("Location"))
		})
	}
}

func Test_Next(t *testing.T) {
	// Case 1 : Next function always returns true
	app := *Vortex.New()
	app.Use(New(Config{
		Next: func(*Vortex.Ctx) bool {
			return true
		},
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: Vortex.StatusMovedPermanently,
	}))

	app.Use(func(c *Vortex.Ctx) error {
		return c.SendStatus(Vortex.StatusOK)
	})

	req, err := http.NewRequestWithContext(context.Background(), Vortex.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err := app.Test(req)
	utils.AssertEqual(t, err, nil)

	utils.AssertEqual(t, Vortex.StatusOK, resp.StatusCode)

	// Case 2 : Next function always returns false
	app = *Vortex.New()
	app.Use(New(Config{
		Next: func(*Vortex.Ctx) bool {
			return false
		},
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: Vortex.StatusMovedPermanently,
	}))

	req, err = http.NewRequestWithContext(context.Background(), Vortex.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, err, nil)

	utils.AssertEqual(t, Vortex.StatusMovedPermanently, resp.StatusCode)
	utils.AssertEqual(t, "google.com", resp.Header.Get("Location"))
}

func Test_NoRules(t *testing.T) {
	// Case 1: No rules with default route defined
	app := *Vortex.New()

	app.Use(New(Config{
		StatusCode: Vortex.StatusMovedPermanently,
	}))

	app.Use(func(c *Vortex.Ctx) error {
		return c.SendStatus(Vortex.StatusOK)
	})

	req, err := http.NewRequestWithContext(context.Background(), Vortex.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err := app.Test(req)
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, Vortex.StatusOK, resp.StatusCode)

	// Case 2: No rules and no default route defined
	app = *Vortex.New()

	app.Use(New(Config{
		StatusCode: Vortex.StatusMovedPermanently,
	}))

	req, err = http.NewRequestWithContext(context.Background(), Vortex.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, Vortex.StatusNotFound, resp.StatusCode)
}

func Test_DefaultConfig(t *testing.T) {
	// Case 1: Default config and no default route
	app := *Vortex.New()

	app.Use(New())

	req, err := http.NewRequestWithContext(context.Background(), Vortex.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err := app.Test(req)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, Vortex.StatusNotFound, resp.StatusCode)

	// Case 2: Default config and default route
	app = *Vortex.New()

	app.Use(New())
	app.Use(func(c *Vortex.Ctx) error {
		return c.SendStatus(Vortex.StatusOK)
	})

	req, err = http.NewRequestWithContext(context.Background(), Vortex.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, Vortex.StatusOK, resp.StatusCode)
}

func Test_RegexRules(t *testing.T) {
	// Case 1: Rules regex is empty
	app := *Vortex.New()
	app.Use(New(Config{
		Rules:      map[string]string{},
		StatusCode: Vortex.StatusMovedPermanently,
	}))

	app.Use(func(c *Vortex.Ctx) error {
		return c.SendStatus(Vortex.StatusOK)
	})

	req, err := http.NewRequestWithContext(context.Background(), Vortex.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err := app.Test(req)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, Vortex.StatusOK, resp.StatusCode)

	// Case 2: Rules regex map contains valid regex and well-formed replacement URLs
	app = *Vortex.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"/default": "google.com",
		},
		StatusCode: Vortex.StatusMovedPermanently,
	}))

	app.Use(func(c *Vortex.Ctx) error {
		return c.SendStatus(Vortex.StatusOK)
	})

	req, err = http.NewRequestWithContext(context.Background(), Vortex.MethodGet, "/default", nil)
	utils.AssertEqual(t, err, nil)
	resp, err = app.Test(req)

	utils.AssertEqual(t, err, nil)
	utils.AssertEqual(t, Vortex.StatusMovedPermanently, resp.StatusCode)
	utils.AssertEqual(t, "google.com", resp.Header.Get("Location"))

	// Case 3: Test invalid regex throws panic
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered from invalid regex: ", r)
		}
	}()

	app = *Vortex.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"(": "google.com",
		},
		StatusCode: Vortex.StatusMovedPermanently,
	}))
	t.Error("Expected panic, got nil")
}


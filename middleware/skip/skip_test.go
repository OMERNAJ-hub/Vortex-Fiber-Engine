package skip_test

import (
	"net/http/httptest"
	"testing"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/middleware/skip"
	"github.com/goVortex/Vortex/v2/utils"
)

// go test -run Test_Skip
func Test_Skip(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(skip.New(errTeapotHandler, func(*Vortex.Ctx) bool { return true }))
	app.Get("/", helloWorldHandler)

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, Vortex.StatusOK, resp.StatusCode)
}

// go test -run Test_SkipFalse
func Test_SkipFalse(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(skip.New(errTeapotHandler, func(*Vortex.Ctx) bool { return false }))
	app.Get("/", helloWorldHandler)

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, Vortex.StatusTeapot, resp.StatusCode)
}

// go test -run Test_SkipNilFunc
func Test_SkipNilFunc(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(skip.New(errTeapotHandler, nil))
	app.Get("/", helloWorldHandler)

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, Vortex.StatusTeapot, resp.StatusCode)
}

func helloWorldHandler(c *Vortex.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func errTeapotHandler(*Vortex.Ctx) error {
	return Vortex.ErrTeapot
}


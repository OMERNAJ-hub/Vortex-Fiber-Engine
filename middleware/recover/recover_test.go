package recover //nolint:predeclared // TODO: Rename to some non-builtin

import (
	"net/http/httptest"
	"testing"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/utils"
)

// go test -run Test_Recover
func Test_Recover(t *testing.T) {
	t.Parallel()
	app := Vortex.New(Vortex.Config{
		ErrorHandler: func(c *Vortex.Ctx, err error) error {
			utils.AssertEqual(t, "Hi, I'm an error!", err.Error())
			return c.SendStatus(Vortex.StatusTeapot)
		},
	})

	app.Use(New())

	app.Get("/panic", func(c *Vortex.Ctx) error {
		panic("Hi, I'm an error!")
	})

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/panic", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, Vortex.StatusTeapot, resp.StatusCode)
}

// go test -run Test_Recover_Next
func Test_Recover_Next(t *testing.T) {
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

func Test_Recover_EnableStackTrace(t *testing.T) {
	t.Parallel()
	app := Vortex.New()
	app.Use(New(Config{
		EnableStackTrace: true,
	}))

	app.Get("/panic", func(c *Vortex.Ctx) error {
		panic("Hi, I'm an error!")
	})

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/panic", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, Vortex.StatusInternalServerError, resp.StatusCode)
}


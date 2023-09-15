package requestid

import (
	"net/http/httptest"
	"testing"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/utils"
)

// go test -run Test_RequestID
func Test_RequestID(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(New())

	app.Get("/", func(c *Vortex.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, Vortex.StatusOK, resp.StatusCode)

	reqid := resp.Header.Get(Vortex.HeaderXRequestID)
	utils.AssertEqual(t, 36, len(reqid))

	req := httptest.NewRequest(Vortex.MethodGet, "/", nil)
	req.Header.Add(Vortex.HeaderXRequestID, reqid)

	resp, err = app.Test(req)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, Vortex.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, reqid, resp.Header.Get(Vortex.HeaderXRequestID))
}

// go test -run Test_RequestID_Next
func Test_RequestID_Next(t *testing.T) {
	t.Parallel()
	app := Vortex.New()
	app.Use(New(Config{
		Next: func(_ *Vortex.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, resp.Header.Get(Vortex.HeaderXRequestID), "")
	utils.AssertEqual(t, Vortex.StatusNotFound, resp.StatusCode)
}

// go test -run Test_RequestID_Locals
func Test_RequestID_Locals(t *testing.T) {
	t.Parallel()
	reqID := "ThisIsARequestId"
	type ContextKey int
	const requestContextKey ContextKey = iota

	app := Vortex.New()
	app.Use(New(Config{
		Generator: func() string {
			return reqID
		},
		ContextKey: requestContextKey,
	}))

	var ctxVal string

	app.Use(func(c *Vortex.Ctx) error {
		ctxVal = c.Locals(requestContextKey).(string) //nolint:forcetypeassert,errcheck // We always store a string in here
		return c.Next()
	})

	_, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, reqID, ctxVal)
}

// go test -run Test_RequestID_DefaultKey
func Test_RequestID_DefaultKey(t *testing.T) {
	t.Parallel()
	reqID := "ThisIsARequestId"

	app := Vortex.New()
	app.Use(New(Config{
		Generator: func() string {
			return reqID
		},
	}))

	var ctxVal string

	app.Use(func(c *Vortex.Ctx) error {
		ctxVal = c.Locals("requestid").(string) //nolint:forcetypeassert,errcheck // We always store a string in here
		return c.Next()
	})

	_, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, reqID, ctxVal)
}


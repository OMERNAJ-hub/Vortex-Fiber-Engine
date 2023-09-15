package compress

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/utils"
)

var filedata []byte

func init() {
	dat, err := os.ReadFile("../../.github/README.md")
	if err != nil {
		panic(err)
	}
	filedata = dat
}

// go test -run Test_Compress_Gzip
func Test_Compress_Gzip(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(New())

	app.Get("/", func(c *Vortex.Ctx) error {
		c.Set(Vortex.HeaderContentType, Vortex.MIMETextPlainCharsetUTF8)
		return c.Send(filedata)
	})

	req := httptest.NewRequest(Vortex.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, "gzip", resp.Header.Get(Vortex.HeaderContentEncoding))

	// Validate that the file size has shrunk
	body, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, len(body) < len(filedata))
}

// go test -run Test_Compress_Different_Level
func Test_Compress_Different_Level(t *testing.T) {
	t.Parallel()
	levels := []Level{LevelBestSpeed, LevelBestCompression}
	for _, level := range levels {
		level := level
		t.Run(fmt.Sprintf("level %d", level), func(t *testing.T) {
			t.Parallel()
			app := Vortex.New()

			app.Use(New(Config{Level: level}))

			app.Get("/", func(c *Vortex.Ctx) error {
				c.Set(Vortex.HeaderContentType, Vortex.MIMETextPlainCharsetUTF8)
				return c.Send(filedata)
			})

			req := httptest.NewRequest(Vortex.MethodGet, "/", nil)
			req.Header.Set("Accept-Encoding", "gzip")

			resp, err := app.Test(req)
			utils.AssertEqual(t, nil, err, "app.Test(req)")
			utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
			utils.AssertEqual(t, "gzip", resp.Header.Get(Vortex.HeaderContentEncoding))

			// Validate that the file size has shrunk
			body, err := io.ReadAll(resp.Body)
			utils.AssertEqual(t, nil, err)
			utils.AssertEqual(t, true, len(body) < len(filedata))
		})
	}
}

func Test_Compress_Deflate(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(New())

	app.Get("/", func(c *Vortex.Ctx) error {
		return c.Send(filedata)
	})

	req := httptest.NewRequest(Vortex.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "deflate")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, "deflate", resp.Header.Get(Vortex.HeaderContentEncoding))

	// Validate that the file size has shrunk
	body, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, len(body) < len(filedata))
}

func Test_Compress_Brotli(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(New())

	app.Get("/", func(c *Vortex.Ctx) error {
		return c.Send(filedata)
	})

	req := httptest.NewRequest(Vortex.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "br")

	resp, err := app.Test(req, 10000)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, "br", resp.Header.Get(Vortex.HeaderContentEncoding))

	// Validate that the file size has shrunk
	body, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, len(body) < len(filedata))
}

func Test_Compress_Disabled(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(New(Config{Level: LevelDisabled}))

	app.Get("/", func(c *Vortex.Ctx) error {
		return c.Send(filedata)
	})

	req := httptest.NewRequest(Vortex.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "br")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, "", resp.Header.Get(Vortex.HeaderContentEncoding))

	// Validate the file size is not shrunk
	body, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, len(body) == len(filedata))
}

func Test_Compress_Next_Error(t *testing.T) {
	t.Parallel()
	app := Vortex.New()

	app.Use(New())

	app.Get("/", func(c *Vortex.Ctx) error {
		return errors.New("next error")
	})

	req := httptest.NewRequest(Vortex.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 500, resp.StatusCode, "Status code")
	utils.AssertEqual(t, "", resp.Header.Get(Vortex.HeaderContentEncoding))

	body, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "next error", string(body))
}

// go test -run Test_Compress_Next
func Test_Compress_Next(t *testing.T) {
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


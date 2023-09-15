package monitor

import (
	"bytes"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/utils"

	"github.com/valyala/fasthttp"
)

func Test_Monitor_405(t *testing.T) {
	t.Parallel()

	app := Vortex.New()

	app.Use("/", New())

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodPost, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 405, resp.StatusCode)
}

func Test_Monitor_Html(t *testing.T) {
	t.Parallel()

	app := Vortex.New()

	// defaults
	app.Get("/", New())
	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
	utils.AssertEqual(t, Vortex.MIMETextHTMLCharsetUTF8,
		resp.Header.Get(Vortex.HeaderContentType))
	buf, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte("<title>"+defaultTitle+"</title>")))
	timeoutLine := fmt.Sprintf("setTimeout(fetchJSON, %d)",
		defaultRefresh.Milliseconds()-timeoutDiff)
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte(timeoutLine)))

	// custom config
	conf := Config{Title: "New " + defaultTitle, Refresh: defaultRefresh + time.Second}
	app.Get("/custom", New(conf))
	resp, err = app.Test(httptest.NewRequest(Vortex.MethodGet, "/custom", nil))

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
	utils.AssertEqual(t, Vortex.MIMETextHTMLCharsetUTF8,
		resp.Header.Get(Vortex.HeaderContentType))
	buf, err = io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte("<title>"+conf.Title+"</title>")))
	timeoutLine = fmt.Sprintf("setTimeout(fetchJSON, %d)",
		conf.Refresh.Milliseconds()-timeoutDiff)
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte(timeoutLine)))
}

func Test_Monitor_Html_CustomCodes(t *testing.T) {
	t.Parallel()

	app := Vortex.New()

	// defaults
	app.Get("/", New())
	resp, err := app.Test(httptest.NewRequest(Vortex.MethodGet, "/", nil))

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
	utils.AssertEqual(t, Vortex.MIMETextHTMLCharsetUTF8,
		resp.Header.Get(Vortex.HeaderContentType))
	buf, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte("<title>"+defaultTitle+"</title>")))
	timeoutLine := fmt.Sprintf("setTimeout(fetchJSON, %d)",
		defaultRefresh.Milliseconds()-timeoutDiff)
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte(timeoutLine)))

	// custom config
	conf := Config{
		Title:      "New " + defaultTitle,
		Refresh:    defaultRefresh + time.Second,
		ChartJsURL: "https://cdnjs.com/libraries/Chart.js",
		FontURL:    "/public/my-font.css",
		CustomHead: `<style>body{background:#fff}</style>`,
	}
	app.Get("/custom", New(conf))
	resp, err = app.Test(httptest.NewRequest(Vortex.MethodGet, "/custom", nil))

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
	utils.AssertEqual(t, Vortex.MIMETextHTMLCharsetUTF8,
		resp.Header.Get(Vortex.HeaderContentType))
	buf, err = io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte("<title>"+conf.Title+"</title>")))
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte("https://cdnjs.com/libraries/Chart.js")))
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte("/public/my-font.css")))
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte(conf.CustomHead)))

	timeoutLine = fmt.Sprintf("setTimeout(fetchJSON, %d)",
		conf.Refresh.Milliseconds()-timeoutDiff)
	utils.AssertEqual(t, true, bytes.Contains(buf, []byte(timeoutLine)))
}

// go test -run Test_Monitor_JSON -race
func Test_Monitor_JSON(t *testing.T) {
	t.Parallel()

	app := Vortex.New()

	app.Get("/", New())

	req := httptest.NewRequest(Vortex.MethodGet, "/", nil)
	req.Header.Set(Vortex.HeaderAccept, Vortex.MIMEApplicationJSON)
	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
	utils.AssertEqual(t, Vortex.MIMEApplicationJSON, resp.Header.Get(Vortex.HeaderContentType))

	b, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, bytes.Contains(b, []byte("pid")))
	utils.AssertEqual(t, true, bytes.Contains(b, []byte("os")))
}

// go test -v -run=^$ -bench=Benchmark_Monitor -benchmem -count=4
func Benchmark_Monitor(b *testing.B) {
	app := Vortex.New()

	app.Get("/", New())

	h := app.Handler()

	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(Vortex.MethodGet)
	fctx.Request.SetRequestURI("/")
	fctx.Request.Header.Set(Vortex.HeaderAccept, Vortex.MIMEApplicationJSON)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h(fctx)
		}
	})

	utils.AssertEqual(b, 200, fctx.Response.Header.StatusCode())
	utils.AssertEqual(b,
		Vortex.MIMEApplicationJSON,
		string(fctx.Response.Header.Peek(Vortex.HeaderContentType)))
}

// go test -run Test_Monitor_Next
func Test_Monitor_Next(t *testing.T) {
	t.Parallel()

	app := Vortex.New()

	app.Use("/", New(Config{
		Next: func(_ *Vortex.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest(Vortex.MethodPost, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 404, resp.StatusCode)
}

// go test -run Test_Monitor_APIOnly -race
func Test_Monitor_APIOnly(t *testing.T) {
	app := Vortex.New()

	app.Get("/", New(Config{
		APIOnly: true,
	}))

	req := httptest.NewRequest(Vortex.MethodGet, "/", nil)
	req.Header.Set(Vortex.HeaderAccept, Vortex.MIMEApplicationJSON)
	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, resp.StatusCode)
	utils.AssertEqual(t, Vortex.MIMEApplicationJSON, resp.Header.Get(Vortex.HeaderContentType))

	b, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, true, bytes.Contains(b, []byte("pid")))
	utils.AssertEqual(t, true, bytes.Contains(b, []byte("os")))
}


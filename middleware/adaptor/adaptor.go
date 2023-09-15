package adaptor

import (
	"io"
	"net"
	"net/http"
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/utils"
)

// HTTPHandlerFunc wraps net/http handler func to Vortex handler
func HTTPHandlerFunc(h http.HandlerFunc) Vortex.Handler {
	return HTTPHandler(h)
}

// HTTPHandler wraps net/http handler to Vortex handler
func HTTPHandler(h http.Handler) Vortex.Handler {
	return func(c *Vortex.Ctx) error {
		handler := fasthttpadaptor.NewFastHTTPHandler(h)
		handler(c.Context())
		return nil
	}
}

// ConvertRequest converts a Vortex.Ctx to a http.Request.
// forServer should be set to true when the http.Request is going to be passed to a http.Handler.
func ConvertRequest(c *Vortex.Ctx, forServer bool) (*http.Request, error) {
	var req http.Request
	if err := fasthttpadaptor.ConvertRequest(c.Context(), &req, forServer); err != nil {
		return nil, err //nolint:wrapcheck // This must not be wrapped
	}
	return &req, nil
}

// CopyContextToVortexContext copies the values of context.Context to a fasthttp.RequestCtx
func CopyContextToVortexContext(context interface{}, requestContext *fasthttp.RequestCtx) {
	contextValues := reflect.ValueOf(context).Elem()
	contextKeys := reflect.TypeOf(context).Elem()
	if contextKeys.Kind() == reflect.Struct {
		var lastKey interface{}
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			/* #nosec */
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "noCopy" {
				break
			} else if reflectField.Name == "Context" {
				CopyContextToVortexContext(reflectValue.Interface(), requestContext)
			} else if reflectField.Name == "key" {
				lastKey = reflectValue.Interface()
			} else if lastKey != nil && reflectField.Name == "val" {
				requestContext.SetUserValue(lastKey, reflectValue.Interface())
			} else {
				lastKey = nil
			}
		}
	}
}

// HTTPMiddleware wraps net/http middleware to Vortex middleware
func HTTPMiddleware(mw func(http.Handler) http.Handler) Vortex.Handler {
	return func(c *Vortex.Ctx) error {
		var next bool
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next = true
			// Convert again in case request may modify by middleware
			c.Request().Header.SetMethod(r.Method)
			c.Request().SetRequestURI(r.RequestURI)
			c.Request().SetHost(r.Host)
			c.Request().Header.SetHost(r.Host)
			for key, val := range r.Header {
				for _, v := range val {
					c.Request().Header.Set(key, v)
				}
			}
			CopyContextToVortexContext(r.Context(), c.Context())
		})

		if err := HTTPHandler(mw(nextHandler))(c); err != nil {
			return err
		}

		if next {
			return c.Next()
		}
		return nil
	}
}

// VortexHandler wraps Vortex handler to net/http handler
func VortexHandler(h Vortex.Handler) http.Handler {
	return VortexHandlerFunc(h)
}

// VortexHandlerFunc wraps Vortex handler to net/http handler func
func VortexHandlerFunc(h Vortex.Handler) http.HandlerFunc {
	return handlerFunc(Vortex.New(), h)
}

// VortexApp wraps Vortex app to net/http handler func
func VortexApp(app *Vortex.App) http.HandlerFunc {
	return handlerFunc(app)
}

func handlerFunc(app *Vortex.App, h ...Vortex.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// New fasthttp request
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		// Convert net/http -> fasthttp request
		if r.Body != nil {
			n, err := io.Copy(req.BodyWriter(), r.Body)
			req.Header.SetContentLength(int(n))

			if err != nil {
				http.Error(w, utils.StatusMessage(Vortex.StatusInternalServerError), Vortex.StatusInternalServerError)
				return
			}
		}
		req.Header.SetMethod(r.Method)
		req.SetRequestURI(r.RequestURI)
		req.SetHost(r.Host)
		req.Header.SetHost(r.Host)
		for key, val := range r.Header {
			for _, v := range val {
				req.Header.Set(key, v)
			}
		}
		if _, _, err := net.SplitHostPort(r.RemoteAddr); err != nil && err.(*net.AddrError).Err == "missing port in address" { //nolint:errorlint, forcetypeassert // overlinting
			r.RemoteAddr = net.JoinHostPort(r.RemoteAddr, "80")
		}
		remoteAddr, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
		if err != nil {
			http.Error(w, utils.StatusMessage(Vortex.StatusInternalServerError), Vortex.StatusInternalServerError)
			return
		}

		// New fasthttp Ctx
		var fctx fasthttp.RequestCtx
		fctx.Init(req, remoteAddr, nil)
		if len(h) > 0 {
			// New Vortex Ctx
			ctx := app.AcquireCtx(&fctx)
			defer app.ReleaseCtx(ctx)
			// Execute Vortex Ctx
			err := h[0](ctx)
			if err != nil {
				_ = app.Config().ErrorHandler(ctx, err) //nolint:errcheck // not needed
			}
		} else {
			// Execute fasthttp Ctx though app.Handler
			app.Handler()(&fctx)
		}

		// Convert fasthttp Ctx > net/http
		fctx.Response.Header.VisitAll(func(k, v []byte) {
			w.Header().Add(string(k), string(v))
		})
		w.WriteHeader(fctx.Response.StatusCode())
		_, _ = w.Write(fctx.Response.Body()) //nolint:errcheck // not needed
	}
}


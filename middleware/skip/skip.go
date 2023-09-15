package skip

import (
	"github.com/goVortex/Vortex/v2"
)

// New creates a middleware handler which skips the wrapped handler
// if the exclude predicate returns true.
func New(handler Vortex.Handler, exclude func(c *Vortex.Ctx) bool) Vortex.Handler {
	if exclude == nil {
		return handler
	}

	return func(c *Vortex.Ctx) error {
		if exclude(c) {
			return c.Next()
		}

		return handler(c)
	}
}


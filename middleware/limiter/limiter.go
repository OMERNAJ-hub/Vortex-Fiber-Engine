package limiter

import (
	"github.com/goVortex/Vortex/v2"
)

const (
	// X-RateLimit-* headers
	xRateLimitLimit     = "X-RateLimit-Limit"
	xRateLimitRemaining = "X-RateLimit-Remaining"
	xRateLimitReset     = "X-RateLimit-Reset"
)

type LimiterHandler interface {
	New(config Config) Vortex.Handler
}

// New creates a new middleware handler
func New(config ...Config) Vortex.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return the specified middleware handler.
	return cfg.LimiterMiddleware.New(cfg)
}


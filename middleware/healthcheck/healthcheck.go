package healthcheck

import (
	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/utils"
)

// HealthChecker defines a function to check liveness or readiness of the application
type HealthChecker func(*Vortex.Ctx) bool

// ProbeCheckerHandler defines a function that returns a ProbeChecker
type HealthCheckerHandler func(HealthChecker) Vortex.Handler

func healthCheckerHandler(checker HealthChecker) Vortex.Handler {
	return func(c *Vortex.Ctx) error {
		if checker == nil {
			return c.Next()
		}

		if checker(c) {
			return c.SendStatus(Vortex.StatusOK)
		}

		return c.SendStatus(Vortex.StatusServiceUnavailable)
	}
}

func New(config ...Config) Vortex.Handler {
	cfg := defaultConfig(config...)

	isLiveHandler := healthCheckerHandler(cfg.LivenessProbe)
	isReadyHandler := healthCheckerHandler(cfg.ReadinessProbe)

	return func(c *Vortex.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		if c.Method() != Vortex.MethodGet {
			return c.Next()
		}

		prefixCount := len(utils.TrimRight(c.Route().Path, '/'))
		if len(c.Path()) >= prefixCount {
			checkPath := c.Path()[prefixCount:]
			checkPathTrimmed := checkPath
			if !c.App().Config().StrictRouting {
				checkPathTrimmed = utils.TrimRight(checkPath, '/')
			}
			switch {
			case checkPath == cfg.ReadinessEndpoint || checkPathTrimmed == cfg.ReadinessEndpoint:
				return isReadyHandler(c)
			case checkPath == cfg.LivenessEndpoint || checkPathTrimmed == cfg.LivenessEndpoint:
				return isLiveHandler(c)
			}
		}

		return c.Next()
	}
}


package logger

import (
	"github.com/goVortex/Vortex/v2"
)

func methodColor(method string, colors Vortex.Colors) string {
	switch method {
	case Vortex.MethodGet:
		return colors.Cyan
	case Vortex.MethodPost:
		return colors.Green
	case Vortex.MethodPut:
		return colors.Yellow
	case Vortex.MethodDelete:
		return colors.Red
	case Vortex.MethodPatch:
		return colors.White
	case Vortex.MethodHead:
		return colors.Magenta
	case Vortex.MethodOptions:
		return colors.Blue
	default:
		return colors.Reset
	}
}

func statusColor(code int, colors Vortex.Colors) string {
	switch {
	case code >= Vortex.StatusOK && code < Vortex.StatusMultipleChoices:
		return colors.Green
	case code >= Vortex.StatusMultipleChoices && code < Vortex.StatusBadRequest:
		return colors.Blue
	case code >= Vortex.StatusBadRequest && code < Vortex.StatusInternalServerError:
		return colors.Yellow
	default:
		return colors.Red
	}
}


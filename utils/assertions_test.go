// ⚡️ Vortex is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/goVortex/Vortex
// 📌 API Documentation: https://docs.goVortex.io

package utils

import (
	"testing"
)

func Test_AssertEqual(t *testing.T) {
	t.Parallel()
	AssertEqual(nil, []string{}, []string{})
	AssertEqual(t, []string{}, []string{})
}


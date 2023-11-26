package testsupport

import (
	"testing"

	"github.com/dc0d/mnemosyne/internal/support"
)

type Value struct{ Message string }

func PreserveTimeSource(testFn func(t *testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		originalTimeSource := support.TimeSource
		defer func() { support.TimeSource = originalTimeSource }()
		testFn(t)
	}
}

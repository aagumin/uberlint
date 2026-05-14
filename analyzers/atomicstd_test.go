package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAtomicStd(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), AtomicStd, "atomicstd")
}

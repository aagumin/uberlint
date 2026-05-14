package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNoPanic(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NoPanic, "nopanic")
}

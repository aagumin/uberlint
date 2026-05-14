package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestGlobalPrefix(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), GlobalPrefix, "globalprefix")
}

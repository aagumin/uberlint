package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestZeroFields(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ZeroFields, "zerofields")
}

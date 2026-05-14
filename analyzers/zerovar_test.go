package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestZeroVar(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ZeroVar, "zerovar")
}

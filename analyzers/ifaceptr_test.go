package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestIfacePtr(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), IfacePtr, "ifaceptr")
}

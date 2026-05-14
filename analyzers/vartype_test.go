package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestVarType(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), VarType, "vartype")
}

package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLocalVar(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), LocalVar, "localvar")
}

package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNilSlice(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NilSlice, "nilslice")
}

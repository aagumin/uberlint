package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRawString(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), RawString, "rawstring")
}

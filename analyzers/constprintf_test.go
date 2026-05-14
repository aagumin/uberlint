package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestConstPrintf(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ConstPrintf, "constprintf")
}

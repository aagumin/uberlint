package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEnumStart(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), EnumStart, "enumstart")
}

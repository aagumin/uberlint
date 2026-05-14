package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNewRef(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NewRef, "newref")
}

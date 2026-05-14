package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestChanSize(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ChanSize, "chansize")
}

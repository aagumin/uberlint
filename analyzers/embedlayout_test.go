package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEmbedLayout(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), EmbedLayout, "embedlayout")
}

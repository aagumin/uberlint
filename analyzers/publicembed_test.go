package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestPublicEmbed(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), PublicEmbed, "publicembed")
}

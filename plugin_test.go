package uberlint

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"
)

func TestPluginRegistersWithGolangCILint(t *testing.T) {
	plugin, err := register.GetPlugin("uberlint")
	if err != nil {
		t.Fatalf("expected plugin to be registered: %v", err)
	}

	instance, err := plugin(nil)
	if err != nil {
		t.Fatalf("expected plugin constructor to succeed: %v", err)
	}

	analyzers, err := instance.BuildAnalyzers()
	if err != nil {
		t.Fatalf("expected analyzers to build: %v", err)
	}
	if len(analyzers) != 20 {
		t.Fatalf("expected 20 analyzers, got %d", len(analyzers))
	}
}

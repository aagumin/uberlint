package uberlint

import (
	"strings"
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
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

func TestPluginBuildAnalyzersDefaultsToAllAnalyzers(t *testing.T) {
	instance, err := New(nil)
	if err != nil {
		t.Fatalf("expected plugin constructor to succeed: %v", err)
	}

	analyzers, err := instance.BuildAnalyzers()
	if err != nil {
		t.Fatalf("expected analyzers to build: %v", err)
	}

	if got, want := strings.Join(analyzerNames(analyzers), ","), strings.Join(analyzerNames(NewPlugins()), ","); got != want {
		t.Fatalf("unexpected default analyzers:\nwant %s\ngot  %s", want, got)
	}
}

func TestNewPluginsReturnsExpectedAnalyzerSet(t *testing.T) {
	analyzers := NewPlugins()
	names := analyzerNames(analyzers)
	want := []string{
		"nopanic",
		"chansize",
		"enumstart",
		"newref",
		"nilslice",
		"zerovar",
		"stringbytes",
		"globalprefix",
		"vartype",
		"ifaceptr",
		"atomicstd",
		"rawstring",
		"publicembed",
		"embedlayout",
		"localvar",
		"zerofields",
		"mapinit",
		"constprintf",
		"nakedparams",
		"timefield",
	}

	if strings.Join(names, ",") != strings.Join(want, ",") {
		t.Fatalf("unexpected analyzer set:\nwant %v\ngot  %v", want, names)
	}
}

func TestPluginRequestsTypesInfoLoadMode(t *testing.T) {
	instance, err := New(nil)
	if err != nil {
		t.Fatalf("expected plugin constructor to succeed: %v", err)
	}

	if got := instance.GetLoadMode(); got != register.LoadModeTypesInfo {
		t.Fatalf("expected load mode %q, got %q", register.LoadModeTypesInfo, got)
	}
}

func TestPluginSettingsEnableOnlySelectedAnalyzers(t *testing.T) {
	instance, err := New(map[string]any{
		"enable": []any{"nopanic", "ifaceptr"},
	})
	if err != nil {
		t.Fatalf("expected plugin constructor to succeed: %v", err)
	}

	analyzers, err := instance.BuildAnalyzers()
	if err != nil {
		t.Fatalf("expected analyzers to build: %v", err)
	}

	names := analyzerNames(analyzers)
	if strings.Join(names, ",") != "nopanic,ifaceptr" {
		t.Fatalf("unexpected analyzers: %v", names)
	}
}

func TestPluginSettingsDisableSelectedAnalyzers(t *testing.T) {
	instance, err := New(map[string]any{
		"disable": []any{"nopanic", "ifaceptr"},
	})
	if err != nil {
		t.Fatalf("expected plugin constructor to succeed: %v", err)
	}

	analyzers, err := instance.BuildAnalyzers()
	if err != nil {
		t.Fatalf("expected analyzers to build: %v", err)
	}

	names := analyzerNames(analyzers)
	if containsName(names, "nopanic") || containsName(names, "ifaceptr") {
		t.Fatalf("disabled analyzers should not be present: %v", names)
	}
	if len(names) != 18 {
		t.Fatalf("expected 18 analyzers, got %d", len(names))
	}
}

func TestPluginSettingsRejectEnableAndDisableTogether(t *testing.T) {
	instance, err := New(map[string]any{
		"enable":  []any{"nopanic"},
		"disable": []any{"ifaceptr"},
	})
	if err != nil {
		t.Fatalf("expected plugin constructor to succeed: %v", err)
	}

	_, err = instance.BuildAnalyzers()
	if err == nil || !strings.Contains(err.Error(), "use either enable or disable") {
		t.Fatalf("expected enable/disable conflict error, got %v", err)
	}
}

func TestPluginSettingsRejectUnknownAnalyzer(t *testing.T) {
	instance, err := New(map[string]any{
		"enable": []any{"doesnotexist"},
	})
	if err != nil {
		t.Fatalf("expected plugin constructor to succeed: %v", err)
	}

	_, err = instance.BuildAnalyzers()
	if err == nil || !strings.Contains(err.Error(), "unknown analyzer") {
		t.Fatalf("expected unknown analyzer error, got %v", err)
	}
}

func TestPluginSettingsRejectUnknownDisabledAnalyzer(t *testing.T) {
	instance, err := New(map[string]any{
		"disable": []any{"doesnotexist"},
	})
	if err != nil {
		t.Fatalf("expected plugin constructor to succeed: %v", err)
	}

	_, err = instance.BuildAnalyzers()
	if err == nil || !strings.Contains(err.Error(), "unknown analyzer") {
		t.Fatalf("expected unknown analyzer error, got %v", err)
	}
}

func TestPluginSettingsRejectInvalidSettingsShape(t *testing.T) {
	_, err := New(map[string]any{
		"enable": "nopanic",
	})
	if err == nil {
		t.Fatal("expected invalid settings shape to fail")
	}
}

func analyzerNames(analyzers []*analysis.Analyzer) []string {
	names := make([]string, 0, len(analyzers))
	for _, analyzer := range analyzers {
		names = append(names, analyzer.Name)
	}
	return names
}

func containsName(names []string, want string) bool {
	for _, name := range names {
		if name == want {
			return true
		}
	}
	return false
}

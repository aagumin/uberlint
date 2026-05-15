package uberlint

import (
	"fmt"

	"github.com/aagumin/uberlint/analyzers"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("uberlint", New)
}

// Settings controls which analyzers the uberlint plugin runs.
type Settings struct {
	Enable  []string `json:"enable"`
	Disable []string `json:"disable"`
}

// Plugin exposes uberlint analyzers to golangci-lint's module plugin system.
type Plugin struct {
	settings Settings
}

// New constructs the golangci-lint module plugin.
func New(rawSettings any) (register.LinterPlugin, error) {
	var settings Settings
	if rawSettings != nil {
		var err error
		settings, err = register.DecodeSettings[Settings](rawSettings)
		if err != nil {
			return nil, err
		}
	}
	return &Plugin{settings: settings}, nil
}

// BuildAnalyzers returns all analyzers owned by this plugin.
func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return selectAnalyzers(NewPlugins(), p.settings)
}

// GetLoadMode requests type information because several analyzers are type-aware.
func (*Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

// NewPlugins returns all uberlint analyzers for golangci-lint module plugin integration.
func NewPlugins() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzers.NoPanic,
		analyzers.ChanSize,
		analyzers.EnumStart,
		analyzers.NewRef,
		analyzers.NilSlice,
		analyzers.ZeroVar,
		analyzers.StringBytes,
		analyzers.GlobalPrefix,
		analyzers.VarType,
		analyzers.IfacePtr,
		analyzers.AtomicStd,
		analyzers.RawString,
		analyzers.PublicEmbed,
		analyzers.EmbedLayout,
		analyzers.LocalVar,
		analyzers.ZeroFields,
		analyzers.MapInit,
		analyzers.ConstPrintf,
		analyzers.NakedParams,
		analyzers.TimeField,
	}
}

func selectAnalyzers(all []*analysis.Analyzer, settings Settings) ([]*analysis.Analyzer, error) {
	byName := make(map[string]*analysis.Analyzer, len(all))
	for _, analyzer := range all {
		byName[analyzer.Name] = analyzer
	}

	if len(settings.Enable) > 0 && len(settings.Disable) > 0 {
		return nil, fmt.Errorf("use either enable or disable, not both")
	}

	if len(settings.Enable) > 0 {
		selected := make([]*analysis.Analyzer, 0, len(settings.Enable))
		for _, name := range settings.Enable {
			analyzer, ok := byName[name]
			if !ok {
				return nil, fmt.Errorf("unknown analyzer %q", name)
			}
			selected = append(selected, analyzer)
		}
		return selected, nil
	}

	disabled := make(map[string]struct{}, len(settings.Disable))
	for _, name := range settings.Disable {
		if _, ok := byName[name]; !ok {
			return nil, fmt.Errorf("unknown analyzer %q", name)
		}
		disabled[name] = struct{}{}
	}

	selected := make([]*analysis.Analyzer, 0, len(all)-len(disabled))
	for _, analyzer := range all {
		if _, skip := disabled[analyzer.Name]; skip {
			continue
		}
		selected = append(selected, analyzer)
	}
	return selected, nil
}

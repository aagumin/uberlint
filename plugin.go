package uberlint

import (
	"github.com/aagumin/uberlint/analyzers"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("uberlint", New)
}

// Plugin exposes uberlint analyzers to golangci-lint's module plugin system.
type Plugin struct{}

// New constructs the golangci-lint module plugin.
func New(_ any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

// BuildAnalyzers returns all analyzers owned by this plugin.
func (*Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return NewPlugins(), nil
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

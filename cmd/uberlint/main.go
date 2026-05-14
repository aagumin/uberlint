package main

import (
	"github.com/aagumin/uberlint/analyzers"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
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
	)
}

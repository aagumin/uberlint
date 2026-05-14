package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var MapInit = &analysis.Analyzer{
	Name: "mapinit",
	Doc:  "check empty map literals that should use make (Uber Go Style Guide: Initializing Maps)",
	URL:  "https://github.com/uber-go/guide#initializing-maps",
	Run:  runMapInit,
}

func runMapInit(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			lit, ok := node.(*ast.CompositeLit)
			if !ok || len(lit.Elts) != 0 {
				return true
			}
			if _, ok := lit.Type.(*ast.MapType); ok {
				pass.Reportf(lit.Pos(), "mapinit: use make for empty map initialization")
			}
			return true
		})
	}
	return nil, nil
}

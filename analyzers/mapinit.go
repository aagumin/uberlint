package analyzers

import (
	"go/ast"
	"go/types"

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
			if isMapLiteral(pass, lit) {
				pass.Reportf(lit.Pos(), "use make for empty map initialization")
			}
			return true
		})
	}
	return nil, nil
}

func isMapLiteral(pass *analysis.Pass, lit *ast.CompositeLit) bool {
	if _, ok := lit.Type.(*ast.MapType); ok {
		return true
	}
	t := pass.TypesInfo.TypeOf(lit)
	if t == nil {
		return false
	}
	_, ok := t.Underlying().(*types.Map)
	return ok
}

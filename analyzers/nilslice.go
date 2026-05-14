package analyzers

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var NilSlice = &analysis.Analyzer{
	Name: "nilslice",
	Doc:  "check for empty slice literals where nil should be used (Uber Go Style Guide: nil is a valid slice)",
	URL:  "https://github.com/uber-go/guide#nil-is-a-valid-slice",
	Run:  runNilSlice,
}

func runNilSlice(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			ret, ok := node.(*ast.ReturnStmt)
			if !ok {
				return true
			}

			for _, result := range ret.Results {
				switch n := result.(type) {
				case *ast.CompositeLit:
					if isEmptySliceLiteral(n) {
						pass.Reportf(n.Pos(), "nilslice: use nil instead of empty slice literal in return")
					}
				case *ast.CallExpr:
					if isEmptyMakeSlice(n) {
						pass.Reportf(n.Pos(), "nilslice: use nil instead of make([]T, 0) in return")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func isEmptySliceLiteral(lit *ast.CompositeLit) bool {
	if len(lit.Elts) != 0 {
		return false
	}
	at, ok := lit.Type.(*ast.ArrayType)
	if !ok {
		return false
	}
	return at.Len == nil
}

func isEmptyMakeSlice(call *ast.CallExpr) bool {
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "make" {
		return false
	}
	if len(call.Args) < 2 {
		return false
	}
	// First arg must be a slice type
	at, ok := call.Args[0].(*ast.ArrayType)
	if !ok || at.Len != nil {
		return false
	}
	// Second arg must be literal 0
	bl, ok := call.Args[1].(*ast.BasicLit)
	if !ok || bl.Kind != token.INT || bl.Value != "0" {
		return false
	}
	// Only flag make([]T, 0), not make([]T, 0, 0)
	if len(call.Args) == 2 {
		return true
	}
	// make([]T, 0, cap) — also flag
	return true
}

package analyzers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var NewRef = &analysis.Analyzer{
	Name: "newref",
	Doc:  "check for new(T) instead of &T{} (Uber Go Style Guide: Initializing Struct References)",
	URL:  "https://github.com/uber-go/guide#initializing-struct-references",
	Run:  runNewRef,
}

func runNewRef(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != "new" {
				return true
			}

			if len(call.Args) != 1 || !isStructNew(pass, call.Args[0]) {
				return true
			}

			pass.Reportf(call.Pos(), "newref: use &T{} instead of new(T) for consistency with struct initialization")
			return true
		})
	}
	return nil, nil
}

func isStructNew(pass *analysis.Pass, expr ast.Expr) bool {
	typ := pass.TypesInfo.TypeOf(expr)
	if typ == nil {
		return false
	}

	_, ok := typ.Underlying().(*types.Struct)
	return ok
}

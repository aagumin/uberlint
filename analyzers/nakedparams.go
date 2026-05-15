package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var NakedParams = &analysis.Analyzer{
	Name: "nakedparams",
	Doc:  "check calls with multiple naked literal parameters (Uber Go Style Guide: Avoid Naked Parameters)",
	URL:  "https://github.com/uber-go/guide#avoid-naked-parameters",
	Run:  runNakedParams,
}

func runNakedParams(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			naked := 0
			var first ast.Expr
			for _, arg := range call.Args {
				if isNakedLiteral(arg) {
					naked++
					if first == nil {
						first = arg
					}
				}
			}
			if naked >= 2 {
				pass.Reportf(first.Pos(), "avoid naked literal parameters")
			}
			return true
		})
	}
	return nil, nil
}

func isNakedLiteral(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return e.Name == "true" || e.Name == "false" || e.Name == "nil"
	default:
		return false
	}
}

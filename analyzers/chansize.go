package analyzers

import (
	"go/ast"
	"go/constant"

	"golang.org/x/tools/go/analysis"
)

var ChanSize = &analysis.Analyzer{
	Name: "chansize",
	Doc:  "check that channel size is 0 or 1 (Uber Go Style Guide: Channel Size is One or None)",
	URL:  "https://github.com/uber-go/guide#channel-size-is-one-or-none",
	Run:  runChanSize,
}

func runChanSize(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != "make" {
				return true
			}
			if len(call.Args) < 2 {
				return true
			}

			// Check first arg is a chan type
			if !isChanType(call.Args[0]) {
				return true
			}

			// Check second arg is a constant literal > 1
			sizeArg := call.Args[1]
			val := pass.TypesInfo.Types[sizeArg].Value
			if val == nil {
				return true // dynamic, not lintable
			}
			intVal, ok := constant.Int64Val(val)
			if !ok {
				return true
			}
			if intVal > 1 {
				pass.Reportf(sizeArg.Pos(), "chansize: channel size should be 0 or 1, got %d", intVal)
			}
			return true
		})
	}
	return nil, nil
}

func isChanType(expr ast.Expr) bool {
	_, ok := expr.(*ast.ChanType)
	return ok
}

package analyzers

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var ConstPrintf = &analysis.Analyzer{
	Name: "constprintf",
	Doc:  "check inline printf format strings (Uber Go Style Guide: Format Strings outside Printf)",
	URL:  "https://github.com/uber-go/guide#format-strings-outside-printf",
	Run:  runConstPrintf,
}

func runConstPrintf(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 || !isPrintfStyleCall(call) {
				return true
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			pass.Reportf(lit.Pos(), "move printf format string to a const")
			return true
		})
	}
	return nil, nil
}

func isPrintfStyleCall(call *ast.CallExpr) bool {
	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		return strings.HasSuffix(fn.Sel.Name, "f")
	case *ast.Ident:
		return strings.HasSuffix(fn.Name, "f")
	default:
		return false
	}
}

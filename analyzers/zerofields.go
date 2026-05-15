package analyzers

import (
	"go/ast"
	"go/constant"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var ZeroFields = &analysis.Analyzer{
	Name: "zerofields",
	Doc:  "check for explicit zero-value fields in struct literals (Uber Go Style Guide: Omit Zero Value Fields in Structs)",
	URL:  "https://github.com/uber-go/guide#omit-zero-value-fields-in-structs",
	Run:  runZeroFields,
}

func runZeroFields(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			lit, ok := node.(*ast.CompositeLit)
			if !ok || !isStructType(pass, lit) {
				return true
			}
			for _, elt := range lit.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				if isZeroValueExpr(pass, kv.Value) {
					pass.Reportf(kv.Value.Pos(), "omit zero-value field from struct literal")
				}
			}
			return true
		})
	}
	return nil, nil
}

func isZeroValueExpr(pass *analysis.Pass, expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "nil" || ident.Name == "false"
	}
	if lit, ok := expr.(*ast.BasicLit); ok {
		switch lit.Kind {
		case token.STRING:
			return lit.Value == `""`
		case token.INT, token.FLOAT, token.IMAG:
			tv := pass.TypesInfo.Types[lit]
			if tv.Value == nil {
				return false
			}
			return constant.Sign(tv.Value) == 0
		}
	}
	return false
}

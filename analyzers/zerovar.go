package analyzers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var ZeroVar = &analysis.Analyzer{
	Name: "zerovar",
	Doc:  "check for T{} short decl where var T should be used (Uber Go Style Guide: Use var for Zero Value Structs)",
	URL:  "https://github.com/uber-go/guide#use-var-for-zero-value-structs",
	Run:  runZeroVar,
}

func runZeroVar(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			assign, ok := node.(*ast.AssignStmt)
			if !ok || assign.Tok != token.DEFINE {
				return true
			}
			if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
				return true
			}
			lit, ok := assign.Rhs[0].(*ast.CompositeLit)
			if !ok {
				return true
			}
			if len(lit.Elts) != 0 {
				return true // has fields, not zero-value
			}
			// Only flag struct types — maps and slices have different nil vs empty semantics
			if !isStructType(pass, lit) {
				return true
			}
			pass.Reportf(lit.Pos(), "use var declaration instead of short assignment with zero-value struct")
			return true
		})
	}
	return nil, nil
}

func isStructType(pass *analysis.Pass, lit *ast.CompositeLit) bool {
	if lit.Type == nil {
		return false
	}
	t := pass.TypesInfo.TypeOf(lit.Type)
	if t == nil {
		return false
	}
	_, ok := t.Underlying().(*types.Struct)
	return ok
}

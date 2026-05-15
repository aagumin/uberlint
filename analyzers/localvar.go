package analyzers

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var LocalVar = &analysis.Analyzer{
	Name: "localvar",
	Doc:  "check for local var declarations that should use := (Uber Go Style Guide: Local Variable Declarations)",
	URL:  "https://github.com/uber-go/guide#local-variable-declarations",
	Run:  runLocalVar,
}

func runLocalVar(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				return true
			}
			ast.Inspect(fn.Body, func(bodyNode ast.Node) bool {
				decl, ok := bodyNode.(*ast.GenDecl)
				if !ok || decl.Tok != token.VAR {
					return true
				}
				for _, spec := range decl.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok || vs.Type != nil || len(vs.Values) == 0 {
						continue
					}
					pass.Reportf(vs.Pos(), "use short variable declaration instead of local var with initializer")
				}
				return true
			})
			return false
		})
	}
	return nil, nil
}

package analyzers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var VarType = &analysis.Analyzer{
	Name: "vartype",
	Doc:  "check for redundant type in top-level var declarations (Uber Go Style Guide: Top-level Variable Declarations)",
	URL:  "https://github.com/uber-go/guide#top-level-variable-declarations",
	Run:  runVarType,
}

func runVarType(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || vs.Type == nil {
					continue
				}
				if len(vs.Values) == 0 {
					continue // no initializer, type is needed
				}
				// Compare declared type with expression type
				declaredType := pass.TypesInfo.TypeOf(vs.Type)
				exprType := pass.TypesInfo.TypeOf(vs.Values[0])
				if declaredType == nil || exprType == nil {
					continue
				}
				if types.Identical(declaredType, exprType) {
					pass.Reportf(vs.Type.Pos(), "vartype: omit redundant type in var declaration, the expression already returns %s", declaredType)
				}
			}
		}
	}
	return nil, nil
}

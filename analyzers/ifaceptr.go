package analyzers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var IfacePtr = &analysis.Analyzer{
	Name: "ifaceptr",
	Doc:  "check for pointer-to-interface types (Uber Go Style Guide: Pointers to Interfaces)",
	URL:  "https://github.com/uber-go/guide#pointers-to-interfaces",
	Run:  runIfacePtr,
}

func runIfacePtr(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.FuncDecl:
				if n.Type.Params != nil {
					for _, field := range n.Type.Params.List {
						checkPointerType(pass, field.Type)
					}
				}
				if n.Type.Results != nil {
					for _, field := range n.Type.Results.List {
						checkPointerType(pass, field.Type)
					}
				}
			case *ast.StructType:
				if n.Fields != nil {
					for _, field := range n.Fields.List {
						checkPointerType(pass, field.Type)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func checkPointerType(pass *analysis.Pass, expr ast.Expr) {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return
	}
	t := pass.TypesInfo.TypeOf(star.X)
	if t == nil {
		return
	}
	if _, isIface := t.Underlying().(*types.Interface); isIface {
		pass.Reportf(star.Pos(), "avoid pointer to interface, interfaces already contain a pointer to data")
	}
}

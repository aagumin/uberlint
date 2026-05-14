package analyzers

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var EnumStart = &analysis.Analyzer{
	Name: "enumstart",
	Doc:  "check that enum constants start at a non-zero value (Uber Go Style Guide: Start Enums at One)",
	URL:  "https://github.com/uber-go/guide#start-enums-at-one",
	Run:  runEnumStart,
}

func runEnumStart(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}
			// Only check parenthesized groups (iota enums)
			if !genDecl.Lparen.IsValid() {
				continue
			}
			if !groupUsesIota(genDecl) {
				continue
			}
			// Check if first spec starts at zero without offset
			firstSpec := genDecl.Specs[0].(*ast.ValueSpec)
			if firstSpec.Values == nil || isPlainIota(firstSpec.Values) {
				if hasEnumStartSuppression(genDecl, firstSpec) {
					continue
				}
				pass.Reportf(firstSpec.Pos(), "enumstart: enum should start at a non-zero value")
			}
		}
	}
	return nil, nil
}

func groupUsesIota(genDecl *ast.GenDecl) bool {
	for _, spec := range genDecl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if vs.Values == nil {
			// Implicit continuation of iota
			return true
		}
		for _, val := range vs.Values {
			if containsIota(val) {
				return true
			}
		}
	}
	return false
}

func isPlainIota(values []ast.Expr) bool {
	if len(values) != 1 {
		return false
	}
	ident, ok := values[0].(*ast.Ident)
	return ok && ident.Name == "iota"
}

func containsIota(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if ok && ident.Name == "iota" {
		return true
	}
	bin, ok := expr.(*ast.BinaryExpr)
	if ok {
		return containsIota(bin.X) || containsIota(bin.Y)
	}
	return false
}

func hasEnumStartSuppression(decl *ast.GenDecl, spec *ast.ValueSpec) bool {
	return commentGroupSuppresses(decl.Doc, "enumstart") ||
		commentGroupSuppresses(spec.Doc, "enumstart") ||
		commentGroupSuppresses(spec.Comment, "enumstart")
}

func commentGroupSuppresses(group *ast.CommentGroup, ruleName string) bool {
	if group == nil {
		return false
	}
	for _, c := range group.List {
		if strings.Contains(c.Text, "//nolint:"+ruleName) ||
			strings.Contains(c.Text, "//nolint:all") ||
			strings.Contains(c.Text, "//uberskip:"+ruleName) {
			return true
		}
	}
	return false
}

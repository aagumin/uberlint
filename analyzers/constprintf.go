package analyzers

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var ConstPrintf = &analysis.Analyzer{
	Name: "constprintf",
	Doc:  "check non-const printf format string variables (Uber Go Style Guide: Format Strings outside Printf)",
	URL:  "https://github.com/uber-go/guide#format-strings-outside-printf",
	Run:  runConstPrintf,
}

func runConstPrintf(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		formatVars := collectPrintfFormatVars(file)
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 || !isPrintfStyleCall(call) {
				return true
			}
			ident, ok := call.Args[0].(*ast.Ident)
			if !ok {
				return true
			}
			if _, ok := formatVars[ident.Name]; !ok {
				return true
			}
			pass.Reportf(ident.Pos(), "make printf format string a const")
			return true
		})
	}
	return nil, nil
}

func collectPrintfFormatVars(file *ast.File) map[string]struct{} {
	vars := map[string]struct{}{}
	ast.Inspect(file, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.AssignStmt:
			if n.Tok != token.DEFINE {
				return true
			}
			for i, rhs := range n.Rhs {
				if i >= len(n.Lhs) || !isStringLiteral(rhs) {
					continue
				}
				ident, ok := n.Lhs[i].(*ast.Ident)
				if ok && ident.Name != "_" {
					vars[ident.Name] = struct{}{}
				}
			}
		case *ast.GenDecl:
			if n.Tok != token.VAR {
				return true
			}
			for _, spec := range n.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for i, value := range vs.Values {
					if i >= len(vs.Names) || !isStringLiteral(value) {
						continue
					}
					name := vs.Names[i]
					if name.Name != "_" {
						vars[name.Name] = struct{}{}
					}
				}
			}
		}
		return true
	})
	return vars
}

func isStringLiteral(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind == token.STRING && hasPrintfVerb(lit.Value)
}

func hasPrintfVerb(format string) bool {
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		if i+1 < len(format) && format[i+1] == '%' {
			i++
			continue
		}
		return true
	}
	return false
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

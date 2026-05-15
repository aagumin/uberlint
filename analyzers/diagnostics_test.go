package analyzers

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"strconv"
	"strings"
	"testing"
)

func TestDiagnosticsDoNotIncludeAnalyzerNamePrefix(t *testing.T) {
	analyzerNames := map[string]struct{}{}
	for _, analyzer := range allAnalyzersForUnitTests() {
		analyzerNames[analyzer.Name] = struct{}{}
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", func(info fs.FileInfo) bool {
		return strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parse analyzer package: %v", err)
	}

	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok || len(call.Args) < 2 {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok || sel.Sel.Name != "Reportf" {
					return true
				}
				lit, ok := call.Args[1].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return true
				}
				message, err := strconv.Unquote(lit.Value)
				if err != nil {
					t.Fatalf("unquote diagnostic at %s: %v", fset.Position(lit.Pos()), err)
				}
				for name := range analyzerNames {
					if strings.HasPrefix(message, name+":") {
						t.Fatalf("%s: diagnostic message should not include analyzer prefix %q: %q", filename, name+":", message)
					}
				}
				return true
			})
		}
	}
}

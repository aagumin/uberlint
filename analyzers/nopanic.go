package analyzers

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var NoPanic = &analysis.Analyzer{
	Name: "nopanic",
	Doc:  "check for panic() calls in production code (Uber Go Style Guide: Don't Panic)",
	URL:  "https://github.com/uber-go/guide#dont-panic",
	Run:  runNoPanic,
}

func runNoPanic(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.HasSuffix(filename, "_test.go") {
			continue
		}

		ast.Inspect(file, nodeVisitor(pass, checkPanicCall))
	}
	return nil, nil
}

func nodeVisitor(pass *analysis.Pass, check func(pass *analysis.Pass, stack []ast.Node)) func(ast.Node) bool {
	var stack []ast.Node
	return func(node ast.Node) bool {
		if node == nil {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			return true
		}
		stack = append(stack, node)
		check(pass, stack)
		return true
	}
}

func checkPanicCall(pass *analysis.Pass, stack []ast.Node) {
	node := stack[len(stack)-1]
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return
	}
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "panic" {
		return
	}

	fnDecl := enclosingFuncDecl(stack)
	if fnDecl != nil {
		name := fnDecl.Name.Name
		if name == "init" || strings.HasPrefix(name, "Must") {
			return
		}
	}

	pass.Reportf(call.Pos(), "avoid panic in production code, return an error instead")
}

func enclosingFuncDecl(stack []ast.Node) *ast.FuncDecl {
	for i := len(stack) - 2; i >= 0; i-- {
		if fd, ok := stack[i].(*ast.FuncDecl); ok {
			return fd
		}
	}
	return nil
}

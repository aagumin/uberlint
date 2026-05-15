package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var StringBytes = &analysis.Analyzer{
	Name: "stringbytes",
	Doc:  "check for repeated string-to-byte conversions inside loops (Uber Go Style Guide: Avoid repeated string-to-byte conversions)",
	URL:  "https://github.com/uber-go/guide#avoid-repeated-string-to-byte-conversions",
	Run:  runStringBytes,
}

func runStringBytes(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			loop, ok := node.(*ast.RangeStmt)
			if ok {
				checkConvInBlock(pass, loop.Body)
				return true
			}
			forStmt, ok := node.(*ast.ForStmt)
			if ok {
				checkConvInBlock(pass, forStmt.Body)
				return true
			}
			return true
		})
	}
	return nil, nil
}

func checkConvInBlock(pass *analysis.Pass, block *ast.BlockStmt) {
	if block == nil {
		return
	}
	ast.Inspect(block, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		// Check for []byte(x) type conversion — Fun is an *ast.ArrayType
		arr, ok := call.Fun.(*ast.ArrayType)
		if !ok || arr.Len != nil {
			return true
		}
		ident, ok := arr.Elt.(*ast.Ident)
		if !ok || ident.Name != "byte" {
			return true
		}
		if len(call.Args) != 1 {
			return true
		}
		// Only flag if argument is a string literal (not a variable)
		if _, isLit := call.Args[0].(*ast.BasicLit); isLit {
			pass.Reportf(call.Pos(), "avoid repeated string-to-byte conversion")
		}
		return true
	})
}

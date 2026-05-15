package analyzers

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var RawString = &analysis.Analyzer{
	Name: "rawstring",
	Doc:  "check for string literals with escaped quotes that could use backticks (Uber Go Style Guide: Use Raw String Literals to Avoid Escaping)",
	URL:  "https://github.com/uber-go/guide#use-raw-string-literals-to-avoid-escaping",
	Run:  runRawString,
}

func runRawString(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			lit, ok := node.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			value := lit.Value
			// Already a raw string (backtick)
			if strings.HasPrefix(value, "`") {
				return true
			}
			// Check if string contains escaped quotes
			if !strings.Contains(value, `\"`) {
				return true
			}

			unquoted, err := strconv.Unquote(value)
			if err != nil {
				return true
			}
			if strings.Contains(unquoted, "`") {
				return true
			}
			if hasEscapeOtherThanQuote(value[1 : len(value)-1]) {
				return true
			}
			// Skip very short strings
			if len(unquoted) < 6 {
				return true
			}
			pass.Reportf(lit.Pos(), "use raw string literal (backticks) to avoid escaping quotes")
			return true
		})
	}
	return nil, nil
}

func hasEscapeOtherThanQuote(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			continue
		}
		if i+1 >= len(s) || s[i+1] != '"' {
			return true
		}
		i++
	}
	return false
}

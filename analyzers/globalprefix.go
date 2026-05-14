package analyzers

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var GlobalPrefix = &analysis.Analyzer{
	Name: "globalprefix",
	Doc:  "check that unexported top-level vars and consts are prefixed with _ (Uber Go Style Guide: Prefix Unexported Globals with _)",
	URL:  "https://github.com/uber-go/guide#prefix-unexported-globals-with-_",
	Run:  runGlobalPrefix,
}

func runGlobalPrefix(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if genDecl.Tok != token.VAR && genDecl.Tok != token.CONST {
				continue
			}
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range vs.Names {
					if name.Name == "_" {
						continue
					}
					if !isUnexported(name.Name) {
						continue // exported, OK
					}
					if strings.HasPrefix(name.Name, "_") {
						continue // already prefixed
					}
					if strings.HasPrefix(name.Name, "err") && len(name.Name) > 3 {
						continue // err prefix for errors
					}
					pass.Reportf(name.Pos(), "globalprefix: prefix unexported global with _")
				}
			}
		}
	}
	return nil, nil
}

func isUnexported(name string) bool {
	if name == "" {
		return false
	}
	return !unicode.IsUpper(rune(name[0]))
}

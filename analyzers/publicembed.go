package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var PublicEmbed = &analysis.Analyzer{
	Name: "publicembed",
	Doc:  "check for embedded fields in public structs (Uber Go Style Guide: Avoid Embedding Types in Public Structs)",
	URL:  "https://github.com/uber-go/guide#avoid-embedding-types-in-public-structs",
	Run:  runPublicEmbed,
}

func runPublicEmbed(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || !typeSpec.Name.IsExported() {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok || structType.Fields == nil {
					continue
				}
				for _, field := range structType.Fields.List {
					if len(field.Names) == 0 {
						pass.Reportf(field.Pos(), "publicembed: avoid embedding types in public structs")
					}
				}
			}
		}
	}
	return nil, nil
}

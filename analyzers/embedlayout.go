package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var EmbedLayout = &analysis.Analyzer{
	Name: "embedlayout",
	Doc:  "check that embedded struct fields are grouped before named fields (Uber Go Style Guide: Embedding in Structs)",
	URL:  "https://github.com/uber-go/guide#embedding-in-structs",
	Run:  runEmbedLayout,
}

func runEmbedLayout(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			structType, ok := node.(*ast.StructType)
			if !ok || structType.Fields == nil {
				return true
			}
			seenNamed := false
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					if seenNamed {
						pass.Reportf(field.Pos(), "embedlayout: embedded fields should be grouped before named fields")
					}
					continue
				}
				seenNamed = true
			}
			return true
		})
	}
	return nil, nil
}

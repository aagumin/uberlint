package analyzers

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var TimeField = &analysis.Analyzer{
	Name: "timefield",
	Doc:  "check serialized numeric time fields for unit suffixes (Uber Go Style Guide: Time with External Systems)",
	URL:  "https://github.com/uber-go/guide#use-time-with-external-systems",
	Run:  runTimeField,
}

func runTimeField(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			structType, ok := node.(*ast.StructType)
			if !ok || structType.Fields == nil {
				return true
			}
			for _, field := range structType.Fields.List {
				if field.Tag == nil || !hasSerializationTag(field.Tag.Value) || isStdTimeType(pass, field.Type) || !isNumericType(pass, field.Type) {
					continue
				}
				for _, name := range field.Names {
					if isTimeLikeName(name.Name) && !hasTimeUnitSuffix(name.Name) {
						pass.Reportf(name.Pos(), "include time unit in serialized numeric field name")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func hasSerializationTag(tag string) bool {
	return strings.Contains(tag, `json:"`) || strings.Contains(tag, `yaml:"`)
}

func isStdTimeType(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}
	if obj.Pkg().Path() != "time" {
		return false
	}
	return obj.Name() == "Duration" || obj.Name() == "Time"
}

func isNumericType(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	info := basic.Info()
	return info&types.IsInteger != 0 || info&types.IsFloat != 0
}

func isTimeLikeName(name string) bool {
	lower := strings.ToLower(name)
	for _, part := range []string{"timeout", "duration", "delay", "ttl", "expires", "interval"} {
		if strings.Contains(lower, part) {
			return true
		}
	}
	return false
}

func hasTimeUnitSuffix(name string) bool {
	lower := strings.ToLower(name)
	for _, suffix := range []string{"ns", "nanos", "nanoseconds", "us", "micros", "microseconds", "ms", "millis", "milliseconds", "sec", "secs", "second", "seconds", "min", "mins", "minute", "minutes", "hour", "hours"} {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}

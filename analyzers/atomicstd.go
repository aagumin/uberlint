package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var AtomicStd = &analysis.Analyzer{
	Name: "atomicstd",
	Doc:  "check for raw sync/atomic function calls instead of typed values (Uber Go Style Guide: Use sync/atomic typed values)",
	URL:  "https://github.com/uber-go/guide#use-gouberorgatomic",
	Run:  runAtomicStd,
}

// rawAtomicFuncs maps raw atomic function names to their typed equivalents.
var _rawAtomicFuncs = map[string]string{
	"AddInt32":             "atomic.Int32",
	"AddInt64":             "atomic.Int64",
	"AddUint32":            "atomic.Uint32",
	"AddUint64":            "atomic.Uint64",
	"AddUintptr":           "atomic.Uintptr",
	"LoadInt32":            "atomic.Int32",
	"LoadInt64":            "atomic.Int64",
	"LoadUint32":           "atomic.Uint32",
	"LoadUint64":           "atomic.Uint64",
	"LoadUintptr":          "atomic.Uintptr",
	"StoreInt32":           "atomic.Int32",
	"StoreInt64":           "atomic.Int64",
	"StoreUint32":          "atomic.Uint32",
	"StoreUint64":          "atomic.Uint64",
	"SwapInt32":            "atomic.Int32",
	"SwapInt64":            "atomic.Int64",
	"SwapUint32":           "atomic.Uint32",
	"SwapUint64":           "atomic.Uint64",
	"CompareAndSwapInt32":  "atomic.Int32",
	"CompareAndSwapInt64":  "atomic.Int64",
	"CompareAndSwapUint32": "atomic.Uint32",
	"CompareAndSwapUint64": "atomic.Uint64",
}

func runAtomicStd(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			typedName, exists := _rawAtomicFuncs[sel.Sel.Name]
			if !exists {
				return true
			}
			// Verify the selector resolves to sync/atomic package
			if !isAtomicPackage(pass, sel) {
				return true
			}
			pass.Reportf(call.Pos(), "use %s type instead of raw atomic.%s", typedName, sel.Sel.Name)
			return true
		})
	}
	return nil, nil
}

func isAtomicPackage(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	obj := pass.TypesInfo.ObjectOf(sel.Sel)
	if obj == nil {
		return false
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}
	return pkg.Path() == "sync/atomic"
}

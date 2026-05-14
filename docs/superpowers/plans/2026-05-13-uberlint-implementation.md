# uberlint Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement 12 Uber Go Style Guide analyzers as a golangci-lint module plugin using go/analysis.

**Architecture:** Each rule is a standalone `go/analysis.Analyzer` in the `analyzers/` package. Tests use `analysistest` with test data files. A standalone binary and golangci-lint module plugin entry point expose all analyzers.

**Tech Stack:** Go 1.24, `golang.org/x/tools/go/analysis`, `analysistest`, golangci-lint module plugin system

---

## File Structure

```
uberlint/
  go.mod
  go.sum
  cmd/uberlint/main.go              # standalone runner
  analyzers/
    nopanic.go                       # 12 analyzer implementations
    nopanic_test.go                  # 12 test files
    chansize.go
    chansize_test.go
    enumstart.go
    enumstart_test.go
    newref.go
    newref_test.go
    nilslice.go
    nilslice_test.go
    zerovar.go
    zerovar_test.go
    stringbytes.go
    stringbytes_test.go
    globalprefix.go
    globalprefix_test.go
    vartype.go
    vartype_test.go
    ifaceptr.go
    ifaceptr_test.go
    atomicstd.go
    atomicstd_test.go
    rawstring.go
    rawstring_test.go
    testdata/
      src/
        nopanic/nopanic.go
        chansize/chansize.go
        enumstart/enumstart.go
        newref/newref.go
        nilslice/nilslice.go
        zerovar/zerovar.go
        stringbytes/stringbytes.go
        globalprefix/globalprefix.go
        vartype/vartype.go
        ifaceptr/ifaceptr.go
        atomicstd/atomicstd.go
        rawstring/rawstring.go
  .golangci.yml                     # example integration config
  .custom-gcl.yml                   # custom binary build config
```

---

### Task 1: Project Scaffold

**Files:**
- Create: `go.mod`
- Create: `cmd/uberlint/main.go`

- [ ] **Step 1: Initialize Go module**

Run:
```bash
cd /Users/arsengumin/projects/uberlint
go mod init github.com/aagumin/uberlint
```

- [ ] **Step 2: Add dependencies**

Run:
```bash
go get golang.org/x/tools
```

- [ ] **Step 3: Create analyzers and testdata directories**

Run:
```bash
mkdir -p analyzers/testdata/src/{nopanic,chansize,enumstart,newref,nilslice,zerovar,stringbytes,globalprefix,vartype,ifaceptr,atomicstd,rawstring}
```

- [ ] **Step 4: Create standalone runner**

Create `cmd/uberlint/main.go`:

```go
package main

import (
	"github.com/aagumin/uberlint/analyzers"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		analyzers.NoPanic,
		analyzers.ChanSize,
		analyzers.EnumStart,
		analyzers.NewRef,
		analyzers.NilSlice,
		analyzers.ZeroVar,
		analyzers.StringBytes,
		analyzers.GlobalPrefix,
		analyzers.VarType,
		analyzers.IfacePtr,
		analyzers.AtomicStd,
		analyzers.RawString,
	)
}
```

- [ ] **Step 5: Verify build compiles**

Run:
```bash
go build ./...
```

Expected: compile errors about missing analyzers — this is expected, they'll be created in following tasks.

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum cmd/ analyzers/testdata/
git commit -m "feat: scaffold project structure with go.mod and cmd entry point"
```

---

### Task 2: nopanic Analyzer

**Files:**
- Create: `analyzers/testdata/src/nopanic/nopanic.go`
- Create: `analyzers/nopanic_test.go`
- Create: `analyzers/nopanic.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/nopanic/nopanic.go`:

```go
package nopanic

import "text/template"

func badPanic() {
	panic("something went wrong") // want "avoid panic in production code"
}

func goodError() error {
	return nil
}

func init() {
	panic("setup failed") // OK: init is allowed
}

func MustParse(s string) *template.Template {
	panic("invalid template") // OK: Must* pattern
}

type Foo struct{}

func (f *Foo) bad() {
	panic("method panic") // want "avoid panic in production code"
}
```

- [ ] **Step 2: Write test**

Create `analyzers/nopanic_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNoPanic(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NoPanic, "nopanic")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestNoPanic -v`
Expected: FAIL — `NoPanic` not defined

- [ ] **Step 4: Implement analyzer**

Create `analyzers/nopanic.go`:

```go
package analyzers

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var NoPanic = &analysis.Analyzer{
	Name:     "nopanic",
	Doc:      "check for panic() calls in production code (Uber Go Style Guide: Don't Panic)",
	URL:      "https://github.com/uber-go/guide#dont-panic",
	Run:      runNoPanic,
	Requires: []*analysis.Analyzer{},
}

func runNoPanic(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.HasSuffix(filename, "_test.go") {
			continue
		}

		ast.Inspect(file, nodeVisitor(pass, "nopanic", checkPanicCall))
	}
	return nil, nil
}

func nodeVisitor(pass *analysis.Pass, name string, check func(pass *analysis.Pass, stack []ast.Node)) func(ast.Node) bool {
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

	pass.Reportf(call.Pos(), "nopanic: avoid panic in production code, return an error instead")
}

func enclosingFuncDecl(stack []ast.Node) *ast.FuncDecl {
	for i := len(stack) - 2; i >= 0; i-- {
		if fd, ok := stack[i].(*ast.FuncDecl); ok {
			return fd
		}
	}
	return nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestNoPanic -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/nopanic.go analyzers/nopanic_test.go analyzers/testdata/src/nopanic/
git commit -m "feat: add nopanic analyzer — flags panic() in production code"
```

---

### Task 3: chansize Analyzer

**Files:**
- Create: `analyzers/testdata/src/chansize/chansize.go`
- Create: `analyzers/chansize_test.go`
- Create: `analyzers/chansize.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/chansize/chansize.go`:

```go
package chansize

var n = 10

func bad() {
	_ = make(chan int, 64)  // want "channel size should be 0 or 1"
	_ = make(chan int, 100) // want "channel size should be 0 or 1"
	_ = make(chan int, 2)   // want "channel size should be 0 or 1"
}

func good() {
	_ = make(chan int)     // OK: unbuffered
	_ = make(chan int, 0)  // OK: unbuffered explicit
	_ = make(chan int, 1)  // OK: size 1
}

func dynamic() {
	_ = make(chan int, n) // OK: dynamic size, not lintable
}
```

- [ ] **Step 2: Write test**

Create `analyzers/chansize_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestChanSize(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ChanSize, "chansize")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestChanSize -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/chansize.go`:

```go
package analyzers

import (
	"go/ast"
	"go/constant"

	"golang.org/x/tools/go/analysis"
)

var ChanSize = &analysis.Analyzer{
	Name: "chansize",
	Doc:  "check that channel size is 0 or 1 (Uber Go Style Guide: Channel Size is One or None)",
	URL:  "https://github.com/uber-go/guide#channel-size-is-one-or-none",
	Run:  runChanSize,
}

func runChanSize(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != "make" {
				return true
			}
			if len(call.Args) < 2 {
				return true
			}

			// Check first arg is a chan type
			if !isChanType(call.Args[0]) {
				return true
			}

			// Check second arg is a constant literal > 1
			sizeArg := call.Args[1]
			val := pass.TypesInfo.Types[sizeArg].Value
			if val == nil {
				return true // dynamic, not lintable
			}
			intVal, ok := constant.Int64Val(val)
			if !ok {
				return true
			}
			if intVal > 1 {
				pass.Reportf(sizeArg.Pos(), "chansize: channel size should be 0 or 1, got %d", intVal)
			}
			return true
		})
	}
	return nil, nil
}

func isChanType(expr ast.Expr) bool {
	_, ok := expr.(*ast.ChanType)
	return ok
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestChanSize -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/chansize.go analyzers/chansize_test.go analyzers/testdata/src/chansize/
git commit -m "feat: add chansize analyzer — flags channels with buffer > 1"
```

---

### Task 4: enumstart Analyzer

**Files:**
- Create: `analyzers/testdata/src/enumstart/enumstart.go`
- Create: `analyzers/enumstart_test.go`
- Create: `analyzers/enumstart.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/enumstart/enumstart.go`:

```go
package enumstart

type Operation int

const (
	Add Operation = iota // want "enum should start at a non-zero value"
	Subtract
	Multiply
)

type Status int

const (
	StatusOK  Status = iota + 1 // OK: starts at 1
	StatusError
)

type LogOutput int

const (
	LogToStdout LogOutput = iota // OK: zero value is intentional default
	LogToFile
	LogToRemote
)

const single = 42 // OK: not a group
```

- [ ] **Step 2: Write test**

Create `analyzers/enumstart_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEnumStart(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), EnumStart, "enumstart")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestEnumStart -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/enumstart.go`:

```go
package analyzers

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var EnumStart = &analysis.Analyzer{
	Name: "enumstart",
	Doc:  "check that enum constants start at a non-zero value (Uber Go Style Guide: Start Enums at One)",
	URL:  "https://github.com/uber-go/guide#start-enums-at-one",
	Run:  runEnumStart,
}

func runEnumStart(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}
			// Only check parenthesized groups (iota enums)
			if !genDecl.Lparen.IsValid() {
				continue
			}
			if !groupUsesIota(genDecl) {
				continue
			}
			// Check if first spec starts at zero without offset
			firstSpec := genDecl.Specs[0].(*ast.ValueSpec)
			if firstSpec.Values == nil {
				// Plain iota — starts at 0
				pass.Reportf(firstSpec.Pos(), "enumstart: enum should start at a non-zero value, use `iota + 1` or add a zero-value comment")
			}
		}
	}
	return nil, nil
}

func groupUsesIota(genDecl *ast.GenDecl) bool {
	for _, spec := range genDecl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if vs.Values == nil {
			// Implicit continuation of iota
			return true
		}
		for _, val := range vs.Values {
			if containsIota(val) {
				return true
			}
		}
	}
	return false
}

func containsIota(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if ok && ident.Name == "iota" {
		return true
	}
	bin, ok := expr.(*ast.BinaryExpr)
	if ok {
		return containsIota(bin.X) || containsIota(bin.Y)
	}
	return false
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestEnumStart -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/enumstart.go analyzers/enumstart_test.go analyzers/testdata/src/enumstart/
git commit -m "feat: add enumstart analyzer — flags iota enums starting at 0"
```

---

### Task 5: newref Analyzer

**Files:**
- Create: `analyzers/testdata/src/newref/newref.go`
- Create: `analyzers/newref_test.go`
- Create: `analyzers/newref.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/newref/newref.go`:

```go
package newref

type MyStruct struct {
	Name string
}

func bad() {
	_ = new(MyStruct) // want "use &T{} instead of new(T)"
	_ = new(int)      // want "use &T{} instead of new(T)"
}

func good() {
	_ = &MyStruct{Name: "foo"}
	_ = new(int) // want "use &T{} instead of new(T)"
	s := MyStruct{}
	_ = &s
}
```

- [ ] **Step 2: Write test**

Create `analyzers/newref_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNewRef(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NewRef, "newref")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestNewRef -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/newref.go`:

```go
package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var NewRef = &analysis.Analyzer{
	Name: "newref",
	Doc:  "check for new(T) instead of &T{} (Uber Go Style Guide: Initializing Struct References)",
	URL:  "https://github.com/uber-go/guide#initializing-struct-references",
	Run:  runNewRef,
}

func runNewRef(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != "new" {
				return true
			}
			pass.Reportf(call.Pos(), "newref: use &T{} instead of new(T) for consistency with struct initialization")
			return true
		})
	}
	return nil, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestNewRef -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/newref.go analyzers/newref_test.go analyzers/testdata/src/newref/
git commit -m "feat: add newref analyzer — flags new(T) instead of &T{}"
```

---

### Task 6: nilslice Analyzer

**Files:**
- Create: `analyzers/testdata/src/nilslice/nilslice.go`
- Create: `analyzers/nilslice_test.go`
- Create: `analyzers/nilslice.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/nilslice/nilslice.go`:

```go
package nilslice

func badEmptyLiteral() []int {
	return []int{} // want "return nil instead of an empty slice"
}

func badMakeEmpty() []string {
	return make([]string, 0) // want "return nil instead of an empty slice"
}

func good() []int {
	return nil
}

func goodNonEmpty() []int {
	return []int{1, 2, 3}
}

func badAssign() {
	_ = []int{} // want "use var or nil instead of empty slice literal"
}
```

- [ ] **Step 2: Write test**

Create `analyzers/nilslice_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNilSlice(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NilSlice, "nilslice")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestNilSlice -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/nilslice.go`:

```go
package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var NilSlice = &analysis.Analyzer{
	Name: "nilslice",
	Doc:  "check for empty slice literals where nil should be used (Uber Go Style Guide: nil is a valid slice)",
	URL:  "https://github.com/uber-go/guide#nil-is-a-valid-slice",
	Run:  runNilSlice,
}

func runNilSlice(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			// Check for []T{} composite literals with no elements
			lit, ok := node.(*ast.CompositeLit)
			if !ok {
				return true
			}
			if !isEmptySliceLiteral(lit) {
				return true
			}
			pass.Reportf(lit.Pos(), "nilslice: use nil or var declaration instead of empty slice literal")
			return true
		})
	}
	return nil, nil
}

func isEmptySliceLiteral(lit *ast.CompositeLit) bool {
	if len(lit.Elts) != 0 {
		return false
	}
	// Check the type is a slice: []T
	at, ok := lit.Type.(*ast.ArrayType)
	if !ok {
		return false
	}
	// Slice has no length (nil length = slice, not array)
	return at.Len == nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestNilSlice -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/nilslice.go analyzers/nilslice_test.go analyzers/testdata/src/nilslice/
git commit -m "feat: add nilslice analyzer — flags empty slice literals instead of nil"
```

---

### Task 7: zerovar Analyzer

**Files:**
- Create: `analyzers/testdata/src/zerovar/zerovar.go`
- Create: `analyzers/zerovar_test.go`
- Create: `analyzers/zerovar.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/zerovar/zerovar.go`:

```go
package zerovar

type User struct {
	Name string
	Age  int
}

func bad() {
	u := User{} // want "use var u User instead of u := User{}"
	_ = u
}

func good() {
	var u User
	_ = u

	initialized := User{Name: "Alice"} // OK: has non-zero fields
	_ = initialized
}
```

- [ ] **Step 2: Write test**

Create `analyzers/zerovar_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestZeroVar(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ZeroVar, "zerovar")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestZeroVar -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/zerovar.go`:

```go
package analyzers

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var ZeroVar = &analysis.Analyzer{
	Name: "zerovar",
	Doc:  "check for T{} short decl where var T should be used (Uber Go Style Guide: Use var for Zero Value Structs)",
	URL:  "https://github.com/uber-go/guide#use-var-for-zero-value-structs",
	Run:  runZeroVar,
}

func runZeroVar(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			assign, ok := node.(*ast.AssignStmt)
			if !ok || assign.Tok != token.DEFINE {
				return true
			}
			if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
				return true
			}
			lit, ok := assign.Rhs[0].(*ast.CompositeLit)
			if !ok {
				return true
			}
			if len(lit.Elts) != 0 {
				return true // has fields, not zero-value
			}
			pass.Reportf(lit.Pos(), "zerovar: use var declaration instead of short assignment with zero-value struct")
			return true
		})
	}
	return nil, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestZeroVar -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/zerovar.go analyzers/zerovar_test.go analyzers/testdata/src/zerovar/
git commit -m "feat: add zerovar analyzer — flags T{} short decl for zero-value structs"
```

---

### Task 8: stringbytes Analyzer

**Files:**
- Create: `analyzers/testdata/src/stringbytes/stringbytes.go`
- Create: `analyzers/stringbytes_test.go`
- Create: `analyzers/stringbytes.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/stringbytes/stringbytes.go`:

```go
package stringbytes

import "io"

func bad(w io.Writer) {
	for i := 0; i < 10; i++ {
		w.Write([]byte("Hello world")) // want "avoid repeated string-to-byte conversion"
	}
}

func good(w io.Writer) {
	data := []byte("Hello world")
	for i := 0; i < 10; i++ {
		w.Write(data)
	}
}

func outsideLoop(w io.Writer) {
	w.Write([]byte("one time")) // OK: not in a loop
}
```

- [ ] **Step 2: Write test**

Create `analyzers/stringbytes_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestStringBytes(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), StringBytes, "stringbytes")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestStringBytes -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/stringbytes.go`:

```go
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
		conv, ok := call.Fun.(*ast.Ident)
		if !ok || conv.Name != "[]byte" {
			return true
		}
		if len(call.Args) != 1 {
			return true
		}
		// Only flag if argument is a string literal (not a variable)
		if _, isLit := call.Args[0].(*ast.BasicLit); isLit {
			pass.Reportf(call.Pos(), "stringbytes: avoid repeated string-to-byte conversion, hoist the conversion outside the loop")
		}
		return true
	})
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestStringBytes -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/stringbytes.go analyzers/stringbytes_test.go analyzers/testdata/src/stringbytes/
git commit -m "feat: add stringbytes analyzer — flags []byte(string) inside loops"
```

---

### Task 9: globalprefix Analyzer

**Files:**
- Create: `analyzers/testdata/src/globalprefix/globalprefix.go`
- Create: `analyzers/globalprefix_test.go`
- Create: `analyzers/globalprefix.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/globalprefix/globalprefix.go`:

```go
package globalprefix

import "errors"

const (
	defaultPort = 8080    // want "prefix unexported global with _"
	defaultUser = "user"  // want "prefix unexported global with _"
)

const (
	_defaultPort = 8080 // OK: prefixed
	_defaultUser = "user"
)

var ErrNotFound = errors.New("not found") // OK: Err prefix for errors
var errInternal = errors.New("internal")  // OK: err prefix for errors

var someGlobal = 42 // want "prefix unexported global with _"

var _okGlobal = 42 // OK: prefixed

const ExportedConst = 100 // OK: exported
```

- [ ] **Step 2: Write test**

Create `analyzers/globalprefix_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestGlobalPrefix(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), GlobalPrefix, "globalprefix")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestGlobalPrefix -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/globalprefix.go`:

```go
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
					if strings.HasPrefix(name.Name, "err") && len(name.Name) > 3 && !unicode.IsUpper(rune(name.Name[3])) {
						continue // err prefix for errors
					}
					pass.Reportf(name.Pos(), "globalprefix: prefix unexported global %q with _", name.Name)
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
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestGlobalPrefix -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/globalprefix.go analyzers/globalprefix_test.go analyzers/testdata/src/globalprefix/
git commit -m "feat: add globalprefix analyzer — flags unexported globals without _ prefix"
```

---

### Task 10: vartype Analyzer

**Files:**
- Create: `analyzers/testdata/src/vartype/vartype.go`
- Create: `analyzers/vartype_test.go`
- Create: `analyzers/vartype.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/vartype/vartype.go`:

```go
package vartype

func F() string { return "A" }

var _s string = F() // want "omit redundant type in var declaration"

// OK: type differs from expression type
var _e error = myError{}

type myError struct{}

func (myError) Error() string { return "error" }
```

- [ ] **Step 2: Write test**

Create `analyzers/vartype_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestVarType(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), VarType, "vartype")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestVarType -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/vartype.go`:

```go
package analyzers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var VarType = &analysis.Analyzer{
	Name: "vartype",
	Doc:  "check for redundant type in top-level var declarations (Uber Go Style Guide: Top-level Variable Declarations)",
	URL:  "https://github.com/uber-go/guide#top-level-variable-declarations",
	Run:  runVarType,
}

func runVarType(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || vs.Type == nil {
					continue
				}
				if len(vs.Values) == 0 {
					continue // no initializer, type is needed
				}
				// Compare declared type with expression type
				declaredType := pass.TypesInfo.TypeOf(vs.Type)
				exprType := pass.TypesInfo.TypeOf(vs.Values[0])
				if declaredType == nil || exprType == nil {
					continue
				}
				if types.Identical(declaredType, exprType) {
					pass.Reportf(vs.Type.Pos(), "vartype: omit redundant type in var declaration, the expression already returns %s", declaredType)
				}
			}
		}
	}
	return nil, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestVarType -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/vartype.go analyzers/vartype_test.go analyzers/testdata/src/vartype/
git commit -m "feat: add vartype analyzer — flags redundant type in var declarations"
```

---

### Task 11: ifaceptr Analyzer

**Files:**
- Create: `analyzers/testdata/src/ifaceptr/ifaceptr.go`
- Create: `analyzers/ifaceptr_test.go`
- Create: `analyzers/ifaceptr.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/ifaceptr/ifaceptr.go`:

```go
package ifaceptr

import "io"

type Reader interface {
	Read()
}

func bad(r *Reader) {} // want "avoid pointer to interface"

type BadStruct struct {
	R *Reader // want "avoid pointer to interface"
}

func good(r Reader) {}

type GoodStruct struct {
	R Reader
}

func badReturn() *io.Reader { // want "avoid pointer to interface"
	return nil
}
```

- [ ] **Step 2: Write test**

Create `analyzers/ifaceptr_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestIfacePtr(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), IfacePtr, "ifaceptr")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestIfacePtr -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/ifaceptr.go`:

```go
package analyzers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var IfacePtr = &analysis.Analyzer{
	Name: "ifaceptr",
	Doc:  "check for pointer-to-interface types (Uber Go Style Guide: Pointers to Interfaces)",
	URL:  "https://github.com/uber-go/guide#pointers-to-interfaces",
	Run:  runIfacePtr,
}

func runIfacePtr(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.FuncDecl:
				if n.Type.Params != nil {
					for _, field := range n.Type.Params.List {
						checkPointerType(pass, field.Type)
					}
				}
				if n.Type.Results != nil {
					for _, field := range n.Type.Results.List {
						checkPointerType(pass, field.Type)
					}
				}
			case *ast.StructType:
				if n.Fields != nil {
					for _, field := range n.Fields.List {
						checkPointerType(pass, field.Type)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func checkPointerType(pass *analysis.Pass, expr ast.Expr) {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return
	}
	t := pass.TypesInfo.TypeOf(star.X)
	if t == nil {
		return
	}
	if _, isIface := t.Underlying().(*types.Interface); isIface {
		pass.Reportf(star.Pos(), "ifaceptr: avoid pointer to interface, interfaces already contain a pointer to data")
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestIfacePtr -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/ifaceptr.go analyzers/ifaceptr_test.go analyzers/testdata/src/ifaceptr/
git commit -m "feat: add ifaceptr analyzer — flags pointer to interface types"
```

---

### Task 12: atomicstd Analyzer

**Files:**
- Create: `analyzers/testdata/src/atomicstd/atomicstd.go`
- Create: `analyzers/atomicstd_test.go`
- Create: `analyzers/atomicstd.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/atomicstd/atomicstd.go`:

```go
package atomicstd

import "sync/atomic"

var counter int64

func bad() {
	atomic.AddInt64(&counter, 1)    // want "use atomic.Int64 type"
	atomic.LoadInt64(&counter)      // want "use atomic.Int64 type"
	atomic.StoreInt64(&counter, 10) // want "use atomic.Int64 type"
}

func good() {
	var ac atomic.Int64
	ac.Add(1)
	ac.Load()
	ac.Store(10)
}
```

- [ ] **Step 2: Write test**

Create `analyzers/atomicstd_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAtomicStd(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), AtomicStd, "atomicstd")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestAtomicStd -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/atomicstd.go`:

```go
package analyzers

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var AtomicStd = &analysis.Analyzer{
	Name: "atomicstd",
	Doc:  "check for raw sync/atomic function calls instead of typed values (Uber Go Style Guide: Use sync/atomic typed values)",
	URL:  "https://github.com/uber-go/guide#use-gouberorgatomic",
	Run:  runAtomicStd,
}

// Raw atomic functions that have typed equivalents.
var rawAtomicFuncs = map[string]string{
	"AddInt32":    "atomic.Int32",
	"AddInt64":    "atomic.Int64",
	"AddUint32":   "atomic.Uint32",
	"AddUint64":   "atomic.Uint64",
	"AddUintptr":  "atomic.Uintptr",
	"LoadInt32":   "atomic.Int32",
	"LoadInt64":   "atomic.Int64",
	"LoadUint32":  "atomic.Uint32",
	"LoadUint64":  "atomic.Uint64",
	"LoadUintptr": "atomic.Uintptr",
	"StoreInt32":  "atomic.Int32",
	"StoreInt64":  "atomic.Int64",
	"StoreUint32": "atomic.Uint32",
	"StoreUint64": "atomic.Uint64",
	"SwapInt32":   "atomic.Int32",
	"SwapInt64":   "atomic.Int64",
	"SwapUint32":  "atomic.Uint32",
	"SwapUint64":  "atomic.Uint64",
	"CompareAndSwapInt32":    "atomic.Int32",
	"CompareAndSwapInt64":    "atomic.Int64",
	"CompareAndSwapUint32":   "atomic.Uint32",
	"CompareAndSwapUint64":   "atomic.Uint64",
}

func runAtomicStd(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Check if file imports sync/atomic
		importsAtomic := false
		for _, imp := range file.Imports {
			if imp.Path != nil && strings.Trim(imp.Path.Value, `"`) == "sync/atomic" {
				importsAtomic = true
				break
			}
		}
		if !importsAtomic {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			typedName, exists := rawAtomicFuncs[sel.Sel.Name]
			if !exists {
				return true
			}
			pass.Reportf(call.Pos(), "atomicstd: use %s type instead of raw atomic.%s", typedName, sel.Sel.Name)
			return true
		})
	}
	return nil, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestAtomicStd -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/atomicstd.go analyzers/atomicstd_test.go analyzers/testdata/src/atomicstd/
git commit -m "feat: add atomicstd analyzer — flags raw sync/atomic calls"
```

---

### Task 13: rawstring Analyzer

**Files:**
- Create: `analyzers/testdata/src/rawstring/rawstring.go`
- Create: `analyzers/rawstring_test.go`
- Create: `analyzers/rawstring.go`

- [ ] **Step 1: Create test data**

Create `analyzers/testdata/src/rawstring/rawstring.go`:

```go
package rawstring

func bad() {
	_ = "unknown name:\"test\""   // want "use raw string literal"
	_ = "path:\"/foo\" bar:\"baz\"" // want "use raw string literal"
}

func good() {
	_ = `unknown name:"test"`
	_ = "simple string"
	_ = "has `backtick` inside" // OK: can't use raw literal
	_ = "ab"                    // OK: too short, not worth it
}
```

- [ ] **Step 2: Write test**

Create `analyzers/rawstring_test.go`:

```go
package analyzers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRawString(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), RawString, "rawstring")
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./analyzers/ -run TestRawString -v`
Expected: FAIL

- [ ] **Step 4: Implement analyzer**

Create `analyzers/rawstring.go`:

```go
package analyzers

import (
	"go/ast"
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
			if !ok || lit.Kind != ast.STRING {
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
			// Check the unquoted content doesn't contain backticks
			unquoted := value[1 : len(value)-1] // strip surrounding quotes
			if strings.Contains(unquoted, "`") {
				return true // can't use raw literal
			}
			// Skip very short strings
			if len(unquoted) < 6 {
				return true
			}
			pass.Reportf(lit.Pos(), "rawstring: use raw string literal (backticks) to avoid escaping quotes")
			return true
		})
	}
	return nil, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./analyzers/ -run TestRawString -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add analyzers/rawstring.go analyzers/rawstring_test.go analyzers/testdata/src/rawstring/
git commit -m "feat: add rawstring analyzer — flags strings with escaped quotes"
```

---

### Task 14: golangci-lint Integration

**Files:**
- Create: `.golangci.yml`
- Create: `.custom-gcl.yml`
- Modify: `cmd/uberlint/main.go` (verify it builds)

- [ ] **Step 1: Verify all tests pass**

Run: `go test ./analyzers/ -v`
Expected: all 12 tests PASS

- [ ] **Step 2: Verify standalone binary builds and works**

Run:
```bash
go build -o uberlint ./cmd/uberlint/
echo 'package main; func main() { panic("x") }' > /tmp/test.go
./uberlint /tmp/test.go
```

Expected: output contains `nopanic` diagnostic

- [ ] **Step 3: Create .custom-gcl.yml for golangci-lint module plugin**

Create `.custom-gcl.yml`:

```yaml
version: v2.0.0
name: custom-gcl
plugins:
  - module: 'github.com/aagumin/uberlint'
    path: .
```

- [ ] **Step 4: Create example .golangci.yml**

Create `.golangci.yml`:

```yaml
version: "2"

linters:
  default: none
  enable:
    - uberlint
  settings:
    custom:
      uberlint:
        type: "module"
        description: "Uber Go Style Guide linter — checks for 12 rules not covered by existing linters"
```

- [ ] **Step 5: Commit**

```bash
git add .golangci.yml .custom-gcl.yml
git commit -m "feat: add golangci-lint integration config"
```

---

### Task 15: Final Verification

- [ ] **Step 1: Run all tests**

Run: `go test ./... -v`
Expected: all tests PASS

- [ ] **Step 2: Run go vet**

Run: `go vet ./...`
Expected: no issues

- [ ] **Step 3: Run the linter on its own code (dogfooding)**

Run: `go build -o uberlint ./cmd/uberlint/ && ./uberlint ./analyzers/`
Expected: no diagnostics from our own analyzers (our code should be clean)

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "chore: final verification — all analyzers pass tests and vet"
```

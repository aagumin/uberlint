# uberlint — Uber Go Style Guide Linter

## Problem

58% of Uber Go Style Guide rules have no automatic linter. Existing tools (golangci-lint with its 50+ linters) cover only ~19% of the guide fully. This leaves developers without automated enforcement for rules like "Don't Panic", "Channel Size is One or None", "Start Enums at One", etc.

## Solution

A golangci-lint plugin (`uberlint`) that implements missing Uber Go Style Guide checks. Each rule is a separate `go/analysis.Analyzer` following the standard Go linter framework.

## Decisions

- **Format:** golangci-lint custom plugin (requires golangci-lint >= 1.53)
- **Framework:** `golang.org/x/tools/go/analysis` — standard framework used by govet, staticcheck
- **Go version:** 1.24+
- **Scope:** MVP of 12 rules, expandable to all 35 unlinted rules

## MVP Rules (12 analyzers)

### Easy (AST pattern matching)

| ID | Rule | What it flags |
|---|---|---|
| `nopanic` | Don't Panic | `panic()` calls outside test files and `func init()` |
| `chansize` | Channel Size is One or None | `make(chan T, N)` where N is a literal > 1 |
| `enumstart` | Start Enums at One | `const` group using `iota` without `+ 1` (unless zero value is intentional) |
| `newref` | Initializing Struct References | `new(T)` instead of `&T{}` |
| `nilslice` | nil is a valid slice | `return []T{}` or `return make([]T, 0)` when `return nil` suffices |
| `zerovar` | Use var for Zero Value Structs | `x := T{}` when all fields are zero, should be `var x T` |

### Medium (requires type info or multi-node analysis)

| ID | Rule | What it flags |
|---|---|---|
| `stringbytes` | Avoid repeated string-to-byte conversions | `[]byte("literal")` inside loop body |
| `globalprefix` | Prefix Unexported Globals with _ | Top-level unexported `var`/`const` without `_` prefix (exception: `err` prefix for errors) |
| `vartype` | Top-level Variable Declarations | `var x T = expr()` when T matches expr's type |
| `ifaceptr` | Pointers to Interfaces | `*InterfaceType` in function params, struct fields, return types |
| `atomicstd` | Use sync/atomic typed values | Raw `atomic.AddInt32`, `atomic.LoadInt64`, etc. instead of `atomic.Int64` methods |
| `rawstring` | Use Raw String Literals to Avoid Escaping | String literals containing `\"` that could use backtick syntax |

## Architecture

```
uberlint/
  cmd/uberlint/main.go       # standalone runner (optional, for dev)
  plugin.go                   # golangci-lint plugin entry: func NewPlugins() []*analysis.Analyzer
  analyzers/
    nopanic.go                # analyzer + diagnostic message
    chansize.go
    enumstart.go
    newref.go
    nilslice.go
    zerovar.go
    stringbytes.go
    globalprefix.go
    vartype.go
    ifaceptr.go
    atomicstd.go
    rawstring.go
  testdata/
    src/
      nopanic/
        nopanic.go            # test input
      chansize/
        ...
  go.mod
  .golangci.yml               # example config
```

### Analyzer contract

Each analyzer:
- Lives in its own file in `analyzers/`
- Exports a `var Analyzer *analysis.Analyzer` variable
- Has a `run` function that walks the AST
- Uses `pass.Reportf()` to report diagnostics
- Has test data in `testdata/src/<rulename>/`
- Is tested using `golang.org/x/tools/go/analysis/analysistest`

### Plugin entry point

```go
// plugin.go
package uberlint

import (
    "github.com/user/uberlint/analyzers"
    "golang.org/x/tools/go/analysis"
)

func NewPlugins() []*analysis.Analyzer {
    return []*analysis.Analyzer{
        analyzers.NoPanic,
        analyzers.ChanSize,
        // ... all 12 analyzers
    }
}
```

### User integration

```yaml
# .golangci.yml
linters:
  custom:
    uberlint:
      type: "module"
      path: /path/to/uberlint
      description: "Uber Go Style Guide linter"
      original-url: "github.com/user/uberlint"
```

## Rule Details

### nopanic

Flags all `panic()` calls in non-test files. Allows:
- `panic()` in `func init()` (program startup)
- `panic()` in `_test.go` files (but suggests `t.Fatal`)
- `panic()` inside functions named `Must*` (e.g., `template.Must`, `regexp.MustCompile`) — these are conventionally allowed

### chansize

Flags `make(chan T, N)` where N is a compile-time constant > 1. Does NOT flag:
- `make(chan T)` (unbuffered, size 0)
- `make(chan T, 1)`
- `make(chan T, runtime.NumCPU())` (dynamic size, not lintable)

### enumstart

Flags `const` groups using `iota` where the first value is 0 (no `+ 1` or explicit non-zero). Silenced with a comment `// uberskip:enumstart` on the const group. Does NOT flag:
- Explicit zero value (`= 0` when zero is intentional)
- Single const (not a group, likely not an enum)

### newref

Flags `new(T)` expressions. Suggests `&T{}` instead. Does NOT flag:
- `new(T)` when T is an interface (rare but valid)

### nilslice

Flags return statements that return `[]T{}` or `make([]T, 0)` when they could return `nil`. Also flags:
- `s == nil` checks when `len(s) == 0` is more correct
- `s := []T{}` when `var s []T` is preferred

### zerovar

Flags `x := T{}` when T is a struct with all zero-value fields. Suggests `var x T`.

### stringbytes

Flags `[]byte("literal")` expressions inside `for` loop bodies. The fix: hoist the conversion before the loop. Does NOT flag:
- `[]byte(variable)` (dynamic, can't hoist)
- `[]byte("literal")` outside loops (one-time conversion is fine)

### globalprefix

Flags top-level `var` and `const` declarations that are unexported and don't start with `_`. Exception: names starting with `err` (per Error Naming rule). Does NOT flag:
- Exported names (they are public API)
- Names already prefixed with `_`
- Error variables prefixed with `err`

### vartype

Flags `var x T = expr()` when T is the same type as expr's return type. Suggests `var x = expr()`. Does NOT flag when types differ (e.g., `var x error = myError{}`).

### ifaceptr

Flags pointer-to-interface types (`*InterfaceType`). Checks:
- Function parameters
- Struct fields
- Return types
Does NOT flag local variables (too noisy, sometimes legitimate).

### atomicstd

Flags raw `sync/atomic` function calls on types that have typed equivalents in Go 1.19+:
- `atomic.AddInt32(&x, 1)` -> use `atomic.Int32`
- `atomic.LoadInt64(&x)` -> use `atomic.Int64`
- `atomic.SwapUint32(&x, v)` -> use `atomic.Uint32`
- Direct read/write of `int32` field marked with `// atomic` comment

### rawstring

Flags string literals containing escaped quotes (`\"`) that could be written as raw string literals (backticks). Does NOT flag:
- Strings with backticks inside (can't use raw literal)
- Strings with newlines AND escape sequences (context-dependent)
- Very short strings (not worth the change)

## Testing Strategy

Each analyzer gets:
1. **Unit tests** via `analysistest.Run` — test input Go files in `testdata/src/<rule>/`
2. **Positive cases** — code that SHOULD be flagged
3. **Negative cases** — code that should NOT be flagged (exceptions)
4. **Suggested fixes** where applicable (using `analysis.SuggestedFix`)

## Future Phases

**Phase 2:** Medium-difficulty rules:
- Verify Interface Compliance
- Copy Slices and Maps at Boundaries
- Handle Errors Once
- Embedding in Structs (position + separator)
- Use Raw String Literals

**Phase 3:** Hard / semantic rules:
- Defer to Clean Up
- Function Grouping and Ordering
- Reduce Scope of Variables
- Test Tables pattern
- Functional Options pattern

## Success Criteria

- All 12 MVP analyzers pass tests
- Plugin loads correctly in golangci-lint >= 1.53
- Zero false positives on a test corpus of popular Go projects
- Each diagnostic includes the style guide section reference

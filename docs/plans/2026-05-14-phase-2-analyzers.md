# Phase 2 Analyzers Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add eight conservative Uber Go Style Guide analyzers for rules currently listed as uncovered by golangci-lint.

**Architecture:** Each rule follows the existing one-file analyzer plus one `analysistest` test file and one testdata package pattern. The analyzers prefer low false-positive checks: report only when the syntax or type information clearly matches the Uber guidance.

**Tech Stack:** Go, `golang.org/x/tools/go/analysis`, `analysistest`, `golangci-lint` module plugin integration.

### Task 1: Add RED tests

**Files:**
- Create: `analyzers/publicembed_test.go`
- Create: `analyzers/embedlayout_test.go`
- Create: `analyzers/localvar_test.go`
- Create: `analyzers/zerofields_test.go`
- Create: `analyzers/mapinit_test.go`
- Create: `analyzers/constprintf_test.go`
- Create: `analyzers/nakedparams_test.go`
- Create: `analyzers/timefield_test.go`
- Create: `analyzers/testdata/src/<rule>/<rule>.go` for each rule

**Steps:**
- [ ] Add one positive and at least one negative case per rule.
- [ ] Run each new test and verify it fails because the analyzer symbol is missing.

### Task 2: Implement analyzers

**Files:**
- Create: `analyzers/publicembed.go`
- Create: `analyzers/embedlayout.go`
- Create: `analyzers/localvar.go`
- Create: `analyzers/zerofields.go`
- Create: `analyzers/mapinit.go`
- Create: `analyzers/constprintf.go`
- Create: `analyzers/nakedparams.go`
- Create: `analyzers/timefield.go`

**Steps:**
- [ ] Implement `publicembed`: report embedded fields in exported structs.
- [ ] Implement `embedlayout`: report embedded fields that appear after named fields.
- [ ] Implement `localvar`: report function-local `var x = expr`.
- [ ] Implement `zerofields`: report explicit zero-value fields in struct literals.
- [ ] Implement `mapinit`: report empty map literals.
- [ ] Implement `constprintf`: report inline string literals passed as format strings to printf-style functions.
- [ ] Implement `nakedparams`: report calls with two or more naked literal arguments.
- [ ] Implement `timefield`: report serialized numeric fields with time-ish names but no unit suffix.

### Task 3: Wire and verify

**Files:**
- Modify: `plugin.go`
- Modify: `cmd/uberlint/main.go`
- Modify: `.golangci.yml`
- Modify: `README.md`
- Modify: `uber-go-styleguide-linter-coverage.md`

**Steps:**
- [ ] Add all analyzers to plugin and standalone runner.
- [ ] Add linter names to the golangci config and README.
- [ ] Update the coverage document from uncovered to `uberlint`.
- [ ] Run targeted tests and `go test ./... -count=1`.

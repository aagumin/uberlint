# Analyzer Review Fixes Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Eliminate the false positives and broken suppression behavior identified in code review for `newref`, `nilslice`, `rawstring`, and `enumstart`.

**Architecture:** Tighten each analyzer at the decision point where it currently over-matches. Drive every change from a regression test in `analysistest` testdata so the project keeps executable examples of both the old bug and the intended behavior.

**Tech Stack:** Go 1.25, `golang.org/x/tools/go/analysis`, `analysistest`

### Task 1: Add regression coverage for `newref`

**Files:**
- Modify: `analyzers/testdata/src/newref/newref.go`
- Test: `analyzers/newref_test.go`

- [ ] Step 1: Add a negative case showing `new(int)` must not be flagged.
- [ ] Step 2: Run `go test ./analyzers -run TestNewRef -count=1` and verify it fails for the new expectation.
- [ ] Step 3: Update `analyzers/newref.go` so it only reports `new(T)` for struct types.
- [ ] Step 4: Re-run `go test ./analyzers -run TestNewRef -count=1` and verify it passes.

### Task 2: Narrow `nilslice` to safe contexts

**Files:**
- Modify: `analyzers/testdata/src/nilslice/nilslice.go`
- Modify: `analyzers/nilslice.go`
- Test: `analyzers/nilslice_test.go`

- [ ] Step 1: Add negative cases for struct literals and other non-return contexts where non-nil empty slices are intentional.
- [ ] Step 2: Run `go test ./analyzers -run TestNilSlice -count=1` and verify it fails for the new expectation.
- [ ] Step 3: Update `analyzers/nilslice.go` to report only return statements that directly return empty slice literals or `make([]T, 0, ...)`.
- [ ] Step 4: Re-run `go test ./analyzers -run TestNilSlice -count=1` and verify it passes.

### Task 3: Make `rawstring` semantics-safe

**Files:**
- Modify: `analyzers/testdata/src/rawstring/rawstring.go`
- Modify: `analyzers/rawstring.go`
- Test: `analyzers/rawstring_test.go`

- [ ] Step 1: Add a negative case for interpreted strings that contain `\"` plus other escapes such as `\n`.
- [ ] Step 2: Run `go test ./analyzers -run TestRawString -count=1` and verify it fails for the new expectation.
- [ ] Step 3: Update `analyzers/rawstring.go` to unquote the literal and only report when the unquoted value equals the raw literal payload.
- [ ] Step 4: Re-run `go test ./analyzers -run TestRawString -count=1` and verify it passes.

### Task 4: Fix `enumstart` suppression on const groups

**Files:**
- Modify: `analyzers/testdata/src/enumstart/enumstart.go`
- Modify: `analyzers/enumstart.go`
- Test: `analyzers/enumstart_test.go`

- [ ] Step 1: Add a negative case with `//nolint:enumstart` on the `const` group itself.
- [ ] Step 2: Run `go test ./analyzers -run TestEnumStart -count=1` and verify it fails for the new expectation.
- [ ] Step 3: Update `analyzers/enumstart.go` to inspect comments on both the `GenDecl` and the first `ValueSpec`.
- [ ] Step 4: Re-run `go test ./analyzers -run TestEnumStart -count=1` and verify it passes.

### Task 5: Final verification

**Files:**
- Verify: `analyzers/*.go`
- Verify: `analyzers/testdata/src/*/*.go`

- [ ] Step 1: Run the four targeted tests together.
- [ ] Step 2: Run `go test ./... -count=1`.
- [ ] Step 3: Review the diff for accidental scope creep.

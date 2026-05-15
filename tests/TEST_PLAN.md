# Test Plan: uberlint

## Executive Summary

`uberlint` is a Go linter package that exposes 20 Uber Go Style Guide analyzers through the `golangci-lint` module plugin system and a standalone `multichecker` runner. This test plan verifies analyzer correctness, plugin integration, configuration behavior, documentation workflows, and release readiness.

Primary testing goals:

- Verify every shipped analyzer reports the intended violations and avoids obvious false positives.
- Verify `uberlint` works as one `golangci-lint` custom linter named `uberlint`.
- Verify individual analyzer selection works through plugin settings, not through top-level `linters.enable`.
- Verify standalone runner behavior remains useful for analyzer development.
- Verify diagnostics are readable, non-duplicated, suppressible, and stable enough for CI use.

## Scope

### In Scope

- Analyzer unit tests using `golang.org/x/tools/go/analysis/analysistest`.
- Plugin registration and settings tests for `github.com/golangci/plugin-module-register`.
- `golangci-lint custom` integration with `.custom-gcl.yml`.
- `.golangci.yml` integration with the single custom linter name `uberlint`.
- Analyzer `enable` and `disable` settings under `linters.settings.custom.uberlint.settings`.
- Standalone runner execution through `go run ./cmd/uberlint ./...`.
- README quick start, troubleshooting, local path usage, and suppression examples.
- Release smoke checks for versioned module usage.

### Out of Scope

- Full correctness testing of third-party linters listed in `uber-go-styleguide-linter-coverage.md`.
- Style-guide rules explicitly documented as not implemented by `uberlint`.
- Performance benchmarking beyond basic runtime sanity checks.
- Editor-specific integration testing.

## Test Environment

Required tools:

- Go version supported by the project module.
- `git`.
- `golangci-lint` v2 with `golangci-lint custom`.
- Network access for release smoke tests that download `github.com/aagumin/uberlint`.

Recommended local commands:

```bash
go test ./...
go run ./cmd/uberlint ./analyzers/testdata/src/zerovar
golangci-lint custom -v
./custom-gcl cache clean
./custom-gcl run ./...
```

Recommended isolated cache for deterministic local runs:

```bash
GOCACHE="$(pwd)/.gocache" go test ./...
GOCACHE="$(pwd)/.gocache" ./custom-gcl run ./...
```

## Entry Criteria

- `README.md` documents the current module path, plugin setup, rule selection, and troubleshooting.
- `uber-go-styleguide-linter-coverage.md` lists the current analyzer set.
- All analyzer source files compile.
- Testdata exists for every shipped analyzer.
- `.custom-gcl.yml` and `.golangci.yml` are available for integration testing.

## Exit Criteria

- All Go tests pass with `go test ./...`.
- Every shipped analyzer has positive and negative testdata coverage.
- The custom `golangci-lint` binary builds successfully.
- `./custom-gcl run ./...` recognizes `uberlint` and does not require individual analyzer names in top-level `linters.enable`.
- Analyzer selection through `settings.enable` and `settings.disable` works.
- Diagnostics do not duplicate analyzer names, for example `zerovar: zerovar: ...`.
- Suppression through `//nolint:<analyzer>` works for representative analyzers.
- All P0 and P1 test cases below pass before release.

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|---|---:|---:|---|
| `golangci-lint` module plugin API changes | Medium | High | Pin custom builder version in `.custom-gcl.yml`; run integration smoke on every release. |
| Analyzer false positives on real projects | High | High | Maintain negative testdata and test against at least one external fixture project before release. |
| Analyzer false negatives due to type information gaps | Medium | Medium | Keep `GetLoadMode()` at `register.LoadModeTypesInfo`; add type-aware testdata. |
| Rule selection misconfigured by users | Medium | Medium | Test README examples and error messages for unknown analyzer names. |
| Stale `golangci-lint` cache shows old diagnostics | Medium | Medium | Include cache-clean step in troubleshooting and release validation. |
| Retagged or cached module versions cause old plugin code to be used | Low | High | Never move published tags; release a new semver tag for every user-visible change. |

## Analyzer Coverage Matrix

| Analyzer | Style Guide Rule | Priority | Unit Test Focus |
|---|---|---:|---|
| `ifaceptr` | Pointers to Interfaces | P0 | Function params, returns, fields, aliases, non-interface pointers. |
| `chansize` | Channel Size is One or None | P0 | `make(chan T, 2+)` reports; `0`, `1`, dynamic sizes do not over-report unless intentionally supported. |
| `enumstart` | Start Enums at One | P1 | `iota` starting at zero reports; explicit non-zero start and documented suppression pass. |
| `nopanic` | Don't Panic | P0 | `panic()` in production code reports; allowed generated/test scenarios behave as documented. |
| `atomicstd` | Use sync/atomic typed values | P1 | Raw `sync/atomic` function calls report; typed atomic values do not. |
| `stringbytes` | Avoid repeated string-to-byte conversions | P2 | Literal conversions in loops report; outside loops and variable conversions do not over-report. |
| `vartype` | Top-level Variable Declarations | P1 | Redundant top-level var types report; needed explicit types pass. |
| `globalprefix` | Prefix Unexported Globals with `_` | P1 | Unexported globals report; constants, exported names, and accepted error patterns pass. |
| `nilslice` | nil is a valid slice | P1 | Empty slice returns report; non-empty slices and non-return contexts pass. |
| `rawstring` | Use Raw String Literals | P2 | Escaped quote strings report when backticks preserve semantics; short or incompatible strings pass. |
| `zerovar` | Use `var` for Zero Value Structs | P1 | `x := T{}` reports; non-empty literals and non-struct literals pass. |
| `newref` | Initializing Struct References | P1 | `new(T)` for structs reports; non-struct `new` uses pass if intended. |
| `publicembed` | Avoid Embedding Types in Public Structs | P1 | Exported structs with embedded fields report; private structs pass if intended. |
| `embedlayout` | Embedding in Structs | P2 | Embedded fields after named fields report; grouped embedded fields pass. |
| `localvar` | Local Variable Declarations | P2 | `var x = value` inside functions reports; package-level vars pass. |
| `zerofields` | Omit Zero Value Fields in Structs | P1 | Explicit zero fields report; non-zero fields and semantically meaningful values pass. |
| `mapinit` | Initializing Maps | P1 | Empty map literals report; non-empty map literals pass. |
| `constprintf` | Format Strings outside Printf | P2 | Non-const format variables report; inline literals and const format strings pass. |
| `nakedparams` | Avoid Naked Parameters | P2 | Multiple naked literal args report; named constants or commented literals pass if supported. |
| `timefield` | Time Units in Serialized Fields | P1 | Serialized numeric time-like fields without units report; unit-suffixed fields pass. |

## Smoke Test Suite

### SMOKE-001: Project Test Suite

**Priority:** P0  
**Type:** Automated / Regression  
**Estimated Time:** 1-3 minutes

#### Objective

Verify the repository builds and all analyzer/plugin tests pass.

#### Steps

1. Run:

   ```bash
   go test ./...
   ```

   **Expected:** All packages pass. No compile errors, failing tests, or analyzer test mismatches.

### SMOKE-002: Standalone Runner Produces a Diagnostic

**Priority:** P1  
**Type:** Automated / Functional  
**Estimated Time:** 1 minute

#### Objective

Verify the standalone multichecker can run analyzers against testdata.

#### Steps

1. Run:

   ```bash
   go run ./cmd/uberlint ./analyzers/testdata/src/zerovar
   ```

   **Expected:** Command exits non-zero and reports a `zerovar` diagnostic for `User{}`.

2. Inspect the diagnostic text.

   **Expected:** The message does not duplicate the analyzer name, for example it must not contain `zerovar: zerovar:`.

### SMOKE-003: Custom golangci-lint Binary Builds

**Priority:** P0  
**Type:** Integration  
**Estimated Time:** 2-5 minutes

#### Objective

Verify the documented module plugin build works.

#### Preconditions

- `.custom-gcl.yml` exists.
- Network access is available unless `.custom-gcl.yml` uses `path`.

#### Steps

1. Remove an old binary:

   ```bash
   rm -f ./custom-gcl
   ```

   **Expected:** No old custom binary remains.

2. Build the binary:

   ```bash
   golangci-lint custom -v
   ```

   **Expected:** Build succeeds and generated imports include `_ "github.com/aagumin/uberlint"`.

3. Inspect module metadata:

   ```bash
   go version -m ./custom-gcl
   ```

   **Expected:** Output includes `github.com/aagumin/uberlint`.

### SMOKE-004: Custom golangci-lint Runs uberlint

**Priority:** P0  
**Type:** Integration  
**Estimated Time:** 1-3 minutes

#### Objective

Verify `golangci-lint` recognizes the custom linter as `uberlint`.

#### Steps

1. Run:

   ```bash
   ./custom-gcl cache clean
   ./custom-gcl run ./...
   ```

   **Expected:** The run starts without `plugin "uberlint" not found`.

2. If findings are present, inspect several `uberlint` diagnostics.

   **Expected:** Findings use the format `<analyzer>: <message> (uberlint)` without duplicated analyzer names.

## Functional Test Cases

### TC-001: All Analyzer Unit Tests Pass

**Priority:** P0  
**Type:** Automated / Functional  
**Status:** Not Run

#### Objective

Verify each analyzer matches its `analysistest` expectations.

#### Preconditions

- Testdata exists under `analyzers/testdata/src/<analyzer>`.

#### Steps

1. Run:

   ```bash
   go test ./analyzers -count=1
   ```

   **Expected:** All analyzer tests pass.

2. Confirm every analyzer listed in README has a matching `*_test.go`.

   **Expected:** There is one test file for each analyzer in the analyzer coverage matrix.

### TC-002: Plugin Registers as `uberlint`

**Priority:** P0  
**Type:** Automated / Integration  
**Status:** Not Run

#### Objective

Verify module plugin registration is available to `golangci-lint`.

#### Steps

1. Run:

   ```bash
   go test . -run TestPluginRegistersWithGolangCILint -count=1
   ```

   **Expected:** Test passes and reports exactly 20 analyzers.

### TC-003: Top-Level Analyzer Names Are Not Accepted as golangci-lint Linters

**Priority:** P1  
**Type:** Manual / Negative Integration  
**Status:** Not Run

#### Objective

Verify documentation is correct: users must enable `uberlint`, not individual analyzer IDs, under top-level `linters.enable`.

#### Test Data

Use a temporary `.golangci.yml` with:

```yaml
version: "2"
linters:
  default: none
  enable:
    - zerovar
```

#### Steps

1. Run:

   ```bash
   ./custom-gcl run ./...
   ```

   **Expected:** The run fails with an unknown linter error for `zerovar`.

2. Replace `zerovar` with `uberlint`.

   **Expected:** The run no longer fails because of an unknown linter name.

### TC-004: `settings.enable` Runs Only Selected Analyzers

**Priority:** P0  
**Type:** Automated / Integration  
**Status:** Not Run

#### Objective

Verify users can gradually adopt individual rules through plugin settings.

#### Test Data

Use:

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
        description: "Uber Go Style Guide linter"
        settings:
          enable:
            - zerovar
```

#### Steps

1. Run against a fixture package containing multiple violations from different analyzers.

   **Expected:** Only `zerovar` findings are reported.

2. Add `mapinit` to `settings.enable`.

   **Expected:** `zerovar` and `mapinit` findings are reported; unrelated analyzer findings are not reported.

### TC-005: `settings.disable` Excludes Selected Analyzers

**Priority:** P0  
**Type:** Automated / Integration  
**Status:** Not Run

#### Objective

Verify users can run all analyzers except explicitly disabled ones.

#### Test Data

Use:

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
        description: "Uber Go Style Guide linter"
        settings:
          disable:
            - zerovar
```

#### Steps

1. Run against a fixture package containing `zerovar` and at least one other `uberlint` violation.

   **Expected:** `zerovar` findings are absent and other enabled analyzer findings are present.

### TC-006: Invalid Analyzer Settings Fail Fast

**Priority:** P1  
**Type:** Automated / Negative Functional  
**Status:** Not Run

#### Objective

Verify invalid analyzer names are rejected with actionable errors.

#### Steps

1. Run:

   ```bash
   go test . -run TestPluginSettingsRejectUnknownAnalyzer -count=1
   ```

   **Expected:** Test passes and validates an `unknown analyzer` error.

2. Manually configure `settings.enable: ["doesnotexist"]` in `.golangci.yml` and run:

   ```bash
   ./custom-gcl run ./...
   ```

   **Expected:** The run fails with an explicit unknown analyzer error.

### TC-007: `enable` and `disable` Cannot Be Used Together

**Priority:** P1  
**Type:** Automated / Negative Functional  
**Status:** Not Run

#### Objective

Verify ambiguous rule selection is rejected.

#### Steps

1. Configure both:

   ```yaml
   settings:
     enable:
       - zerovar
     disable:
       - mapinit
   ```

2. Run:

   ```bash
   ./custom-gcl run ./...
   ```

   **Expected:** The run fails with an error equivalent to `use either enable or disable, not both`.

### TC-008: Diagnostics Do Not Duplicate Analyzer Names

**Priority:** P0  
**Type:** Regression  
**Status:** Not Run

#### Objective

Prevent regressions where messages appear as `zerovar: zerovar: ...`.

#### Steps

1. Run:

   ```bash
   ./custom-gcl cache clean
   ./custom-gcl run ./analyzers/testdata/src/zerovar
   ```

   **Expected:** Output contains `zerovar: use var declaration` and does not contain `zerovar: zerovar:`.

2. Repeat for at least three more analyzers, for example `mapinit`, `publicembed`, and `constprintf`.

   **Expected:** No diagnostic contains `<analyzer>: <analyzer>:` for any analyzer.

### TC-009: Suppression Comments Work

**Priority:** P1  
**Type:** Integration / Functional  
**Status:** Not Run

#### Objective

Verify users can suppress intentional exceptions through normal `golangci-lint` comments.

#### Test Data

Create a temporary fixture with:

```go
package fixture

type User struct{}

func example() {
	//nolint:zerovar // testing suppression for intentional style exception
	u := User{}
	_ = u
}
```

#### Steps

1. Run:

   ```bash
   ./custom-gcl run ./path/to/fixture
   ```

   **Expected:** No `zerovar` issue is reported for the suppressed line.

2. Remove the suppression comment and run again.

   **Expected:** `zerovar` issue is reported.

### TC-010: README Quick Start Works in a Fresh External Project

**Priority:** P0  
**Type:** Manual / End-to-End  
**Status:** Not Run

#### Objective

Verify a new user can follow README instructions successfully.

#### Preconditions

- Create a temporary Go module outside this repository.
- Add `.custom-gcl.yml` and `.golangci.yml` exactly as documented.

#### Steps

1. Add a small Go file with a known violation, for example `x := T{}`.

   **Expected:** The fixture builds as a valid Go package.

2. Run:

   ```bash
   golangci-lint custom -v
   ./custom-gcl cache clean
   ./custom-gcl run ./...
   ```

   **Expected:** The custom binary builds, `uberlint` runs, and the known violation is reported.

3. Change `.custom-gcl.yml` from version mode to local `path` mode.

   **Expected:** The custom binary still builds and uses the local checkout.

### TC-011: Stale Cache Troubleshooting Works

**Priority:** P1  
**Type:** Manual / Regression  
**Status:** Not Run

#### Objective

Verify README troubleshooting resolves stale diagnostics after plugin upgrades.

#### Steps

1. Run a custom binary that reports at least one `uberlint` issue.

   **Expected:** `golangci-lint` cache is populated.

2. Rebuild the custom binary after changing the plugin version or local path.

   **Expected:** New binary is created.

3. Run:

   ```bash
   ./custom-gcl cache clean
   ./custom-gcl run ./...
   ```

   **Expected:** Output reflects current diagnostics and does not show stale duplicated messages.

### TC-012: Release Tag Smoke Test

**Priority:** P0  
**Type:** Release / End-to-End  
**Status:** Not Run

#### Objective

Verify a published module version can be consumed by another project.

#### Preconditions

- A new Git tag has been pushed.
- Go module proxy or direct Git access can resolve the tag.

#### Steps

1. In a temporary project, set:

   ```yaml
   version: v2.0.0
   name: custom-gcl
   plugins:
     - module: github.com/aagumin/uberlint
       version: <new-version>
   ```

2. Run:

   ```bash
   golangci-lint custom -v
   go version -m ./custom-gcl
   ./custom-gcl cache clean
   ./custom-gcl run ./...
   ```

   **Expected:** The binary includes `github.com/aagumin/uberlint` at the intended version and runs without plugin registration errors.

## Analyzer-Specific Regression Checklist

For every analyzer in the coverage matrix, maintain these test categories:

- Positive case: at least one style-guide violation must be reported.
- Negative case: idiomatic code recommended by the style guide must not be reported.
- Boundary case: minimum/maximum or ambiguous inputs must be covered when applicable.
- Type-aware case: if the analyzer depends on type info, include aliases, named types, or package-qualified types.
- Suppression case: representative findings must be suppressible through `//nolint:<analyzer>` when run through `golangci-lint`.
- Diagnostic case: message must be clear, stable, and must not include a duplicated analyzer prefix.

## Regression Suite

### Fast Regression

Run before every commit that changes analyzers, plugin registration, or testdata.

```bash
go test ./...
```

Expected result: all tests pass.

### Integration Regression

Run before merging changes that affect README instructions, `.custom-gcl.yml`, `.golangci.yml`, `plugin.go`, or release packaging.

```bash
rm -f ./custom-gcl
golangci-lint custom -v
./custom-gcl cache clean
./custom-gcl run ./...
```

Expected result: custom binary builds, recognizes `uberlint`, and produces current diagnostics.

### Release Regression

Run before publishing a new tag.

```bash
go test ./...
rm -f ./custom-gcl
golangci-lint custom -v
go version -m ./custom-gcl
./custom-gcl cache clean
./custom-gcl run ./...
```

Expected result:

- Tests pass.
- Custom binary contains `github.com/aagumin/uberlint`.
- `uberlint` is registered as a single custom linter.
- Diagnostics are current and non-duplicated.
- README quick start still matches actual behavior.

## Test Data Requirements

Each analyzer fixture under `analyzers/testdata/src/<analyzer>` should include:

- A minimal violating example with `// want` assertion.
- A compliant example near the violation to prevent broad false positives.
- At least one realistic example shaped like production code.
- Comments explaining intentional edge cases when the expected behavior is not obvious.

For integration tests, maintain at least one external-style fixture project with:

- A `go.mod`.
- `.custom-gcl.yml`.
- `.golangci.yml`.
- A package with violations from at least five analyzers.
- A package with no expected `uberlint` violations.

## Reporting Template

Use this template for each test run:

```markdown
# Test Run: uberlint <version or commit>

Date:
Tester:
Go version:
golangci-lint version:
OS:
Commit/tag:

## Summary

Total test cases:
Passed:
Failed:
Blocked:
Not run:

## Failures

| Test Case | Severity | Actual Result | Link/Evidence |
|---|---|---|---|
| | | | |

## Release Recommendation

Go / No-Go:
Reason:
Required fixes:
```

## Bug Report Template

```markdown
# BUG: [Analyzer or Integration Area] short title

Severity: Critical | High | Medium | Low
Priority: P0 | P1 | P2 | P3
Type: Analyzer False Positive | Analyzer False Negative | Plugin Integration | Documentation | Release

## Environment

- OS:
- Go version:
- golangci-lint version:
- uberlint commit/tag:
- Run mode: analysistest | standalone | custom-gcl

## Description

Describe the issue and why it violates expected behavior.

## Steps to Reproduce

1.
2.
3.

## Expected Result

What should happen.

## Actual Result

What happened instead.

## Evidence

Command output, fixture code, screenshots, or links.

## Impact

User impact, CI impact, and workaround if available.
```

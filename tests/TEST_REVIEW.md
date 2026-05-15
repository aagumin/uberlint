# Test Review: uberlint

## Current Automated Coverage

The project already has baseline unit coverage for all 20 shipped analyzers. Each analyzer has:

- A dedicated `analyzers/<name>_test.go` file.
- A matching `analyzers/testdata/src/<name>/<name>.go` fixture.
- At least one positive `// want` assertion.
- At least one compliant example for the common non-violation path.

Covered analyzers:

| Analyzer | Unit Test | Testdata |
|---|---|---|
| `atomicstd` | Yes | Yes |
| `chansize` | Yes | Yes |
| `constprintf` | Yes | Yes |
| `embedlayout` | Yes | Yes |
| `enumstart` | Yes | Yes |
| `globalprefix` | Yes | Yes |
| `ifaceptr` | Yes | Yes |
| `localvar` | Yes | Yes |
| `mapinit` | Yes | Yes |
| `nakedparams` | Yes | Yes |
| `newref` | Yes | Yes |
| `nilslice` | Yes | Yes |
| `nopanic` | Yes | Yes |
| `publicembed` | Yes | Yes |
| `rawstring` | Yes | Yes |
| `stringbytes` | Yes | Yes |
| `timefield` | Yes | Yes |
| `vartype` | Yes | Yes |
| `zerofields` | Yes | Yes |
| `zerovar` | Yes | Yes |

The root plugin package also has unit coverage for:

- Module plugin registration as `uberlint`.
- Building the analyzer list.
- `settings.enable`.
- `settings.disable`.
- Unknown analyzer rejection for `settings.enable`.

## Gaps Found

The analyzer-level test layout is complete, but several plugin and regression contracts were not covered:

- Exact shipped analyzer ID set was not locked by a test.
- `GetLoadMode()` was not tested, even though several analyzers require type information.
- `settings.enable` and `settings.disable` conflict handling was not tested.
- Unknown analyzer rejection for `settings.disable` was not tested.
- The recent diagnostic format regression was not protected by an automated test.

## Tests Added

Added root plugin tests:

- `TestNewPluginsReturnsExpectedAnalyzerSet`
- `TestPluginRequestsTypesInfoLoadMode`
- `TestPluginSettingsRejectEnableAndDisableTogether`
- `TestPluginSettingsRejectUnknownDisabledAnalyzer`

Added analyzer regression test:

- `TestDiagnosticsDoNotIncludeAnalyzerNamePrefix`

Added analyzer inventory tests:

- `TestAnalyzersHaveUniqueNamesAndMetadata`
- `TestEveryAnalyzerHasTestdata`

This protects all analyzer diagnostics from returning to the duplicated format:

```text
zerovar: zerovar: use var declaration...
```

Expanded `analysistest` fixtures with additional unit-testable cases:

- `atomicstd`: aliased `sync/atomic`, CAS/swap operations, non-atomic selector false positive guard.
- `chansize`: constant expressions, const identifiers, non-channel `make` false positive guard.
- `constprintf`: selector and local printf-style calls, non-printf call guard.
- `embedlayout`: pointer embedded fields and multiple embedded fields.
- `enumstart`: `iota + 0`, same-line suppression, non-iota const groups.
- `ifaceptr`: pointer-to-interface in type aliases, top-level vars, and local vars.
- `localvar`: grouped local var declarations and no-initializer guard.
- `mapinit`: named map types.
- `nakedparams`: boolean/nil/zero literal combinations and single-literal guard.
- `newref`: named struct aliases.
- `nilslice`: `make([]T, 0, cap)` returns and allocated non-empty buffer guard.
- `publicembed`: embedded pointer fields in exported structs.
- `stringbytes`: range loops.
- `timefield`: float/int64 serialized time-like fields and non-serialized/non-numeric guards.
- `zerofields`: float and imaginary zero values.
- `zerovar`: named struct aliases.

New tests exposed and fixed these implementation gaps:

- `enumstart` now reports enum groups starting with `iota + 0`.
- `ifaceptr` now reports pointer-to-interface declarations in `var` and `type` declarations.
- `mapinit` now reports empty literals for named map types.

## Can Everything Be Unit Tested?

No. Most analyzer logic should be unit-tested with `analysistest`, and plugin selection logic can be unit-tested directly. However, several important behaviors require integration or release smoke tests.

Good fit for unit tests:

- Individual analyzer true positives and false positives.
- Type-aware analyzer behavior using `analysistest`.
- Analyzer list membership and ordering.
- Plugin settings validation.
- Diagnostic message invariants.
- Standalone helper logic, if added later.

Requires integration tests:

- `golangci-lint custom -v` actually building a custom binary.
- `./custom-gcl run ./...` discovering the module plugin.
- Top-level `.golangci.yml` behavior where only `uberlint` is a valid custom linter name.
- `//nolint:<analyzer>` behavior in the real `golangci-lint` runner.
- Cache behavior after plugin upgrades.
- Versioned module consumption from `github.com/aagumin/uberlint@<tag>`.

Recommended split:

- Keep analyzer correctness in fast unit tests.
- Keep plugin settings in root package unit tests.
- Add a small scripted integration test for `golangci-lint custom` and README quickstart.
- Run release smoke tests manually or in CI before publishing a tag.

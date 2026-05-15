# uberlint

Uber Go Style Guide, enforced by `golangci-lint`.

> **Disclaimer:** `uberlint` is an independent open source project. It is not affiliated with, endorsed by, sponsored by, or connected to Uber Technologies, Inc.

`uberlint` is a focused set of Go analyzers for teams that like the Uber Go Style Guide but do not want to rely on memory, review comments, or tribal knowledge to enforce it. It catches the style-guide rules that common `golangci-lint` setups usually miss: panic usage, oversized channel buffers, zero-valued enums, pointer-to-interface types, raw `sync/atomic` calls, and more.

The goal is simple: make the preferred Go style automatic, reviewable, and CI-friendly.

## Why Use It

Code review should be about design, correctness, and trade-offs, not repeatedly asking for `&T{}` instead of `new(T)` or for enum values to start at `iota + 1`.

`uberlint` gives you:

- **Uber Go Style Guide coverage** for rules that are missing or only partially covered by standard linters.
- **Native `golangci-lint` integration** through the module plugin system.
- **Small, explainable analyzers** built on `golang.org/x/tools/go/analysis`.
- **Rule selection through plugin settings**, so teams can adopt strictness gradually.
- **Standalone runner support** for development and debugging outside `golangci-lint`.

## Included Linters

| Linter | Checks |
|---|---|
| `ifaceptr` | Avoid pointers to interfaces. |
| `chansize` | Keep channel buffer sizes at `0` or `1`. |
| `enumstart` | Start iota-style enums at a non-zero value. |
| `nopanic` | Avoid `panic()` in production code. |
| `atomicstd` | Prefer typed `sync/atomic` values over raw atomic functions. |
| `stringbytes` | Avoid repeated `[]byte("literal")` conversions inside loops. |
| `vartype` | Omit redundant top-level var types when the initializer already provides the type. |
| `globalprefix` | Prefix unexported top-level globals with `_`. |
| `nilslice` | Return `nil` instead of empty slices where `nil` is sufficient. |
| `rawstring` | Prefer raw string literals when they avoid quote escaping without changing semantics. |
| `zerovar` | Use `var x T` for zero-value structs instead of `x := T{}`. |
| `newref` | Prefer `&T{}` over `new(T)` for struct references. |
| `publicembed` | Avoid embedded fields in exported structs. |
| `embedlayout` | Keep embedded fields grouped before named fields. |
| `localvar` | Prefer `:=` for local variables with inferred types. |
| `zerofields` | Omit explicit zero-value fields in struct literals. |
| `mapinit` | Prefer `make(map[K]V)` for empty map initialization. |
| `constprintf` | Use constants for printf format strings stored outside call sites. |
| `nakedparams` | Avoid calls with multiple naked literal arguments. |
| `timefield` | Require unit suffixes for serialized numeric time fields. |

## Quick Start

Use `uberlint` in your Go project as a `golangci-lint` module plugin. The setup has two files:

- `.custom-gcl.yml` tells `golangci-lint` how to build a custom binary with `uberlint` inside.
- `.golangci.yml` enables the `uberlint` rules your project wants to run.

After that, you run `./custom-gcl run ./...` instead of the stock `golangci-lint run ./...`.

Requirements:

- Go
- git
- `golangci-lint` v2 with `golangci-lint custom`

In the Go project you want to lint, add `.custom-gcl.yml`:

```yaml
version: v2.0.0
name: custom-gcl
plugins:
  - module: github.com/aagumin/uberlint
    version: v0.1.4
```

If you are trying a local checkout before a release is available, use `path` instead:

```yaml
version: v2.0.0
name: custom-gcl
plugins:
  - module: github.com/aagumin/uberlint
    path: /absolute/path/to/uberlint
```

Then add `.golangci.yml`:

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
```

Build and run:

```bash
golangci-lint custom -v
./custom-gcl run ./...
```

The first command reads `.custom-gcl.yml` and builds a local `custom-gcl` binary. The second command runs that binary with the `uberlint` plugin enabled by `.golangci.yml`.

Important: enable `uberlint`, not individual analyzer names like `nopanic` or `ifaceptr`. The module plugin registers one golangci-lint linter named `uberlint`; that linter runs the analyzers listed below.

## Selecting Rules

`golangci-lint` sees one custom linter named `uberlint`. Individual checks are selected through plugin settings.

To run only specific analyzers:

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
            - nopanic
            - ifaceptr
            - enumstart
```

To run all analyzers except a few:

```yaml
version: "2"

linters:
  enable:
    - uberlint
  settings:
    custom:
      uberlint:
        type: "module"
        description: "Uber Go Style Guide linter"
        settings:
          disable:
            - nakedparams
            - constprintf
```

Use either `enable` or `disable`, not both. Unknown analyzer names fail the run with an explicit error.

## Troubleshooting

### `plugin "uberlint" not found`

This means the `custom-gcl` binary was built without the `uberlint` module plugin registration. Check these points:

- `.golangci.yml` must enable `uberlint`, not individual analyzer names:

```yaml
linters:
  enable:
    - uberlint
  settings:
    custom:
      uberlint:
        type: "module"
        description: "Uber Go Style Guide linter"
```

- Rebuild the custom binary after changing `.custom-gcl.yml` or upgrading `uberlint`:

```bash
rm -f ./custom-gcl
golangci-lint custom -v
./custom-gcl run ./...
```

- If diagnostics still look stale after upgrading, for example `zerovar: zerovar: ...`,
  clear the `golangci-lint` cache and run again:

```bash
./custom-gcl cache clean
./custom-gcl run ./...
```

- If you are testing a local checkout, make sure `.custom-gcl.yml` points to that checkout with `path`.

## CI Example

For CI, build the custom binary and run it like any other `golangci-lint` command:

```yaml
name: lint

on:
  pull_request:
  push:
    branches: [main]

jobs:
  golangci-lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"

      - name: Install golangci-lint
        run: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

      - name: Build custom golangci-lint
        run: golangci-lint custom -v

      - name: Run linters
        run: ./custom-gcl run ./...
```

For faster CI, cache Go modules and the `golangci-lint` cache using your CI provider's normal Go caching setup.

## Standalone Runner

For development, you can also run `uberlint` without `golangci-lint`:

```bash
go run ./cmd/uberlint ./...
```

This is useful when developing new analyzers or debugging a rule. For normal project adoption, prefer the `golangci-lint` module plugin path so output, CI behavior, exclusions, and editor integrations stay consistent with the rest of your lint setup.

## Suppressing Findings

Use normal `golangci-lint` suppression comments when a rule is intentionally not applicable:

```go
//nolint:enumstart // zero is the wire-compatible default
const (
	Unknown Status = iota
	Ready
	Failed
)
```

Keep suppressions narrow and documented. A suppression should explain why this instance is exceptional, not why the rule is inconvenient.

## Current Status

`uberlint` currently ships 20 analyzers covering the first two implementation slices of the Uber Go Style Guide. The analyzers are tested with `analysistest`, and the project also includes a standalone `multichecker` binary for local development.

Run the project tests with:

```bash
go test ./...
```

## Links

- Uber Go Style Guide: https://github.com/uber-go/guide/blob/master/style.md
- golangci-lint module plugin docs: https://golangci-lint.run/docs/plugins/module-plugins/

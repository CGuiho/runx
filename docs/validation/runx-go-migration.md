---
name: RunX Go Migration Validation
purpose: Record completion evidence for the production Go/Cobra migration.
description: Formatting, module integrity, tests, vet, CLI smokes, installer syntax, XDocs, Mirror planning, and the standard cross-build asset matrix.
created: 2026-07-26
owner: runx-validation
flags: []
tags:
  - validation
  - go
  - cli
keywords:
  - RunX Go migration
  - 11 artifacts
  - RFC 0034
---

# RunX Go Migration Validation

## Scope

Validated the saved `main` checkout after replacing Bun/TypeScript production,
CI, installer, and publication authority with the native Go/Cobra
implementation. No version bump, tag, push, package publication, GitHub Release,
or approval of the waiting 0.8.0 workflow occurred.

## Results

| Check | Result |
| --- | --- |
| `gofmt` verification | Passed for `main.go`, `cmd`, `pkg`, `embed`, and Go devops programs. |
| `go mod tidy` clean-diff verification | Passed. |
| `go test ./...` | Passed for all production packages. |
| `go vet ./...` | Passed. |
| Native Windows AMD64 `--version` | Passed: `0.8.0`. |
| Native no-argument welcome | Passed: `Hello Windows - runx v0.8.0`. |
| Bare-command agent bootstrap | Passed in unit and compiled-native smoke checks: both global skill copies; both/single/new repository instruction targets; CRLF and unmanaged-content preservation; malformed-marker rejection; idempotent second run; help/version/agent/uninstall exclusions. |
| Native help tree | Passed with live Unicode Cobra hierarchy and depth limiting. |
| Native `init`, `check --format json`, `list --format json` | Passed in an isolated directory; configuration paths were absolute and JSON was deterministic. |
| Bash installer syntax | Passed with Git Bash `bash -n`. |
| PowerShell installer syntax | Passed with the PowerShell parser. |
| `go run devops/build-binaries.go ...` | Passed for all eight `CGO_ENABLED=0` targets. |
| `go run devops/verify-release-assets.go` | Passed: exactly 11 assets, 10 deterministic checksum entries, and a valid skill ZIP. |
| Foreign target runtime | Build-only locally; native CI runners are responsible for runtime smoke coverage. |
| XDocs strict metadata/tree/doctor | Passed; doctor reported zero errors and zero warnings. |
| Mirror config/current/minor plan | Passed read-only: current `0.8.0`; minor plans `0.9.0` and tag `@guiho/runx/v0.9.0`. Nothing was applied. |

## Release Matrix Decision

The shared `guiho-s-0035-cli-engineer-go` contract is authoritative: Linux
AMD64/ARM64/ARMv7/ARMv6, Darwin AMD64/ARM64, Windows AMD64/ARM64,
`guiho-s-runx.zip`, `guiho-i-runx.md`, and `checksums.txt`. The retired x64
default/baseline/modern assets and AMD64 V3 builds are not produced.

## Boundary

The implementation passed the local migration and release-contract gate.
Mirror application and every external side effect remain explicitly deferred.

---
name: RunX Windows Git Bash Automatic Shell Implementation
purpose: Record RX-SHELL-1 implementation choices, regression coverage, validation evidence, and delivery boundaries.
description: Tracks caller-aware Windows shell:auto resolution and its safe fallback contract without changing configured commands or explicit shells.
created: "2026-08-11"
flags:
  - completed
tags:
  - todo
  - implementation
  - windows
  - execution
keywords:
  - RX-SHELL-1
  - Git Bash
  - MSYSTEM
  - shell auto
  - issue 47
owner: runx-todo-done
---

# RunX Windows Git Bash Automatic Shell Implementation

## Summary

RX-SHELL-1 makes Windows `shell: auto` caller-aware. A recognized `MSYSTEM`
Git Bash/MSYS marker and a resolved non-System32 Bash executable select that
exact executable for the Bash transport. Missing or unverified caller evidence,
missing Bash, and the Windows System32/WSL launcher all fall back to the
existing `cmd.exe` transport. Non-Windows automatic selection remains `sh`, and
explicit shell values remain authoritative.

## Design

- `BuildShellExecution` remains the public production entry point.
- A private `shellResolutionDeps` seam injects the runtime platform,
  environment lookup, and executable lookup so Windows cases are deterministic
  on every test platform.
- Only recognized `MSYSTEM` values identify MSYS/Git Bash callers. Automatic
  lookup asks for `bash` and rejects System32/Sysnative `bash.exe` and `wsl.exe`
  launchers before selecting it.
- The resolved executable path is used as the Bash program without rewriting
  configured command text, child arguments, or paths.

## Regression Coverage

Focused executor tests cover resolved Git Bash selection, the non-secret
`/c/GUIHO` argument boundary, ordinary Windows fallback, missing Bash, empty
resolution, System32/Sysnative/WSL launchers, explicit-shell precedence,
non-Windows automatic selection, and recognized/unrecognized markers.

## Progress

- `in progress`: RX-SHELL-1 was claimed from planning head `71868a6` in the
  isolated `codex/runx-git-bash-auto` worktree.
- `testing`: implementation and focused regressions were committed in
  `10774b1`; public docs, canonical/embedded skill guidance, and executor
  XDocs were committed in `bcafba5`.

## Validation Evidence

- `gofmt -w pkg/executor/executor.go pkg/executor/executor_test.go pkg/executor/shell_resolution_test.go` - passed.
- `go mod tidy` followed by a clean `go.mod`/`go.sum` diff and `git diff --check` - passed.
- `go test -count=1 ./pkg/executor` and `go test -count=1 ./...` - passed.
- `go vet ./...` and `go build ./...` - passed; the build required the
  repository-local safe-directory Git metadata configuration in the isolated
  worktree and was not run with `-buildvcs=false`.
- `go run devops/build-binaries.go --version 0.12.1 --commit bcafba5
  --build-date 2026-08-11T00:00:00Z` - built eight targets and three support
  artifacts; `go run devops/verify-release-assets.go` confirmed exactly 11
  assets and every SHA-256 checksum.
- `xdocs scan`, strict `xdocs meta . --documents --strict`, `xdocs tree`, and
  `xdocs doctor --warnings-as-errors` - passed with zero errors and warnings.
- `mirror config check` - passed; `mirror version plan patch` resolved
  `0.12.0` to `0.12.1` with tag `@guiho/runx/v0.12.1`. No version apply, tag,
  release, push, or GitHub mutation occurred.

## Delivery Boundary

Task 22 completed after independent review, exact-head validation, protected
integration, and the Mirror-managed `0.12.1` release. Environment and key
paths remained opaque throughout.

## References

- [Task specification](./windows-git-bash-auto-shell.md)
- [Implementation plan](../../plans/windows-git-bash-auto-shell.md)
- [GitHub issue 47](https://github.com/CGuiho/runx/issues/47)

## Final Acceptance

PR #55 merged at `a297aa25b94d83d11d361a0fab0b5415bdd1ba20`.
Release-note PR #56 produced release source
`2c96e8aec1052b35c6fb4061eba5a6b54405bbcb`. RunX 0.12.1 published with
exactly 11 assets; the downloaded Windows AMD64 checksum/version and released
Git Bash path probe passed, npm latest is 0.12.1, and issue 47 is closed.

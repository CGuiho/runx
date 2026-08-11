---
name: RunX Windows Git Bash Automatic Shell Review
purpose: Materialize the independent exact-head review of the issue 47 Windows automatic-shell correction.
description: Accepts caller-aware Git Bash selection, cmd fallback, explicit-shell precedence, raw command transport, and deterministic regression coverage.
created: "2026-08-11"
flags:
  - accepted
tags:
  - review
  - implementation
  - windows
keywords:
  - issue 47
  - Git Bash
  - shell auto
  - PR 55
owner: runx-implementation-reviews
---

# RunX Windows Git Bash Automatic Shell Review

## Verdict

Accepted with no findings at exact implementation head
`a5c9b1cf2f2f28131a2be74ff3c414d0e10e749a` in PR
[CGuiho/runx#55](https://github.com/CGuiho/runx/pull/55).

## Authorship Boundary

- Root planning baseline: `71868a6`.
- GPT-5.6 Luna at maximum reasoning effort authored executable code, tests,
  contract documentation, skill alignment, and implementation evidence in
  commits `10774b1`, `bcafba5`, and `a5c9b1c`.
- Root performed the independent review and did not edit implementation code.

## Review Results

- `BuildShellExecution` remains the public production entry and injects only a
  private platform, environment, and executable-lookup seam for deterministic
  tests.
- Windows `shell: auto` selects the resolved Bash executable only for
  recognized MSYS caller markers and rejects System32/Sysnative Bash or WSL
  launchers before falling back to `cmd.exe`.
- Explicit `cmd`, `powershell`, `bash`, and `sh` remain authoritative;
  non-Windows automatic selection remains `sh`.
- Configured command text, child arguments, selector behavior, confirmation,
  dry runs, reveal output, and exit propagation are unchanged.
- Regression tests cover the `/c/GUIHO` boundary, ordinary Windows fallback,
  missing and rejected Bash paths, explicit precedence, markers, and
  non-Windows behavior without accessing any secret-bearing file.

## Immutable Evidence

- Accepted review:
  [4904615114](https://github.com/CGuiho/runx/pull/55#pullrequestreview-4904615114).
- Exact-head validation:
  [5251208887](https://github.com/CGuiho/runx/pull/55#issuecomment-5251208887).
- CI: [31476447395](https://github.com/CGuiho/runx/actions/runs/31476447395),
  successful on Windows, Ubuntu, and the canonical release-contract job.
- Merge commit:
  [`a297aa25b94d83d11d361a0fab0b5415bdd1ba20`](https://github.com/CGuiho/runx/commit/a297aa25b94d83d11d361a0fab0b5415bdd1ba20).


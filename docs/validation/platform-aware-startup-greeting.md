---
name: Platform-Aware RunX Greeting Validation
purpose: Record verification evidence for GitHub issue 21.
description: Captures deterministic platform assertions, runtime smoke tests, full package checks, and remaining released-Linux verification.
created: 2026-07-20
flags:
  - release-ready
tags:
  - validation
  - cli
keywords:
  - issue 21
  - Windows
  - Linux
  - macOS
owner: runx-validation
---

# Platform-Aware RunX Greeting Validation

## Summary

The implementation is ready for patch release. The published Linux executable
is the final platform-specific closure gate.

## Commands Run

| Check | Result |
| --- | --- |
| Focused `bun test source/cli.spec.ts` banner assertions | Passed |
| Independently rerun background-maintenance timing test | Passed |
| `bun run typecheck` | Passed |
| `bun test` | Passed: 54 tests, 373 expectations |
| `bun run build` | Passed |
| `bun run binaries` | Passed |
| Compiled `runx-windows-x64-baseline.exe` no arguments | Passed: `Hello Windows - runx v0.5.1` |
| Compiled Windows `--version` | Passed: `0.5.1` |
| XDocs strict metadata, tree, and doctor | Passed |

## Deterministic Assertions

- `win32` renders `Hello Windows`.
- `linux` renders `Hello Linux`.
- `darwin` renders `Hello macOS`.
- Cached update output remains before the startup greeting.

## Remaining Gate

After publication, run the released Linux binary without arguments and record
`Hello Linux - runx v<released-version>` in issue 21 before closure.

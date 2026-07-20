---
name: Use The Runtime Platform In The RunX Greeting
purpose: Track the platform-aware no-argument greeting requested by GitHub issue 21.
description: Defines Windows, Linux, and macOS greeting labels, output ordering, regression coverage, and public verification.
created: 2026-07-20
flags:
  - testing
owner: runx-todo
tags:
  - cli
  - startup
keywords:
  - issue 21
  - operating system
  - greeting
---

# Use The Runtime Platform In The RunX Greeting

## Outcome

`runx` with no arguments names the operating system it is actually running on:
Windows, Linux, or macOS.

## External Tracker

- [CGuiho/runx#21](https://github.com/CGuiho/runx/issues/21)

## Acceptance Signals

- Windows prints `Hello Windows - runx v<version>`.
- Linux prints `Hello Linux - runx v<version>`.
- macOS prints `Hello macOS - runx v<version>`.
- Cached update notices remain before the greeting.
- Help, version, command routing, output channels, and exit codes are unchanged.
- Deterministic tests cover all three supported platform labels.

## Contract Note

This issue intentionally replaces the historical RFC 0034 literal-Windows
banner. The developer explicitly approved the platform-aware behavior.

---
name: No-Argument Welcome Window
purpose: Record the approved replacement for the legacy one-line startup greeting
description: Defines deterministic bare-invocation output and cached update notice placement without foreground network work.
created: 2026-07-22
flags:
  - decision
  - approved
tags:
  - decision
  - cli
keywords:
  - welcome
  - update cache
owner: runx-decisions
---

# No-Argument Welcome Window

## Decision

Bare RunX invocation uses a stable UTF-8 bordered welcome window with no ANSI
styling. It includes the RunX product identity, purpose, GUIHO organization,
platform, architecture, version, and help guidance.

A cached notice is rendered after the body only when validated SemVer proves
the cached latest stable version is newer than the running version. Reading the
cache and completing the local worker-spawn handoff may be awaited; the remote
request never is.

## Consequences

- This intentionally supersedes the historical exact `Hello <platform>` line.
- Help, version, Markdown, JSON, and hidden worker routes do not render it.
- Snapshot tests own width, ordering, final newline, and update placement.

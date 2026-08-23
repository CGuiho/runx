---
name: Guarantee Windows Installer Staging Cleanup
purpose: Define the completion contract for RunX TODO task 25.
description: Requires unconditional cleanup of the unique Windows installer staging directory after success or any terminating failure.
created: 2026-08-23
flags:
  - testing
tags:
  - installer
  - windows
  - cleanup
keywords:
  - .guiho/.temp
  - finally
  - rollback
owner: runx-todo
---

# Guarantee Windows Installer Staging Cleanup

## Status

- State: testing
- Updated: `2026-08-23T13:47:00+02:00`

## Incident

Failed pre-0.14.9 Windows installations left unique
`$HOME/.guiho/.temp/runx-install-*` directories because cleanup ran only on the
success path. Convention 0001 requires cleanup after both success and failure.

## Plan Unit

Wrap every post-staging installer operation in `try/finally`, remove only the
current operation's unique staging directory in `finally`, and assert cleanup
after clean install, injected activation failure, reinstall, and upgrade.

Existing Convention 0001 architecture and task 24 supply the requirements and
plan context; no separate architecture or plan is necessary.

## Acceptance Criteria

1. Success leaves no staging child for the operation.
2. Injected post-activation failure rolls back and leaves no staging child.
3. Cleanup remains confined to the operation's exact staging path.
4. Native Windows CI and release gates enforce the behavior.

## Validation

The native Windows integration passes script parsing, clean installation,
injected post-activation rollback, same-version reinstall, different-version
whole-release upgrade, same-version idempotence, and zero leftover
`runx-install-*` staging children after each installer outcome.

## Release Decision

Patch 0.14.9 required after complete CI and public failure-cleanup smoke.

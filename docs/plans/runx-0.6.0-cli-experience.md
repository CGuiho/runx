---
name: RunX 0.6.0 CLI Experience Plan
purpose: Execute the approved RunX portion of issues 23, 24, and 25
description: Sequences durable documentation, welcome rendering, safe child argument forwarding, installer copy, workflow simplification, validation, versioning, release, and issue closure.
created: 2026-07-22
flags:
  - approved
  - in-progress
tags:
  - plan
  - cli
keywords:
  - RunX 0.6.0
  - welcome
  - argument forwarding
owner: runx-plans
---

# RunX 0.6.0 CLI Experience Plan

## Units

1. Record issue-linked requirements, decisions, task specs, validation target,
   and XDocs ownership.
2. Add a pure welcome renderer and cached SemVer update model; route bare
   startup through it while preserving clean special outputs.
3. Add a TypeBox-validated `runx run` partition module and make Citty receive
   only the RunX-owned prefix.
4. Extend execution adapters to forward an immutable argument array safely for
   Bash, sh, PowerShell, cmd, and auto shell selection.
5. Extend text and JSON dry runs with lossless forwarded arguments.
6. Update README, canonical documentation, bundled skill, prompt, and XDocs.
7. Replace the public POSIX bootstrap with the exact simplified curl command
   while retaining installer-internal hardening.
8. Remove only the GitHub Environment approval gate from publishing; retain
   typecheck, tests, builds, release-note extraction, npm OIDC, and exact assets.
9. Run focused and complete tests, typecheck, builds, binaries, asset checks,
   banned-import scan, XDocs validation, and implementation review.
10. Use Mirror to plan and apply `0.6.0`, release with only its changelog
    section, verify fourteen assets and public installation, then comment on and
    close issues 23 through 25.

## Acceptance

- `runx` shows the approved welcome without a manifest.
- `runx run cli-ts -v` forwards `-v`.
- Hostile-looking argument values remain literal values.
- Dry runs are non-executing and machine-readable.
- Public install and native smoke checks pass.
- Issue 22 and all Go rewrite work remain untouched.

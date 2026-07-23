---
name: RunX Manifest V2 Composition Plan
purpose: Sequence issue 26 from contract through public release and closure.
description: Detailed implementation, migration, validation, release, and evidence plan for recursive commands and catalog composition.
created: 2026-07-23
flags:
  - executing
tags:
  - plans
  - cli
keywords:
  - manifest v2
  - issue 26
  - child catalogs
owner: runx-plans
---

# RunX Manifest V2 Composition Plan

1. Freeze issue 26 as a breaking v2 contract and explicitly exclude issue 22.
2. Replace `project.name` with required `namespace`; remove top-level `groups`.
3. Define TypeBox recursive command/group unions and reject extra legacy keys.
4. Validate sibling names, namespace collisions, group shape, global UIDs, and
   composed selectors before any command may run.
5. Resolve local references relative to the declaring catalog and normalize
   supported full GitHub URLs to bounded raw fetches.
6. Require reciprocal child `parent` declarations and reject mismatches,
   missing files, cycles, excessive depth, oversized payloads, and unsafe URLs.
7. Flatten validated catalogs while retaining local/foreign source, owning
   catalog path, namespace alias, canonical selector, and execution base.
8. Update list, describe, check, dry-run, run, and init to consume v2 without a
   second parser or compatibility branch.
9. Add focused tests for nested groups, renamed child mounts, direct child
   loading, cwd ownership, legacy rejection, collisions, reciprocity, and URL
   policy; retain every prior CLI and forwarding regression.
10. Update README, DOCS, architecture, skill, prompt, XDocs descriptors, TODO,
    changelog, and validation evidence with exact examples and migration rules.
11. Run typecheck, the complete suite, build, native matrix, exact fourteen
    assets, XDocs strict/doctor, and Mirror minor planning.
12. Use Mirror to release `0.7.0`, approve the bound npm production identity,
    verify npm provenance, exact GitHub assets, version-scoped notes, public
    installers, composed-catalog execution, then comment and close issue 26.

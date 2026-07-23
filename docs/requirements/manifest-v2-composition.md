---
name: RunX Manifest V2 Composition Requirements
purpose: Define the breaking colocated command tree and explicit parent-child catalog contract for issue 26.
description: Requires namespace identity, recursive groups, local and GitHub child catalogs, reciprocal parents, collisions, safe loading, selectors, and migration.
created: 2026-07-23
flags:
  - approved
tags:
  - requirements
  - cli
keywords:
  - manifest v2
  - namespace
  - child catalogs
owner: runx-requirements
---

# RunX Manifest V2 Composition Requirements

## Required Manifest

Manifest major version 2 removes `project` and top-level `groups`. Every file
has one identifier-safe `namespace`, one relative `scripts.directory`, an
optional explicit `parent`, and one recursive `commands` list.

A list entry is either a command with stable `uid` and local `id`, or a group
with `group`, `summary`, and exactly one of nested `commands` or `runx`.

```yaml
version: "2.0.0"
namespace: laboratory
scripts:
  directory: scripts
commands:
  - group: cli
    summary: CLI commands.
    commands:
      - uid: cli-test
        id: test
        summary: Test the CLI.
        description: Run the CLI test suite.
        command: bun test
```

## Composition

- A `runx` group mounts a child file under the group name. That name is the
  child namespace alias and may differ from the child's declared namespace.
- A child declares its reciprocal parent with a relative path or full GitHub
  URL. A mounted child without the matching parent declaration is invalid.
- Relative paths resolve from the declaring file. Foreign references accept
  only HTTPS GitHub blob or raw-content URLs and are marked `foreign`; local
  references are marked `local`.
- The loader follows only explicit references. It never searches ancestor
  directories and never silently merges unrelated catalogs.
- Foreign catalogs have bounded fetch time, bounded size, process-local
  lifetime, cycle protection, and no persistent stale cache.

## Identity And Execution

- Command IDs, group names, and mounted namespace aliases are unique among
  siblings. The manifest namespace cannot equal a top-level command or group.
- UIDs and full selectors are unique across the composed graph.
- Full selectors follow the recursive mount/group path, such as
  `worker-alias/build/compile`; a unique UID or command ID remains selectable.
- Local child commands resolve `cwd` from their own catalog directory. Foreign
  commands resolve from the local mount root because a URL has no local cwd.
- Cycles, depth beyond 32, non-GitHub URLs, absolute filesystem references,
  invalid parent reciprocity, and escaping command cwd fail with exit code 3.

## Breaking Migration

Manifest v1 is rejected. `runx init` creates v2. README, canonical docs,
architecture, bundled skill/prompt, tests, changelog, release assets, and public
validation must all teach the same v2 contract. The Go rewrite in issue 22 is
explicitly excluded.

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

`runx init` emits `.scripts` as the default scripts directory. A manifest that
sets `scripts.directory` explicitly retains that value.

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
  URL. Local-to-local and foreign-to-foreign mounts require exact canonical
  identity. When a local working copy mounts a foreign child, RunX loads the
  child's declared foreign parent and validates that complete upstream parent
  graph declares the child before accepting the working-copy mount.
- Relative paths resolve from the declaring file. Foreign references accept
  only HTTPS GitHub blob or raw-content URLs and are marked `foreign`; local
  references are marked `local`.
- The loader follows only explicit references. It never searches ancestor
  directories and never silently merges unrelated catalogs.
- Foreign catalogs have bounded fetch time, bounded size, process-local
  lifetime, cycle protection, and no persistent stale cache.
- Relative references inside a foreign graph cannot escape its GitHub owner,
  repository, and ref root; crossing those boundaries requires a full URL.

## Identity And Execution

- Command IDs, group names, and mounted namespace aliases are unique among
  siblings. The manifest namespace cannot equal a top-level command or group.
- UIDs remain globally unique and canonical selectors remain unique across the
  composed graph. Local command IDs are scoped by their containing group and
  become unqualified shorthands only when they have one owner. A UID may equal
  another command's ID; exact UID lookup wins before canonical selectors or ID
  shorthands. Duplicate unqualified IDs remain ambiguous and fail rather than
  selecting an arbitrary command.
- Full selectors follow the recursive mount/group path, such as
  `worker-alias/build/compile`; an exact UID or unambiguous command ID remains
  selectable.
- Local child commands resolve `cwd` from their own catalog directory. Foreign
  commands resolve from the local mount root because a URL has no local cwd.
- Cycles, depth beyond 32, non-GitHub URLs, absolute filesystem references,
  invalid parent reciprocity, and escaping command cwd fail with exit code 3.
- All catalog scripts directories and command cwd values are containment-
  validated during composition, so check/list fail before execution.

## Breaking Migration

Manifest v1 is rejected. `runx init` creates v2. README, canonical docs,
architecture, bundled skill/prompt, tests, changelog, release assets, and public
validation must all teach the same v2 contract. The Go rewrite in issue 22 is
explicitly excluded.

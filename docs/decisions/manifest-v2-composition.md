---
name: RunX Manifest V2 Composition Decision
purpose: Record the accepted breaking design for colocated groups and explicit catalog composition.
description: Replaces split group maps and project.name with recursive commands, namespace identity, reciprocal mounts, and bounded GitHub loading.
created: 2026-07-23
flags:
  - accepted
tags:
  - decisions
  - cli
keywords:
  - manifest v2
  - namespace
  - composition
owner: runx-decisions
---

# RunX Manifest V2 Composition Decision

## Decision

RunX uses manifest major version 2. `namespace` replaces `project.name` and
`commands` owns both executable leaves and recursively nested groups. A group
uses exactly one of `commands` or `runx`; a `runx` group explicitly mounts a
child and its group name is the namespace alias.

Children declare the parent reference that mounts them. RunX validates both
directions. Loading a child directly validates its parent but keeps the child
as the active catalog; loading the parent composes child commands below the
mount selector. This makes ownership explicit without implicit ancestor search.

The effective runtime catalog is flattened only after validation. Each command
retains its canonical selector, owning catalog, source kind, and local execution
base. Existing list, describe, dry-run, and run flows therefore consume one
validated model without duplicating graph logic.

## Foreign Catalog Boundary

Only HTTPS `github.com/.../blob/...` and `raw.githubusercontent.com` references
are accepted. Blob URLs normalize to raw URLs. Fetches time out after ten
seconds, reject payloads above one MiB, live only for the process, and share the
same cycle/depth/collision validation as local files. Foreign commands use the
local mount root for cwd resolution.

## Consequences

V1 manifests fail closed and must migrate. The previous public-group
requirement, split `groups` map, `group` field on command leaves, and
`project.name` are removed. Stable UIDs remain the preferred automation key.
This decision supersedes the manifest-shape portions of the alpha and
interactive-init decisions; it does not change execution confirmation, shell
safety, CLI technology, release assets, or issue 22.

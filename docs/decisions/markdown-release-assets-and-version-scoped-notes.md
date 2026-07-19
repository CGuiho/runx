---
name: RunX Markdown Release Assets And Version-Scoped Notes
purpose: Preserve the accepted RunX GitHub Release filename and release-description policy.
description: Records the project override requiring .md agent asset filenames and exact-version changelog sections for idempotent GitHub Release notes.
created: 2026-07-19
flags:
  - accepted
  - release-policy
tags:
  - decisions
  - releases
  - cli
keywords:
  - guiho-s-runx.md
  - guiho-i-runx.md
  - GitHub Release notes
  - changelog extraction
owner: runx-decisions
---

# RunX Markdown Release Assets And Version-Scoped Notes

## Status

Accepted by the developer on 2026-07-19.

## Context

The bundled RFC 0034 skill originally named the two Markdown release artifacts
without extensions. That makes their content type less obvious and conflicts
with the explicit RunX release contract chosen by the developer. The publish
workflow also passed the complete changelog to GitHub, causing one release to
describe unrelated versions.

## Decision

Every RunX GitHub Release contains exactly fourteen assets: twelve native
binaries plus `guiho-s-runx.md` and `guiho-i-runx.md`.

The `.md` suffixes are public filename requirements. Installers download those
names, install the skill content as `SKILL.md`, and use the instruction asset
content for managed instruction blocks.

Release notes contain only the exact `## <version> - <date>` section from
`CHANGELOG.md`, ending before the next `## ` heading. Publishing fails closed
when the exact section is missing or empty. Rerunning a tag updates the existing
release notes before replacing its assets.

Downloaded agent artifacts are treated as untrusted release inputs. Installers
reject empty payloads, Windows executable headers, NUL-containing binary data,
invalid UTF-8, missing YAML frontmatter, and mismatched resource identities
before writing any global skill or project instruction file.

## Alternatives Considered

- Extensionless agent assets: rejected because the developer explicitly
  requires recognizable Markdown filenames.
- Full-changelog release notes: rejected because they include unrelated release
  history.
- Creating notes only on first publish: rejected because reruns must repair
  stale release descriptions idempotently.
- Trusting the filename extension alone: rejected because a `.md` download can
  still contain an executable or corrupted payload.

## Consequences

- RunX intentionally differs from the original extensionless RFC asset spelling.
- All build, verification, workflow, installer, test, and public documentation
  surfaces must use the `.md` filenames.
- Every version tag requires an exact matching changelog heading before its
  protected publish workflow can succeed.
- Installer validation must complete before any downloaded agent resource is
  copied or applied.

## Reversal Or Revisit Conditions

Revisit only through an explicit developer decision that changes the public
release contract and migrates all consumers together.

## References

- [RFC migration plan](../plans/rfc-0034-cli-compliance-migration.md)
- [RFC task specification](../todo/rfc-0034-cli-compliance-migration.md)
- [RunX documentation](../../DOCS.md)

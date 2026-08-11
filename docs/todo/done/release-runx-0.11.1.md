---
name: Publish RunX 0.11.1
purpose: Preserve the completed patch-release scope, Mirror evidence, publication record, and production boundary for TODO task 20.
description: Records the successful Mirror-managed RunX 0.11.1 release for the confirmation-policy skill correction.
created: 2026-08-11
flags:
  - completed
tags:
  - todo
  - done
  - release
  - mirror
  - runx
keywords:
  - runx
  - 0.11.1
  - confirm
  - patch
  - mirror
owner: runx-todo-done
---

# Publish RunX 0.11.1

## Completion

- Task: `20. Publish RunX 0.11.1`
- Status: completed
- Release-preparation PR: [CGuiho/runx#49](https://github.com/CGuiho/runx/pull/49), merge commit [00789a4a8225bdae8fd4e78f418834af8d3662db](https://github.com/CGuiho/runx/commit/00789a4a8225bdae8fd4e78f418834af8d3662db).
- XDocs correction PR: [CGuiho/runx#51](https://github.com/CGuiho/runx/pull/51), main head before tagging [c84f6432d1e4a1169a356600265d0c583515371b](https://github.com/CGuiho/runx/commit/c84f6432d1e4a1169a356600265d0c583515371b).
- Final evidence: [RunX 0.11.1 release validation](../../validation/runx-0.11.1-release.md).

## Mirror and publication

The clean synchronized clone reported current `0.11.0` and exact patch plan
`0.11.1` with tag `@guiho/runx/v0.11.1`. Only the authorized command
`mirror version apply patch --yes` was run; it created annotated tag object
`7d17153b44185f4215079f2e792effee87f36e9f`, peeled to source commit
`c84f6432d1e4a1169a356600265d0c583515371b`.

The tag-triggered [Publish workflow 31469968479](https://github.com/CGuiho/runx/actions/runs/31469968479)
completed successfully after the explicitly authorized `production`
package/source environment approval. The public [GitHub Release](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.11.1)
is non-draft and non-prerelease with exactly 11 canonical assets. Independent
downloads matched `checksums.txt` and every GitHub SHA-256 digest; Windows
AMD64 version/help/help-tree smoke passed; and `npm.cmd view @guiho/runx@latest
version` returned `0.11.1`.

## Production boundary

The approved production environment was used only for source/package
publication. No deployment, promotion, traffic, DNS, database, secret, or
other production-infrastructure mutation occurred. There are no release
residuals.

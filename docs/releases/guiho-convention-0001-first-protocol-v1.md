---
name: GUIHO Convention 0001 First Protocol-v1 Release
purpose: Durable operation record for RunX release 0.13.0, the first protocol-v1 release.
description: Records the R00 release window for tag @guiho/runx/v0.13.0 — Mirror apply, protected tag, GitHub Release publication, workflow runs, asset set, and remote installer smoke evidence.
created: 2026-08-22
owner: runx
flags:
  - released
tags: [release, r00]
keywords: [0.13.0, protocol v1, mirror, github release]
---

# R00 — First Protocol-v1 Release (RunX 0.13.0)

## Binding

| Field | Value |
| --- | --- |
| Release | 0.13.0 |
| Protected tag | `@guiho/runx/v0.13.0` |
| Tag commit (`RELEASE_COMMIT_SHA`) | `5690a4d` |
| Tag object (remote) | `4668fd10f628051f72e7c6ebf5bc7d01b1e11efd` |
| Reachability | Tag points at exact HEAD of protected `main`; main pushed through `5690a4d`. |

## Authorization

The developer issued a durable execution order on 2026-08-22 ("start everything
right now! Do not stop until everything is done") followed by "Okay, you can
proceed with everything" after the R00 gate was presented explicitly. The
production environment deployment approvals were performed as CGuiho through
the authenticated session, once per run.

## Mirror

```text
mirror config check        -> ok (after migrating legacy mirror.yaml to mirror.config.toml)
mirror version plan 0.13.0 -> 0.12.2 -> 0.13.0; tag @guiho/runx/v0.13.0; push tags=true
mirror version apply 0.13.0 --yes -> applied
```

No version-bearing files exist outside git (source=git, output=git), so the
release commit is the tagged HEAD itself; only the tag was created and pushed.
CHANGELOG.md 0.13.0 section was written by the agent per
`agents.write_changelog` before the successful tag.

## Workflow runs (Publish)

| Run | Result | Note |
| --- | --- | --- |
| 32596100605 | failed | Missing CHANGELOG 0.13.0 section at tag; retagged. |
| 32597387362 | failed | `checksums.txt` did not cover `artifacts.json`; contract fixed in builder/verifier. |
| 32597652808 | failed | Release created as draft by `gh release create`, so remote asset URLs 404'd; published draft and hardened workflow with `--draft=false --latest`. |
| 32599277619 | failed | Installer succeeded end-to-end; final summary printf line-break bug mangled output (fixed), plus verification HOME fix in workflow. |
| 32599683000 | **success** | Full pipeline including remote installer smoke. |

Final URL: https://github.com/CGuiho/runx/actions/runs/32599683000

## Published release

URL: https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.13.0

- draft: false; published 2026-08-22T21:30:59Z
- 22 assets = 8 `runx-payload-*` + 8 `runx-launcher-*` + `guiho-s-runx.zip`
  + `guiho-i-runx.md` + `runx.schema.json` + `runx.global.schema.json`
  + `artifacts.json` + `checksums.txt`
- `checksums.txt` covers all 21 non-checksum artifacts (verified remotely).

## Remote installer smoke (CI, ubuntu runner)

Executed by the publish workflow against the live release:

1. Isolated `HOME`; `devops/install.sh --version 0.13.0`.
2. Downloaded payload, launcher, manifest, schemas, skill, instruction;
   SHA-256 verified against `checksums.txt` and `artifacts.json`.
3. Staged payload `--version` matched 0.13.0; hidden `__self-test` passed.
4. Installed immutable `versions/0.13.0/`, activated atomic pointer, installed
   both global skill copies, verified launcher `--version` = `0.13.0`, and ran
   the launcher `__self-test`.

## Production mutation

GitHub Release publication only. No servers, traffic, data, or secrets were
touched; the repository's production environment approval was the sole gate.

## Follow-ups recorded

- Legacy direct-binary installations upgrade through the retained legacy path
  until they reinstall via the new installers.
- H00 (transition-support removal) remains open and is now unblocked.

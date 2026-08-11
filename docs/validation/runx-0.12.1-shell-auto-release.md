---
name: RunX 0.12.1 Shell Auto Release Validation
purpose: Preserve final release evidence for runx reveal and the Windows Git Bash automatic-shell correction.
description: Records exact-head review, CI, protected merges, Mirror publications, 11-asset checks, npm state, native checksums, reveal output, and the non-secret Git Bash path probe.
created: "2026-08-11"
flags:
  - completed
tags:
  - validation
  - release
  - windows
keywords:
  - RunX 0.12.0
  - RunX 0.12.1
  - runx reveal
  - Git Bash
  - issue 47
owner: runx-validation
---

# RunX 0.12.1 Shell Auto Release Validation

## Outcome

RunX reveal and its motivating Windows shell-boundary correction are publicly
accepted. Issue [CGuiho/runx#47](https://github.com/CGuiho/runx/issues/47) is
closed as completed.

## Reveal Release

- Implementation PR:
  [#50](https://github.com/CGuiho/runx/pull/50), merged as
  [`cccc87348da130df90665e954f48bad9bd652ceb`](https://github.com/CGuiho/runx/commit/cccc87348da130df90665e954f48bad9bd652ceb).
- Repository evidence PR:
  [#53](https://github.com/CGuiho/runx/pull/53), merged as
  [`7baa214c5514c39fdc5af8118925ffde015be5ae`](https://github.com/CGuiho/runx/commit/7baa214c5514c39fdc5af8118925ffde015be5ae).
- Release-note recovery PR:
  [#54](https://github.com/CGuiho/runx/pull/54), merged as
  [`107f3f69a71c05a877a9854fee063ddce2ecd0ce`](https://github.com/CGuiho/runx/commit/107f3f69a71c05a877a9854fee063ddce2ecd0ce).
- Final release:
  [`@guiho/runx/v0.12.0`](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.12.0),
  annotated tag object `2403da6081b8c138fd28defb34ba8569b180e107`,
  peeled source `107f3f69a71c05a877a9854fee063ddce2ecd0ce`.
- Publish workflow:
  [31474091754](https://github.com/CGuiho/runx/actions/runs/31474091754),
  successful.

The first 0.12.0 tag run stopped before publication because the changelog note
was still under `Unreleased`. Public checks confirmed that no GitHub Release or
npm 0.12.0 existed before the unpublished failed tag was removed and Mirror
recreated it at corrected main.

## Windows Shell Correction

- Exact implementation head:
  `a5c9b1cf2f2f28131a2be74ff3c414d0e10e749a`.
- Review: [runx-windows-git-bash-auto-shell-review.md](../reviews/implementation/runx-windows-git-bash-auto-shell-review.md).
- Implementation PR:
  [#55](https://github.com/CGuiho/runx/pull/55), merged as
  [`a297aa25b94d83d11d361a0fab0b5415bdd1ba20`](https://github.com/CGuiho/runx/commit/a297aa25b94d83d11d361a0fab0b5415bdd1ba20).
- Release-note PR:
  [#56](https://github.com/CGuiho/runx/pull/56), merged as release source
  [`2c96e8aec1052b35c6fb4061eba5a6b54405bbcb`](https://github.com/CGuiho/runx/commit/2c96e8aec1052b35c6fb4061eba5a6b54405bbcb).
- Final release:
  [`@guiho/runx/v0.12.1`](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.12.1),
  annotated tag object `6b37f30c2935e3e2182ada3886c25fc71d4cb4b9`,
  peeled source `2c96e8aec1052b35c6fb4061eba5a6b54405bbcb`.
- Publish workflow:
  [31476981185](https://github.com/CGuiho/runx/actions/runs/31476981185),
  successful.

## Validation

- Focused/full Go tests, gofmt, tidy cleanliness, vet, and build passed.
- Strict XDocs metadata, tree, and doctor passed with zero errors/warnings.
- CI run 31476447395 passed Windows, Ubuntu, and canonical release-contract
  jobs for the exact reviewed implementation head.
- Both protected publication workflows built eight native binaries and three
  supporting assets; GitHub exposes exactly 11 assets per release.
- The 0.12.1 publisher passed version-scoped notes, exact asset verification,
  exact-version installer smoke, and npm OIDC publication.
- Public npm `latest` is `0.12.1`.
- Downloaded 0.12.1 Windows AMD64 SHA-256
  `694aff5bfa71334f54e297eff8dbfc2fad67cb4b1a916300a19778c0b01a0cdb`
  matched `checksums.txt`; the binary reported `0.12.1`.
- Released `runx reveal dev` against the GUIHO Core catalog printed the exact
  stored Dotenvx command without executing it.
- Non-secret native Node path probe from Git Bash:
  - published 0.12.0 automatic shell: `/c/GUIHO`;
  - reviewed and released 0.12.1 automatic shell: `C:/GUIHO`;
  - reveal output remained the original unmodified command.

No environment file, key file, or secret value was accessed. The protected
GitHub `production` environment was used only to authorize source/package
publication; no application deployment, traffic, DNS, database, or secret
mutation occurred.


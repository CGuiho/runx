---
name: RunX 0.11.1 Release Validation
purpose: Preserve independent publication, artifact, checksum, native smoke, and package evidence for RunX 0.11.1.
description: Records the final clean-clone validation and successful Mirror-managed 0.11.1 publication for the confirmation-policy skill correction.
created: 2026-08-11
flags:
  - accepted
  - published
tags:
  - validation
  - release
  - mirror
  - runx
  - 0.11.1
keywords:
  - runx
  - 0.11.1
  - confirm
  - opt-in
  - checksums
  - npm
owner: runx-validation
---

# RunX 0.11.1 Release Validation

## Verdict

RunX `0.11.1` is published and independently verified. The Mirror-managed
annotated tag `@guiho/runx/v0.11.1` peels to source commit
[c84f6432d1e4a1169a356600265d0c583515371b](https://github.com/CGuiho/runx/commit/c84f6432d1e4a1169a356600265d0c583515371b),
the final release workflow succeeded, the GitHub Release is public and
non-draft/non-prerelease, all canonical assets are present, and npm latest
resolves to `0.11.1`.

## Integration and release-preparation evidence

- Confirmation-policy implementation PR: [CGuiho/runx#46](https://github.com/CGuiho/runx/pull/46), accepted review [4903741635](https://github.com/CGuiho/runx/pull/46#pullrequestreview-4903741635), READY validation [5250161313](https://github.com/CGuiho/runx/pull/46#issuecomment-5250161313), merge commit [37327d5e80d405ecee6d8fa29bb9dd3bdc1302b8](https://github.com/CGuiho/runx/commit/37327d5e80d405ecee6d8fa29bb9dd3bdc1302b8).
- Protected integration-evidence PR [#48](https://github.com/CGuiho/runx/pull/48) merged at [93dc2f201045003b64360e2b214a6f7c86c176e6](https://github.com/CGuiho/runx/commit/93dc2f201045003b64360e2b214a6f7c86c176e6).
- Protected release-preparation PR [#49](https://github.com/CGuiho/runx/pull/49) merged at [00789a4a8225bdae8fd4e78f418834af8d3662db](https://github.com/CGuiho/runx/commit/00789a4a8225bdae8fd4e78f418834af8d3662db), finalizing the dated `0.11.1` changelog section and empty `Unreleased` section.
- Protected XDocs correction PR [#51](https://github.com/CGuiho/runx/pull/51) merged at [c84f6432d1e4a1169a356600265d0c583515371b](https://github.com/CGuiho/runx/commit/c84f6432d1e4a1169a356600265d0c583515371b).

## Clean-clone gates before tagging

The synchronized clean clone was detached at `c84f6432d1e4a1169a356600265d0c583515371b` and had no dirty files. The following passed before applying the version:

- `gofmt`, `go mod tidy` with a clean diff, `go test -count=1 ./...`, `go vet ./...`, and `go build ./...`.
- `xdocs scan`, `xdocs tree`, strict `xdocs meta`, and `xdocs doctor --warnings-as-errors` (0 errors, 0 warnings).
- `mirror` bootstrap, `mirror config check`, `mirror version current` (`0.11.0`), and `mirror version plan patch` (`0.11.1`, tag `@guiho/runx/v0.11.1`).
- Native release build and verifier produced exactly 11 assets; downloaded Windows AMD64 smoke later repeated `--version`, `--help`, and `--help-tree` successfully.

Only the authorized command `mirror version apply patch --yes` was run. It
created and pushed annotated tag object
`7d17153b44185f4215079f2e792effee87f36e9f`, which peels to the source commit
above. No manual tag was created.

## Publication evidence

- GitHub Release: [@guiho/runx/v0.11.1](https://github.com/CGuiho/runx/releases/tag/%40guiho/runx/v0.11.1), release id `368402599`, `draft=false`, `prerelease=false`.
- Publish workflow: [run 31469968479](https://github.com/CGuiho/runx/actions/runs/31469968479), head `c84f6432d1e4a1169a356600265d0c583515371b`, completed `success`; publish job `93710916755` completed successfully after the authorized `production` package/source environment approval.
- npm: `npm.cmd view @guiho/runx@latest version` returned `0.11.1` with exit code 0.

The exact 11 authored assets are:

`checksums.txt`, `guiho-i-runx.md`, `guiho-s-runx.zip`, `runx-darwin-amd64`,
`runx-darwin-arm64`, `runx-linux-amd64`, `runx-linux-arm64`,
`runx-linux-armv6`, `runx-linux-armv7`, `runx-windows-amd64.exe`, and
`runx-windows-arm64.exe`.

The downloaded `checksums.txt` values matched independent SHA-256 hashes for
all ten payloads. Every downloaded asset, including `checksums.txt`, also
matched its GitHub API `digest`:

| Asset | SHA-256 |
| --- | --- |
| `guiho-i-runx.md` | `fd3cd1dd9e853686f5451ccd8a00eeb72368e8fbc6aa800c9d08e4f91d9f62b2` |
| `guiho-s-runx.zip` | `3c2f7173bdddc00f3fa767b4399ac5ceffc7d3db15e93e5f693f00ab67d5e786` |
| `runx-darwin-amd64` | `baa5de8ca5aea2c206a4139b49a1cf4e7c91c410d2346e592ba56150039f1ee7` |
| `runx-darwin-arm64` | `14be684530be9cf1731fb22a672599cd3a4dfffe9401e86e62023d1809f5afb1` |
| `runx-linux-amd64` | `7d4f2d91df814d743b0c8b14747db59da36171487403da991af3eba27a572fe6` |
| `runx-linux-arm64` | `0d895bf4ddf81f00d8ec734a60b5b6f406dd7ae107ea87c634228f76f355907e` |
| `runx-linux-armv6` | `364cf2a0a3ca17f81b2776b51ae8936cc9335b3fc4ce7d396a24647a73bdb7c7` |
| `runx-linux-armv7` | `60623d11f2eef076f1e3f73a3c2b161807421beda9be356ba0886cca99f604ed` |
| `runx-windows-amd64.exe` | `b888366c4f2305c2ef89c4e936761b0abfb31e7355d28165b16411a338fb63f6` |
| `runx-windows-arm64.exe` | `37e5588c517fa445de1f35bc3f2e983a0e5dafa6cb18823420beb07aeaaa4f08` |
| `checksums.txt` | `0897093d88e3de94dbf49049db6c121f7881219d3519a95fa240bd3ba7df85d5` |

Downloaded native smoke:

- `runx-windows-amd64.exe --version` printed `0.11.1` and exited 0.
- `runx-windows-amd64.exe --help` exited 0.
- `runx-windows-amd64.exe --help-tree` exited 0.

## Production boundary and residuals

The approved `production` environment was used only for the repository's
tag-triggered source/package publication workflow. No deployment, promotion,
traffic, DNS, database, secret, or other production-infrastructure mutation
occurred. Release tasks 19 and 20 can now be archived as completed; no release
residual remains.

# Handoff — RunX Smart release plumbing (Unit 2)

Mode: dnd

## Task identity

- GitHub issue: https://github.com/CGuiho/runx/issues/66
- Task: `30. RunX Smart Command and rxs Launcher`
- Project: GUIHO #2, Component `runx`, Status `In Progress` (already set)
- Spec: `docs/todo/2026-10-22-runx-smart-and-rxs-launcher.md`

## Repository and required reading

- Repository: `/root/swe/runx` (Go + Cobra, branch `main`).
- Read `AGENTS.md`, then
  `/root/swe/guiho/conventions/guiho-convention-0001-cli.md` (complete release
  artifact set, installer sequence, uninstaller contract, shared `.guiho`
  ownership).
- Load skills: `guiho-s-0001-guiho`, `guiho-s-0004-working-with-cg`,
  `guiho-s-0032-git-commit`, `guiho-s-xdocs`.

## Model decision

- Provider/model: `xiaomi/mimo-v2.6-pro`, provider-default thinking, no variant.

## Exclusive paths

You own only:

- `devops/build-binaries.go`
- `devops/verify-release-assets.go`
- `devops/install.sh`
- `devops/install.ps1`
- `devops/uninstall.sh`
- `devops/uninstall.ps1`
- `devops/install_ps1_test.go`
- `devops/devops.xdocs.md`
- `cmd/uninstall.go`
- `.github/workflows/publish.yml`

Do not touch any other `cmd/` file, `docs/`, `README.md`, `DOCS.md`, `skills/`,
`embed/`, `prompts/`, `runx.yaml`, or `devops/extract-release-notes.go`.
A sibling worker creates `cmd/rxs/` and `cmd/smart.go`; never edit or create
those.

## Artifact contract (fixed — do not invent variations)

- Eight new assets: `rxs-launcher-<platform>` with the existing platform names
  (`linux-amd64`, `linux-arm64`, `linux-armv7`, `linux-armv6`, `darwin-amd64`,
  `darwin-arm64`, `windows-amd64.exe`, `windows-arm64.exe`).
- Built from `./cmd/rxs` with the existing cross-build helper and build stamp.
- `artifacts.json` entry per platform:
  - `id`: `rxs-launcher/<platform>`
  - `file`: `rxs-launcher-<platform>[.exe]`
  - `kind`: `launcher`
  - `installedPath`: `./.guiho/bin/rxs[.exe]`
  - `os`/`arch`/`version` as the other launchers.
- No `rx` artifact and no compatibility alias for `rxs`.
- Resulting counts: 46 declared artifacts, 47 checksum entries. Update the
  hardcoded counts in `verify-release-assets.go` and any other hardcoded counts
  (including generated-summary text).
- `install.sh` and `install.ps1`: download and verify the `rxs` launcher for the
  platform, install it next to `runx` under `$HOME/.guiho/bin`, back up and
  restore it on rollback together with the `runx` launcher, verify
  `rxs --version` equals the resolved version after activation, and print the
  installed `rxs` path.
- `uninstall.sh`, `uninstall.ps1`, and `runx uninstall` remove the `rxs`
  launcher as a CLI-owned artifact while preserving the shared
  `$HOME/.guiho/bin` directory.
- `.github/workflows/publish.yml`: after the exact-version install smoke, also
  require the installed `rxs --version` to equal `$VERSION` (and keep the
  existing `runx` checks). Do not add a bundle-directory `rxs` smoke, because
  the launcher needs an activated pointer.
- `devops/install_ps1_test.go`: the fake release and upgrade scenarios must
  declare and serve the `rxs` launcher; whole-release upgrade assertions must
  cover `rxs`.

## Dependency note

`go run devops/build-binaries.go` and `verify-release-assets.go` can only pass
after the sibling core unit lands `cmd/rxs`. Write the plumbing now, run
`bash -n` on both shell scripts, `gofmt` the Go files, and run the devops tests
that do not need the real matrix. If a build fails only because `./cmd/rxs` is
missing, record that exact failure and stop; the parent validates the matrix
after integration. Do not create `cmd/rxs` yourself.

## Validation you must run and report

- `gofmt -l devops` clean
- `bash -n devops/install.sh devops/uninstall.sh`
- `go vet ./...` (repository-wide, read-only)
- `go test -count=1 ./devops/...` for the tests that can run on Linux
- `go run devops/install_ps1_test.go` must not be run; the PowerShell test runs
  on Windows only.
- Any exact failure caused by the missing `cmd/rxs` unit.

## Commit authority

- Commit each coherent unit with `guiho-s-0032-git-commit` (explicit paths).
- Never push. The parent reviews and pushes.
- If a Git index lock collides with a sibling worker, wait briefly and retry.
- Report commit hashes, commands run, and residual risk.

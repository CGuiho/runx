---
name: GUIHO Convention 0001 First Protocol-v1 Release Validation
purpose: Record validation evidence for the RunX 0.13.0 protocol-v1 release.
description: Local gate evidence, release-contract verification, and remote smoke results bound to the released commit.
created: 2026-08-22
owner: runx-validation
flags:
  - passed
tags: [validation, r00]
keywords: [0.13.0, validation, protocol v1]
---

# Validation — RunX 0.13.0 (First Protocol-v1 Release)

## Bound head

`RELEASE_COMMIT_SHA` = `5690a4d` (tag `@guiho/runx/v0.13.0`, remote-verified).

## Local gates (all green at the tagged commit)

| Gate | Result |
| --- | --- |
| `gofmt -l main.go cmd pkg devops` | clean |
| `go mod tidy -diff` | clean |
| `go vet ./...` | clean |
| `go build ./...` | ok |
| `go test -count=1 ./...` | 10/10 packages ok |
| `runx check --format json` | valid, namespace runx, 6 commands |
| `bash -n devops/install.sh devops/uninstall.sh` | ok |
| PowerShell scriptblock parse of both .ps1 scripts | ok |
| `git diff --check` | clean |

## Release contract

```text
go run devops/build-binaries.go --version 0.13.0 --commit <sha> --build-date <rfc3339>
go run devops/verify-release-assets.go
```

Result: 20 declared artifacts + manifest + checksums; manifest digests match
files; checksums.txt covers all 21 non-checksum artifacts; skill archive
contains `guiho-s-runx/SKILL.md`.

## Remote evidence (CI publish run 32599683000 — success)

1. Built and verified the matrix from the protected tag.
2. Extracted version-scoped notes from CHANGELOG.md 0.13.0.
3. Published the GitHub Release with exactly the 22-file asset set
   (`expected-assets.txt` vs actual assets diff passed).
4. Remote exact-version installer smoke in an isolated home: download,
   checksum + manifest verification, staged payload `--version` match,
   hidden self-test pass, immutable install, atomic activation, global skill
   copies, launcher `--version` = `0.13.0`, launcher self-test pass.

## Residual risk

- Windows installer path was validated by syntax checks and local structural
  tests only; no native Windows runner smoke ran for this release window
  (labelled build-only for platform coverage).
- Legacy direct-binary installations continue on the legacy upgrade path until
  they reinstall through the new installers (documented behavior).

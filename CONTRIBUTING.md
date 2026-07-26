---
name: Contributing to RunX
purpose: Explain how to change and validate the native RunX CLI safely.
description: Go validation, documentation, agent-resource, and protected release expectations.
created: 2026-07-12
owner: runx
flags: []
tags:
  - documentation
  - contributing
keywords:
  - runx
  - go
  - contributing
---

# Contributing to RunX

RunX production behavior is implemented in Go. Keep changes focused and
preserve the distinction between catalog inspection and command execution.

1. Use `gofmt` on changed Go files.
2. Run `go mod tidy`, `go test ./...`, and `go vet ./...`.
3. For release-sensitive work, build and verify the standard 11 artifacts with
   the Go devops programs.
4. Update `DOCS.md`, the bundled `guiho-s-runx` skill, and XDocs descriptors
   whenever CLI behavior, manifest fields, installers, or agent workflows change.
5. Preserve the protected tag-driven release workflow. GitHub Release assets
   and the optional npm native bootstrap publish only from the `production`
   environment.
6. Use Mirror for every version plan or application. Do not edit versions or
   create release tags manually.

The retained `source/` and TypeScript devops files document the former CLI and
must not become production, CI, or release authority again.

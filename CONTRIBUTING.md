# Contributing to RunX

RunX is a Bun/TypeScript CLI. Keep changes focused and preserve the distinction
between catalog inspection and command execution.

1. Install dependencies with `bun install`.
2. Run `bun run typecheck` and `bun test` before opening a pull request.
3. Update `DOCS.md`, the bundled `guiho-s-runx` skill, and XDocs descriptors
   whenever CLI behavior, manifest fields, installers, or agent workflows change.
4. Do not add npm publishing or release publishing automation without an
   explicit project decision.
5. Use Mirror for version planning and version application; do not edit a
   release version manually.

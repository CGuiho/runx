# Security Policy

## Reporting

Report a potential security issue privately to `cg@guiho.co`. Do not include
credentials, tokens, or private manifests in public issues.

## RunX Trust Boundary

RunX manifests contain shell commands and are trusted local executable code.
RunX validates structure and makes commands visible; it cannot make arbitrary
shell code safe. Do not store secrets in `runx.yaml`, inspect a command before
running it, and use `confirm: always` for operations that benefit from an
explicit `--yes` acknowledgement.

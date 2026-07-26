---
subject: runx-manifest
description: Strict manifest-v2 YAML decoding, semantic validation, exact configuration resolution, and bounded local or GitHub catalog composition.
parent: runx-packages
children: []
files:
  types.go: Defines typed source manifests and resolved catalog, command, group, child, and provenance results.
  parser.go: Strictly decodes one YAML document and validates version, identity, paths, shells, confirmation, and entry shapes.
  composition.go: Resolves configuration precedence, composes reciprocal child catalogs, bounds remote reads and graphs, and indexes identities.
  manifest_test.go: Covers strict decoding, semantic failures, precedence without ancestor search, and reciprocal local composition.
documents: {}
tags:
  - go
  - yaml
  - manifest
keywords:
  - runx.yaml
  - KnownFields
  - composition
flags: []
status: stable
---

All catalog inspection and execution consumes the same validated `Catalog`.


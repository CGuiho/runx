---
subject: runx-devops
description: Native release-asset compilation and direct RunX installer scripts.
parent: runx
children: []
files:
  build-binaries.ts: Cross-compiles supported Bun native executable assets.
  install.ps1: Installs the latest Windows executable from GitHub Releases.
  install.sh: Installs the latest macOS or Linux executable from GitHub Releases.
documents: {}
tags:
  - devops
  - installers
  - releases
keywords:
  - runx
  - windows
  - macos
  - linux
  - native binary
flags: []
status: stable
---

These scripts support direct native installation without an npm publication.

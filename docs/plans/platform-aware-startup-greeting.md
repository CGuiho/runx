---
owner: runx-plans
tags:
  - cli
  - startup
keywords:
  - issue 21
  - platform label
  - regression test
---

# Platform-Aware RunX Greeting Plan

## Unit 1: Centralize Greeting Rendering

1. Add a pure greeting renderer to the existing Citty CLI module.
2. Map `win32`, `linux`, and `darwin` to Windows, Linux, and macOS.
3. Route the no-argument command through the renderer without changing cached
   update notice ordering.

## Unit 2: Prove Every Supported Platform

1. Add deterministic renderer assertions for Windows, Linux, and macOS.
2. Update child-process expectations to use the runtime platform.
3. Run the focused CLI suite and retry only independently diagnosed timing
   flakes.
4. Run the full RunX validation gates.

## Unit 3: Document, Release, And Close

1. Replace literal-Windows startup guidance in repository instructions and CLI
   documentation.
2. Update the owning XDocs descriptors.
3. Include the correction in the same Mirror patch as issue 20.
4. Verify the released Linux binary prints `Hello Linux`, post evidence to
   issue 21, and close it.

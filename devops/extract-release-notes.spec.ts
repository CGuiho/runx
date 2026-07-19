/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { describe, expect, test } from 'bun:test'
import { extractReleaseNotes } from './extract-release-notes.js'

describe('extractReleaseNotes', () => {
  test('extracts only the exact unbracketed version section', () => {
    const changelog = `---
name: Example
---

# Changelog

## 0.4.10 - 2026-07-20

### Fixed

- Newer patch.

## 0.4.1 - 2026-07-18

### Fixed

- Requested patch.

## 0.4.0 - 2026-07-18

### Added

- Older release.
`

    expect(extractReleaseNotes(changelog, '0.4.1')).toBe(`## 0.4.1 - 2026-07-18

### Fixed

- Requested patch.
`)
  })

  test('fails closed when the exact version heading is absent', () => {
    expect(() => extractReleaseNotes('## 0.4.1 - 2026-07-18\n\n- Fixed.\n', '0.4.2')).toThrow(
      'Expected exactly one changelog heading for version 0.4.2.',
    )
  })

  test('rejects an empty version section', () => {
    expect(() => extractReleaseNotes('## 0.4.2 - 2026-07-19\n## 0.4.1 - 2026-07-18\n\n- Older.\n', '0.4.2')).toThrow(
      'Changelog section for version 0.4.2 has no release notes.',
    )
  })
})

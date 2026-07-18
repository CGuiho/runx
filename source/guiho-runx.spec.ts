/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { describe, expect, test } from 'bun:test'
import { expectedReleaseAssetNames } from '../devops/build-binaries.js'

describe('RunX RFC structural contracts', () => {
  test('core source contains no prohibited Node built-ins', async () => {
    const source = Bun.fileURLToPath(new URL('.', import.meta.url))
    const files = [...new Bun.Glob('*.ts').scanSync({ cwd: source, onlyFiles: true })].filter((file) => !file.endsWith('.spec.ts'))
    const prohibited = /node:(?:fs(?:\/promises)?|child_process|path|os)/
    for (const file of files) expect(await Bun.file(`${source}/${file}`).text()).not.toMatch(prohibited)
  })

  test('release matrix contains exactly twelve binaries and two agent assets', () => {
    expect(expectedReleaseAssetNames).toHaveLength(14)
    expect(new Set(expectedReleaseAssetNames).size).toBe(14)
    expect(expectedReleaseAssetNames.filter((name) => name.startsWith('runx-'))).toHaveLength(12)
    expect(expectedReleaseAssetNames).toContain('guiho-s-runx')
    expect(expectedReleaseAssetNames).toContain('guiho-i-runx')
    expect(expectedReleaseAssetNames.some((name) => name.includes('macos'))).toBe(false)
  })

  test('npm bootstrap is the only Node runtime exception', async () => {
    const wrapper = await Bun.file(new URL('../scripts/runx-bin.mjs', import.meta.url)).text()
    expect(wrapper).toContain("from 'node:child_process'")
    expect(wrapper).toContain('runx-${os}-x64-baseline')
    expect(wrapper).not.toContain('source/guiho-runx-bin')
    expect(wrapper).not.toContain('library/guiho-runx-bin')
  })
})

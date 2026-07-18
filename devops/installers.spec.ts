/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { describe, expect, test } from 'bun:test'

describe('RunX direct installers', () => {
  test('PowerShell installer exposes the complete RFC sequence', async () => {
    const script = await Bun.file(new URL('./install.ps1', import.meta.url)).text()
    expect(script).toContain('Initiating GUIHO CLI Upgrade / Installation Sequence...')
    expect(script).toContain('guiho-s-runx')
    expect(script).toContain('guiho-i-runx')
    expect(script).toContain('.agents\\skills\\guiho-s-runx')
    expect(script).toContain('.claude\\skills\\guiho-s-runx')
    expect(script).toContain('BEGIN RUNX — DO NOT EDIT THIS SECTION')
    expect(script).toContain('Test-NativeBinary')
    expect(script).toContain('Install-Transactional')
  })

  test('POSIX installer selects Darwin assets and has no Bun dependency', async () => {
    const script = await Bun.file(new URL('./install.sh', import.meta.url)).text()
    expect(script).toContain("Darwin) printf 'darwin")
    expect(script).toContain('--progress-bar')
    expect(script).toContain('guiho-s-runx')
    expect(script).toContain('guiho-i-runx')
    expect(script).toContain('.agents/skills/guiho-s-runx')
    expect(script).toContain('.claude/skills/guiho-s-runx')
    expect(script).toContain('BEGIN RUNX — DO NOT EDIT THIS SECTION')
    expect(script).not.toContain(' bun ')
  })
})

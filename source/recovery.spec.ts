import { describe, expect, test } from 'bun:test'
import { createRecoveryInstructions } from './recovery.js'

describe('recovery instructions', () => {
  test('pins Windows installation and keeps process termination separate', () => {
    const recovery = createRecoveryInstructions('1.2.0-alpha.3', 'windows')
    expect(recovery.installCommand).toContain('-Version "1.2.0-alpha.3"')
    expect(recovery.installCommand).toContain('install.ps1')
    expect(recovery.targetSource).toBe('resolved')
    expect(recovery.stopProcessCommand).toContain('Get-Process runx')
    expect(recovery.stopProcessCommand).not.toBe(recovery.installCommand)
  })

  test('pins POSIX installation and exact process matching', () => {
    const recovery = createRecoveryInstructions('1.2.3', 'linux')
    expect(recovery.installCommand).toContain("--version '1.2.3'")
    expect(recovery.installCommand).toContain('| bash -s --')
    expect(recovery.stopProcessCommand).toBe('pkill -x runx')
  })
})

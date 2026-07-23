/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { describe, expect, test } from 'bun:test'
import { renderWelcome } from './welcome.js'

describe('RunX welcome', () => {
  test('renders deterministic product, platform, architecture, version, and help guidance', () => {
    const output = renderWelcome({ platform: 'win32', architecture: 'x64', version: '1.2.3' })

    expect(output).toContain('RUNX')
    expect(output).toContain('Documented command catalog')
    expect(output).toContain('GUIHO · Cristóvão GUIHO')
    expect(output).toContain('platform      Windows x64')
    expect(output).toContain('version       v1.2.3')
    expect(output).toContain('Run `runx --help` to see available commands.')
    expect(output.endsWith('\n')).toBe(true)
    expect(output).not.toContain('\u001B[')
  })

  test('uses canonical platform labels and appends a cached notice after the body', () => {
    expect(renderWelcome({ platform: 'linux', architecture: 'arm64', version: '1.0.0' })).toContain('Linux arm64')
    expect(renderWelcome({ platform: 'darwin', architecture: 'x64', version: '1.0.0' })).toContain('macOS x64')

    const output = renderWelcome({
      platform: 'linux',
      architecture: 'x64',
      version: '1.0.0',
      updateNotice: '  ⚠ New version available: v1.1.0\n    Run `runx upgrade` to update.',
    })
    expect(output.indexOf('Run `runx --help`')).toBeLessThan(output.indexOf('⚠ New version available'))
    expect(output.match(/New version available/g)).toHaveLength(1)
  })
})

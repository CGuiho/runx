import { describe, expect, test } from 'bun:test'
import { renderReleaseCatalog, renderUpgradeEvent, renderUpgradeHeading, renderUpgradePlan, renderUpgradeResult } from './upgrade-reporting.js'
import type { ReleaseCatalog, UpgradeEnvelope, UpgradePlan } from './upgrade-types.js'

const plan: UpgradePlan = {
  currentVersion: '0.2.7',
  targetVersion: '0.2.8',
  os: 'windows',
  arch: 'x64',
  assetName: 'runx-windows-x64-baseline.exe',
  assetUrl: 'https://example.test/runx.exe',
  executablePath: 'C:\\Users\\crist\\.local\\bin\\runx.exe',
}

describe('upgrade reporting', () => {
  test('renders the plan and phases in the required order', () => {
    const output = [
      renderUpgradeHeading(),
      renderUpgradePlan(plan),
      renderUpgradeEvent({ sequence: 3, phase: 'download', status: 'started' }),
      renderUpgradeEvent({ sequence: 5, phase: 'validate', status: 'started' }),
      renderUpgradeEvent({ sequence: 7, phase: 'replace', status: 'started' }),
      renderUpgradeEvent({ sequence: 9, phase: 'verify', status: 'started' }),
    ].join('')
    expect(output.indexOf('target  : 0.2.8')).toBeLessThan(output.indexOf('Downloading...'))
    expect(output.indexOf('Downloading...')).toBeLessThan(output.indexOf('Validating...'))
    expect(output.indexOf('Validating...')).toBeLessThan(output.indexOf('Replacing...'))
    expect(output.indexOf('Replacing...')).toBeLessThan(output.indexOf('Verifying...'))
  })

  test('always renders pinned recovery and a separate stop command', () => {
    const result: UpgradeEnvelope = {
      schemaVersion: 1,
      command: 'runx upgrade',
      outcome: 'upgraded',
      plan,
      events: [],
      result: { installedVersion: '0.2.8', cleanupDeferred: false },
      recovery: { targetVersion: '0.2.8', targetSource: 'resolved', installCommand: 'install 0.2.8', stopProcessCommand: 'stop runx' },
      error: null,
    }
    const output = renderUpgradeResult(result)
    expect(output).toContain('Upgrade complete: 0.2.7 -> 0.2.8')
    expect(output).toContain('install 0.2.8')
    expect(output).toContain('stop runx')
    expect(Object.keys(result)).toEqual(['schemaVersion', 'command', 'outcome', 'plan', 'events', 'result', 'recovery', 'error'])
    expect(Object.keys(result.recovery)).toEqual(['targetVersion', 'targetSource', 'installCommand', 'stopProcessCommand'])
  })

  test('visibly labels fallback-current recovery when discovery has no plan', () => {
    const result: UpgradeEnvelope = {
      schemaVersion: 1,
      command: 'runx upgrade',
      outcome: 'failed',
      plan: null,
      events: [{ sequence: 1, phase: 'plan', status: 'failed' }],
      result: null,
      recovery: { targetVersion: '0.2.7', targetSource: 'fallback-current', installCommand: 'repair 0.2.7', stopProcessCommand: 'stop runx' },
      error: { code: 'release_lookup_failed', phase: 'plan', message: 'offline' },
    }
    expect(renderUpgradeResult(result)).toContain('Repair reinstall (target lookup failed; pinned to installed RunX 0.2.7)')
  })

  test('renders a complete aligned release table', () => {
    const catalog: ReleaseCatalog = {
      schemaVersion: 1,
      command: 'runx upgrade list',
      currentVersion: '0.2.7',
      latestStableVersion: '0.2.8',
      releases: [
        { tag: '@guiho/runx@0.3.0-alpha.1', version: '0.3.0-alpha.1', channel: 'alpha', prerelease: true, publishedAt: '2026-07-15T00:00:00Z', current: false, latestStable: false, compatibleAsset: null },
        { tag: '@guiho/runx@0.2.8', version: '0.2.8', channel: 'stable', prerelease: false, publishedAt: '2026-07-14T00:00:00Z', current: false, latestStable: true, compatibleAsset: { name: 'runx.exe', url: 'https://example.test' } },
      ],
    }
    const output = renderReleaseCatalog(catalog)
    expect(output).toContain('0.3.0-alpha.1  alpha')
    expect(output).toContain('CURRENT')
    expect(output).toContain('LATEST')
  })
})

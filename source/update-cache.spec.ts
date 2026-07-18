/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { afterEach, describe, expect, test } from 'bun:test'
import { mkdir, mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { readVersion } from './help.js'
import { readCachedUpdateNotice, runUpdateWorker, updateCachePath } from './update-cache.js'

const originalHome = Bun.env.HOME
let home = ''

afterEach(async () => {
  if (originalHome === undefined) delete Bun.env.HOME
  else Bun.env.HOME = originalHome
  if (home) await rm(home, { recursive: true, force: true })
  home = ''
})

describe('RunX detached update cache', () => {
  test('foreground reads only local decoded cache data', async () => {
    home = await mkdtemp(join(tmpdir(), 'runx-cache-'))
    Bun.env.HOME = home
    await mkdir(join(home, '.guiho', 'runx'), { recursive: true })
    await Bun.write(updateCachePath(), JSON.stringify({
      newVersionAvailable: true,
      latestVersion: '99.0.0',
      upgradeCommand: 'runx upgrade',
      lastCheck: '2026-07-18T00:00:00.000Z',
    }))
    expect(await readCachedUpdateNotice()).toBe('New version available. Run this command to upgrade: runx upgrade')
    await Bun.write(updateCachePath(), '{')
    expect(await readCachedUpdateNotice()).toBeNull()
  })

  test('worker validates releases and writes update and up-to-date shapes', async () => {
    home = await mkdtemp(join(tmpdir(), 'runx-cache-'))
    Bun.env.HOME = home
    const update = await runUpdateWorker({ fetchImpl: releaseFetch('99.0.0'), now: () => new Date('2026-07-18T00:00:00Z') })
    expect(update.newVersionAvailable).toBe(true)
    expect(update.upgradeCommand).toBe('runx upgrade')
    const current = await runUpdateWorker({ fetchImpl: releaseFetch(readVersion()), now: () => new Date('2026-07-18T00:00:00Z') })
    expect(current.newVersionAvailable).toBe(false)
    expect(current.upgradeCommand).toBeUndefined()
  })
})

function releaseFetch(version: string): typeof fetch {
  return (async () => Response.json([{
    tag_name: `@guiho/runx@${version}`,
    draft: false,
    prerelease: false,
    published_at: '2026-07-18T00:00:00Z',
    assets: [{ name: currentAsset(), browser_download_url: 'https://example.test/runx' }],
  }])) as typeof fetch
}

function currentAsset(): string {
  const os = process.platform === 'win32' ? 'windows' : process.platform === 'darwin' ? 'darwin' : 'linux'
  const suffix = os === 'windows' ? '.exe' : ''
  return process.arch === 'arm64' ? `runx-${os}-arm64${suffix}` : `runx-${os}-x64-baseline${suffix}`
}

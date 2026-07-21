/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { afterEach, describe, expect, test } from 'bun:test'
import { mkdir, mkdtemp, rm, utimes } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { readVersion } from './help.js'
import { readCachedUpdateNotice, runUpdateWorker, spawnUpdateWorker, updateCachePath, updateWorkerLockPath } from './update-cache.js'

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

  test('64 concurrent schedules create exactly one detached worker and honor the four-hour cache TTL', async () => {
    home = await mkdtemp(join(tmpdir(), 'runx cache '))
    Bun.env.HOME = home
    const leaseTokens: string[] = []
    const spawn = ((command: string[], options: { env?: Record<string, string | undefined> }) => {
      expect(command.filter((argument) => argument === '--check-updates-worker')).toHaveLength(1)
      leaseTokens.push(options.env?.RUNX_UPDATE_WORKER_LEASE_TOKEN ?? '')
      return { unref() {} }
    }) as unknown as typeof Bun.spawn

    const scheduled = await Promise.all(Array.from({ length: 64 }, () => spawnUpdateWorker({ spawn })))
    expect(scheduled.filter(Boolean)).toHaveLength(1)
    expect(leaseTokens).toHaveLength(1)
    expect(leaseTokens[0]).not.toBe('')
    expect(await Bun.file(join(updateWorkerLockPath(), 'lease.json')).exists()).toBe(true)

    const updated = await runUpdateWorker({
      fetchImpl: releaseFetch(readVersion()),
      leaseToken: leaseTokens[0],
      now: () => new Date('2026-07-21T00:00:00Z'),
    })
    expect(updated?.newVersionAvailable).toBe(false)
    expect(await Bun.file(join(updateWorkerLockPath(), 'lease.json')).exists()).toBe(false)

    let freshCacheSpawns = 0
    expect(await spawnUpdateWorker({
      now: () => new Date('2026-07-21T03:59:59Z').getTime(),
      spawn: (() => {
        freshCacheSpawns += 1
        return { unref() {} }
      }) as unknown as typeof Bun.spawn,
    })).toBe(false)
    expect(freshCacheSpawns).toBe(0)
  })

  test('stale lease recovery is ownership-safe and a suspended worker cannot remove its successor', async () => {
    home = await mkdtemp(join(tmpdir(), 'runx-cache-'))
    Bun.env.HOME = home
    const leaseTokens: string[] = []
    const spawn = ((_command: string[], options: { env?: Record<string, string | undefined> }) => {
      leaseTokens.push(options.env?.RUNX_UPDATE_WORKER_LEASE_TOKEN ?? '')
      return { unref() {} }
    }) as unknown as typeof Bun.spawn

    expect(await spawnUpdateWorker({ now: () => 1_000, spawn })).toBe(true)
    let continueOldWorker!: () => void
    let reportOldWorkerStarted!: () => void
    const oldWorkerGate = new Promise<void>((resolve) => { continueOldWorker = resolve })
    const oldWorkerStarted = new Promise<void>((resolve) => { reportOldWorkerStarted = resolve })
    const oldFetch = (async (input: string | URL | Request, init?: RequestInit) => {
      reportOldWorkerStarted()
      await oldWorkerGate
      return releaseFetch(readVersion())(input, init)
    }) as typeof fetch
    const oldWorker = runUpdateWorker({ fetchImpl: oldFetch, leaseToken: leaseTokens[0], timeoutMilliseconds: 1_000 })
    await oldWorkerStarted

    const reclaimers = await Promise.all(Array.from({ length: 32 }, () =>
      spawnUpdateWorker({ now: () => 31_001, spawn })))
    expect(reclaimers.filter(Boolean)).toHaveLength(1)
    expect(leaseTokens).toHaveLength(2)
    expect(leaseTokens[0]).not.toBe(leaseTokens[1])

    continueOldWorker()
    expect(await oldWorker).not.toBeNull()
    expect(await Bun.file(join(updateWorkerLockPath(), 'lease.json')).exists()).toBe(true)
    expect(await runUpdateWorker({ fetchImpl: releaseFetch(readVersion()), leaseToken: leaseTokens[1] })).not.toBeNull()
    expect(await Bun.file(join(updateWorkerLockPath(), 'lease.json')).exists()).toBe(false)
  })

  test('hard deadline releases the lease and scheduling failures remain isolated from the foreground', async () => {
    home = await mkdtemp(join(tmpdir(), 'runx-cache-'))
    Bun.env.HOME = home

    const started = performance.now()
    await expect(runUpdateWorker({
      fetchImpl: (() => new Promise<Response>(() => {})) as typeof fetch,
      timeoutMilliseconds: 20,
    })).rejects.toThrow('timed out after 20ms')
    expect(performance.now() - started).toBeLessThan(1_000)
    expect(await Bun.file(join(updateWorkerLockPath(), 'lease.json')).exists()).toBe(false)

    const originalHome = Bun.env.HOME
    const blockedHome = join(home, 'not-a-directory')
    await Bun.write(blockedHome, 'blocked')
    Bun.env.HOME = blockedHome
    try {
      expect(await spawnUpdateWorker()).toBe(false)
    } finally {
      Bun.env.HOME = originalHome
    }
  })

  test('stale missing and malformed lease metadata recover without accumulating workers', async () => {
    home = await mkdtemp(join(tmpdir(), 'runx-cache-'))
    Bun.env.HOME = home

    for (const metadata of [null, '{']) {
      await mkdir(updateWorkerLockPath(), { recursive: true })
      if (metadata !== null) await Bun.write(join(updateWorkerLockPath(), 'lease.json'), metadata)
      const old = new Date('2026-07-20T23:59:00Z')
      await utimes(updateWorkerLockPath(), old, old)
      let leaseToken = ''
      const spawn = ((_command: string[], options: { env?: Record<string, string | undefined> }) => {
        leaseToken = options.env?.RUNX_UPDATE_WORKER_LEASE_TOKEN ?? ''
        return { unref() {} }
      }) as unknown as typeof Bun.spawn

      expect(await spawnUpdateWorker({ now: () => new Date('2026-07-21T00:00:00Z').getTime(), spawn })).toBe(true)
      expect(leaseToken).not.toBe('')
      expect(await runUpdateWorker({ fetchImpl: releaseFetch(readVersion()), leaseToken })).not.toBeNull()
      expect(await Bun.file(join(updateWorkerLockPath(), 'lease.json')).exists()).toBe(false)
    }
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

/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'
import { gt, valid } from 'semver'
import { readVersion } from './help.js'
import { fetchReleaseCatalog, resolveUpgradePlatform } from './release-catalog.js'
import { globalRunXDirectory, movePath, readTextIfExists, writeTextFile } from './storage.js'
import { joinPath } from './path-utils.js'

import type { Static } from '@sinclair/typebox'

export {
  cacheSchema,
  readCachedUpdateNotice,
  runUpdateWorker,
  spawnUpdateWorker,
  updateCachePath,
}
export type {
  UpdateCache,
}

const cacheSchema = Type.Object({
  newVersionAvailable: Type.Boolean(),
  latestVersion: Type.String(),
  upgradeCommand: Type.Optional(Type.String()),
  lastCheck: Type.String(),
}, { additionalProperties: false })

type UpdateCache = Static<typeof cacheSchema>

function updateCachePath(): string {
  return joinPath(globalRunXDirectory(), 'cache.json')
}

async function readCachedUpdateNotice(verbose = false): Promise<string | null> {
  const raw = await readTextIfExists(updateCachePath())
  if (raw === null) return null
  try {
    const cache = Value.Decode(cacheSchema, JSON.parse(raw))
    return cache.newVersionAvailable
      ? 'New version available. Run this command to upgrade: runx upgrade'
      : null
  } catch (error) {
    if (verbose) process.stderr.write(`warning: ignored invalid RunX update cache: ${error instanceof Error ? error.message : String(error)}\n`)
    return null
  }
}

function spawnUpdateWorker(): void {
  if (Bun.env.RUNX_DISABLE_UPDATE_WORKER === '1') return
  const command = process.execPath.toLowerCase().includes('bun')
    ? [process.execPath, Bun.main, '--check-updates-worker']
    : [process.execPath, '--check-updates-worker']
  Bun.spawn(command, {
    detached: true,
    stdout: 'ignore',
    stderr: 'ignore',
    stdin: 'ignore',
    env: process.env,
  }).unref()
}

async function runUpdateWorker(options: { fetchImpl?: typeof fetch, now?: () => Date } = {}): Promise<UpdateCache> {
  const currentVersion = readVersion()
  const platform = resolveUpgradePlatform()
  const catalog = await fetchReleaseCatalog({ ...platform, currentVersion, fetchImpl: options.fetchImpl })
  const latestVersion = catalog.latestStableVersion ?? currentVersion
  const newVersionAvailable = Boolean(valid(currentVersion) && valid(latestVersion) && gt(latestVersion, currentVersion))
  const cache: UpdateCache = {
    newVersionAvailable,
    latestVersion,
    ...(newVersionAvailable ? { upgradeCommand: 'runx upgrade' } : {}),
    lastCheck: (options.now?.() ?? new Date()).toISOString(),
  }
  const path = updateCachePath()
  const temporary = `${path}.${process.pid}.${crypto.randomUUID()}.tmp`
  await writeTextFile(temporary, `${JSON.stringify(cache, null, 2)}\n`)
  await movePath(temporary, path)
  return cache
}

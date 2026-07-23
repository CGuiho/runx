/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { $ } from 'bun'
import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'
import { gt, valid } from 'semver'
import { readVersion } from './help.js'
import { fetchReleaseCatalog, resolveUpgradePlatform } from './release-catalog.js'
import { ensureDirectory, globalRunXDirectory, movePath, readTextIfExists, removePath, writeTextFile } from './storage.js'
import { joinPath } from './path-utils.js'

import type { Static } from '@sinclair/typebox'

export {
  cacheSchema,
  readCachedUpdateNotice,
  runUpdateWorker,
  spawnUpdateWorker,
  updateCachePath,
  updateWorkerLockPath,
}
export type {
  UpdateCache,
}

type UpdateWorkerLease = {
  readonly createdAt: string
  readonly path: string
  readonly pid: number
  readonly token: string
}

type UpdateWorkerSpawnOptions = {
  readonly now?: () => number
  readonly spawn?: typeof Bun.spawn
}

type UpdateWorkerRunOptions = {
  readonly fetchImpl?: typeof fetch
  readonly leaseToken?: string
  readonly now?: () => Date
  readonly timeoutMilliseconds?: number
}

const cacheSchema = Type.Object({
  newVersionAvailable: Type.Boolean(),
  latestVersion: Type.String(),
  upgradeCommand: Type.Optional(Type.String()),
  lastCheck: Type.String(),
}, { additionalProperties: false })

type UpdateCache = Static<typeof cacheSchema>

const cacheTtlMilliseconds = 4 * 60 * 60 * 1000
const fetchTimeoutMilliseconds = 15_000
const staleLockMilliseconds = 30_000
const workerLeaseTokenEnvironment = 'RUNX_UPDATE_WORKER_LEASE_TOKEN'
const updateWorkerLeaseSchema = Type.Object({
  token: Type.String({ minLength: 1 }),
  pid: Type.Integer({ minimum: 1 }),
  createdAt: Type.String({ minLength: 1 }),
}, { additionalProperties: false })

function updateCachePath(): string {
  return joinPath(globalRunXDirectory(), 'cache.json')
}

function updateWorkerLockPath(): string {
  return joinPath(globalRunXDirectory(), '.update-check.lock')
}

async function readCachedUpdateNotice(verbose = false): Promise<string | null> {
  const raw = await readTextIfExists(updateCachePath())
  if (raw === null) return null
  try {
    const cache = Value.Decode(cacheSchema, JSON.parse(raw))
    const currentVersion = readVersion()
    return cache.newVersionAvailable && Boolean(valid(currentVersion) && valid(cache.latestVersion) && gt(cache.latestVersion, currentVersion))
      ? `  ⚠ New version available: v${cache.latestVersion}\n    Run \`${cache.upgradeCommand ?? 'runx upgrade'}\` to update.`
      : null
  } catch (error) {
    if (verbose) process.stderr.write(`warning: ignored invalid RunX update cache: ${error instanceof Error ? error.message : String(error)}\n`)
    return null
  }
}

async function spawnUpdateWorker(options: UpdateWorkerSpawnOptions = {}): Promise<boolean> {
  let lease: UpdateWorkerLease | null = null
  try {
    if (Bun.env.RUNX_DISABLE_UPDATE_WORKER === '1') return false
    const now = options.now?.() ?? Date.now()
    if (await hasFreshCache(now)) return false
    lease = await acquireUpdateWorkerLease(now)
    if (lease === null) return false
    const command = process.execPath.toLowerCase().includes('bun')
      ? [process.execPath, Bun.main, '--check-updates-worker']
      : [process.execPath, '--check-updates-worker']
    ;(options.spawn ?? Bun.spawn)(command, {
      detached: true,
      stdout: 'ignore',
      stderr: 'ignore',
      stdin: 'ignore',
      env: { ...process.env, [workerLeaseTokenEnvironment]: lease.token },
    }).unref()
    return true
  } catch {
    if (lease !== null) await releaseUpdateWorkerLease(lease).catch(() => undefined)
    return false
  }
}

async function runUpdateWorker(options: UpdateWorkerRunOptions = {}): Promise<UpdateCache | null> {
  const lease = options.leaseToken
    ? await readUpdateWorkerLease(options.leaseToken)
    : await acquireUpdateWorkerLease(Date.now())
  if (lease === null) return null

  const controller = new AbortController()
  const timeoutMilliseconds = options.timeoutMilliseconds ?? fetchTimeoutMilliseconds
  const timeout = setTimeout(() => controller.abort(), timeoutMilliseconds)
  const fetchImpl = options.fetchImpl ?? fetch
  const boundedFetch = ((input, init) => fetchImpl(input, { ...init, signal: controller.signal })) as typeof fetch

  try {
    const currentVersion = readVersion()
    const platform = resolveUpgradePlatform()
    const catalog = await Promise.race([
      fetchReleaseCatalog({ ...platform, currentVersion, fetchImpl: boundedFetch }),
      new Promise<never>((_, reject) => {
        const abort = () => reject(new Error(`RunX update check timed out after ${timeoutMilliseconds}ms.`))
        controller.signal.addEventListener('abort', abort, { once: true })
      }),
    ])
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
  } finally {
    clearTimeout(timeout)
    await releaseUpdateWorkerLease(lease)
  }
}

async function hasFreshCache(now: number): Promise<boolean> {
  const raw = await readTextIfExists(updateCachePath())
  if (raw === null) return false
  try {
    const cache = Value.Decode(cacheSchema, JSON.parse(raw))
    const lastCheck = Date.parse(cache.lastCheck)
    const age = now - lastCheck
    return Number.isFinite(lastCheck) && age >= 0 && age < cacheTtlMilliseconds
  } catch {
    return false
  }
}

async function acquireUpdateWorkerLease(now: number): Promise<UpdateWorkerLease | null> {
  const guardPath = await acquireUpdateMutationGuard()
  if (guardPath === null) return null
  try {
    return await acquireUpdateWorkerLeaseInsideGuard(now)
  } finally {
    await removePath(guardPath).catch(() => undefined)
  }
}

async function acquireUpdateWorkerLeaseInsideGuard(now: number): Promise<UpdateWorkerLease | null> {
  const path = updateWorkerLockPath()
  await ensureDirectory(globalRunXDirectory())
  for (let attempt = 0; attempt < 2; attempt += 1) {
    const created = await createLockDirectory(path)
    if (created) return writeUpdateWorkerLease(path, now)
    const existing = await readUpdateWorkerLease()
    if (existing) {
      const createdAt = Date.parse(existing.createdAt)
      if (!Number.isFinite(createdAt) || now - createdAt < staleLockMilliseconds) return null
      if (!(await releaseUpdateWorkerLeaseInsideGuard(existing))) return null
    } else {
      const modifiedAt = await lockModifiedAt(path)
      if (modifiedAt === null || now - modifiedAt < staleLockMilliseconds) return null
      await removePath(path)
    }
  }
  return null
}

async function writeUpdateWorkerLease(path: string, now: number): Promise<UpdateWorkerLease | null> {
  const lease = {
    createdAt: new Date(now).toISOString(),
    path,
    pid: process.pid,
    token: crypto.randomUUID(),
  }
  try {
    await Bun.write(joinPath(path, 'lease.json'), `${JSON.stringify({
      token: lease.token,
      pid: lease.pid,
      createdAt: lease.createdAt,
    }, null, 2)}\n`)
    return lease
  } catch {
    await removePath(path)
    return null
  }
}

async function createLockDirectory(path: string): Promise<boolean> {
  return (await $`mkdir ${path}`.quiet().nothrow()).exitCode === 0
}

async function releaseUpdateWorkerLease(lease: UpdateWorkerLease): Promise<void> {
  const guardPath = await acquireUpdateMutationGuard()
  if (guardPath === null) return
  try {
    await releaseUpdateWorkerLeaseInsideGuard(lease)
  } finally {
    await removePath(guardPath).catch(() => undefined)
  }
}

async function lockModifiedAt(path: string): Promise<number | null> {
  try {
    return (await Bun.file(path).stat()).mtimeMs
  } catch {
    return null
  }
}

async function releaseUpdateWorkerLeaseInsideGuard(lease: UpdateWorkerLease): Promise<boolean> {
  const current = await readUpdateWorkerLease()
  if (current?.token !== lease.token) return false
  await removePath(lease.path)
  return true
}

async function acquireUpdateMutationGuard(): Promise<string | null> {
  const path = joinPath(globalRunXDirectory(), '.update-check.mutation.lock')
  await ensureDirectory(globalRunXDirectory())
  return await createLockDirectory(path) ? path : null
}

async function readUpdateWorkerLease(expectedToken?: string, path = updateWorkerLockPath()): Promise<UpdateWorkerLease | null> {
  const raw = await readTextIfExists(joinPath(path, 'lease.json'))
  if (raw === null) return null
  try {
    const value = Value.Decode(updateWorkerLeaseSchema, JSON.parse(raw))
    if (expectedToken !== undefined && value.token !== expectedToken) return null
    return { ...value, path }
  } catch {
    return null
  }
}

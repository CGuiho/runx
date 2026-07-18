/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { Buffer } from 'node:buffer'
import { chmod, rename, rm } from 'node:fs/promises'
import { basename } from 'node:path'
import { compare, gt, valid } from 'semver'
import { RunXError } from './errors.js'
import { readVersion } from './help.js'
import { fetchReleaseCatalog, resolveUpgradePlatform } from './release-catalog.js'
import { createRecoveryInstructions } from './recovery.js'
import type { ReleaseCatalog, ReleaseCatalogEntry, RecoveryInstructions, UpgradeEnvelope, UpgradeError, UpgradeErrorCode, UpgradeEvent, UpgradePhase, UpgradePlan, UpgradeResult } from './upgrade-types.js'
import type { UpdateResult } from './types.js'

export { checkForLatestVersion, listAvailableVersions, uninstallSelf, upgradeSelf, validateNativeBinary }

export type UpgradeFileOperations = {
  rename: (from: string, to: string) => Promise<void>
  remove: (path: string) => Promise<void>
  makeExecutable: (path: string) => Promise<void>
}

export type UpgradeSelfOptions = {
  dryRun: boolean
  currentVersion?: string
  onPlan?: (plan: UpgradePlan) => void
  onEvent?: (event: UpgradeEvent) => void
  fileOperations?: UpgradeFileOperations
}

type ReplacementState = { backupPath: string, originalMoved: boolean }

class UpgradeOperationError extends Error {
  readonly code: UpgradeErrorCode

  constructor(code: UpgradeErrorCode, message: string) {
    super(message)
    this.name = 'UpgradeOperationError'
    this.code = code
  }
}

const defaultFileOperations: UpgradeFileOperations = {
  rename: async (from, to) => rename(from, to),
  remove: async (path) => rm(path, { force: true }),
  makeExecutable: async (path) => chmod(path, 0o755),
}

const checkForLatestVersion = async (): Promise<UpdateResult> => {
  const catalog = await listAvailableVersions()
  const latest = catalog.releases.find((release) => release.latestStable)
  return {
    currentVersion: catalog.currentVersion,
    latestVersion: catalog.latestStableVersion ?? catalog.currentVersion,
    updateAvailable: Boolean(latest && valid(latest.version) && valid(catalog.currentVersion) && gt(latest.version, catalog.currentVersion)),
    url: latest?.compatibleAsset?.url,
  }
}

const listAvailableVersions = async (): Promise<ReleaseCatalog> => {
  const platform = resolveUpgradePlatform()
  return fetchReleaseCatalog({ ...platform, currentVersion: readVersion() })
}

const upgradeSelf = async (input: boolean | UpgradeSelfOptions): Promise<UpgradeEnvelope> => {
  const options: UpgradeSelfOptions = typeof input === 'boolean' ? { dryRun: input } : input
  const fileOperations = options.fileOperations ?? defaultFileOperations
  const currentVersion = options.currentVersion ?? readVersion()
  const platform = resolveUpgradePlatform()
  const executablePath = resolveSelfExecutable()
  const events: UpgradeEvent[] = []
  let phase: UpgradePhase = 'plan'
  let plan: UpgradePlan | null = null
  let recovery = createRecoveryInstructions(currentVersion, platform.os, 'fallback-current')
  let temporaryPath: string | null = null
  let replacement: ReplacementState | null = null
  let mutationCode: UpgradeErrorCode | null = null

  const emit = (eventPhase: UpgradePhase, status: UpgradeEvent['status'], message?: string): void => {
    phase = eventPhase
    const event: UpgradeEvent = { sequence: events.length + 1, phase: eventPhase, status, ...(message ? { message } : {}) }
    events.push(event)
    options.onEvent?.(event)
  }

  emit('plan', 'started')
  try {
    let catalog: ReleaseCatalog
    try {
      catalog = await fetchReleaseCatalog({ ...platform, currentVersion })
    } catch (error) {
      throw new UpgradeOperationError(classifyPlanError(error), errorMessage(error))
    }
    const stableTarget = catalog.releases.find((release) => release.latestStable)
    if (!stableTarget || !catalog.latestStableVersion) throw new UpgradeOperationError('release_lookup_failed', 'No stable RunX release is available for upgrade.')
    const target = preventDowngrade(currentVersion, stableTarget, catalog.releases)
    const targetIsCurrent = target.version === currentVersion
    recovery = createRecoveryInstructions(target.version, platform.os, 'resolved')
    if (!targetIsCurrent && !target.compatibleAsset) {
      throw new UpgradeOperationError('no_compatible_asset', `RunX ${target.version} has no compatible ${platform.os} ${platform.arch} asset.`)
    }

    plan = {
      currentVersion,
      targetVersion: target.version,
      os: platform.os,
      arch: platform.arch,
      assetName: target.compatibleAsset?.name ?? '',
      assetUrl: target.compatibleAsset?.url ?? '',
      executablePath,
    }
    emit('plan', 'succeeded')
    options.onPlan?.(plan)

    if (basename(executablePath).toLowerCase().startsWith('bun')) {
      throw new UpgradeOperationError('verification_failed', 'Self-management requires a native RunX executable installed from a GitHub release.')
    }
    if (targetIsCurrent) {
      emitSkippedMutation(events, options.onEvent)
      return envelope('up-to-date', plan, events, { installedVersion: currentVersion, cleanupDeferred: false }, recovery)
    }
    if (options.dryRun) {
      emitSkippedMutation(events, options.onEvent)
      return envelope('dry-run', plan, events, null, recovery)
    }

    emit('download', 'started')
    const response = await fetch(plan.assetUrl)
    if (!response.ok) throw new UpgradeOperationError('download_failed', `Could not download RunX ${plan.targetVersion}: HTTP ${response.status}`)
    const bytes = new Uint8Array(await response.arrayBuffer())
    if (bytes.length === 0) throw new UpgradeOperationError('download_invalid', 'Downloaded RunX executable is empty.')
    emit('download', 'succeeded')

    emit('validate', 'started')
    try {
      validateNativeBinary(bytes, platform.os)
      temporaryPath = `${executablePath}.new-${process.pid}-${crypto.randomUUID()}`
      await Bun.write(temporaryPath, bytes)
      if (platform.os !== 'windows') await fileOperations.makeExecutable(temporaryPath)
    } catch (error) {
      throw asOperationError('download_invalid', error)
    }
    emit('validate', 'succeeded')

    emit('replace', 'started')
    replacement = { backupPath: `${executablePath}.old-${process.pid}-${crypto.randomUUID()}`, originalMoved: false }
    mutationCode = 'backup_failed'
    try {
      await fileOperations.rename(executablePath, replacement.backupPath)
      replacement.originalMoved = true
      mutationCode = 'replace_failed'
      await fileOperations.rename(temporaryPath, executablePath)
      temporaryPath = null
    } catch (error) {
      throw asOperationError(mutationCode, error)
    }
    emit('replace', 'succeeded')

    emit('verify', 'started')
    try {
      await verifyExecutableVersion(executablePath, plan.targetVersion)
    } catch (error) {
      throw asOperationError('verification_failed', error)
    }
    emit('verify', 'succeeded')
    emit('cache', 'skipped', 'RunX does not use an upgrade cache.')

    emit('cleanup', 'started')
    const cleanupDeferred = await cleanupBackup(replacement.backupPath, platform.os, fileOperations)
    emit('cleanup', cleanupDeferred ? 'skipped' : 'succeeded', cleanupDeferred ? 'Old executable deletion was deferred; replacement already succeeded.' : undefined)
    return envelope('upgraded', plan, events, { installedVersion: plan.targetVersion, cleanupDeferred }, recovery)
  } catch (error) {
    const primary = asOperationError(classifyFailure(phase, mutationCode), error)
    if (!events.some((event) => event.phase === phase && event.status === 'failed')) emit(phase, 'failed', primary.message)
    if (temporaryPath) await fileOperations.remove(temporaryPath).catch(() => undefined)
    if (replacement?.originalMoved) {
      try {
        await rollbackReplacement(replacement.backupPath, executablePath, fileOperations)
        return envelope('rolled-back', plan, events, { installedVersion: currentVersion, cleanupDeferred: false }, recovery, {
          code: primary.code,
          phase,
          message: primary.message,
        })
      } catch (rollbackError) {
        return envelope('failed', plan, events, null, recovery, {
          code: 'rollback_failed',
          phase,
          message: `${primary.message}. Automatic rollback failed: ${errorMessage(rollbackError)}. Canonical path: ${executablePath}. Backup path: ${replacement.backupPath}.`,
        })
      }
    }
    return envelope('failed', plan, events, null, recovery, { code: primary.code, phase, message: primary.message })
  }
}

const uninstallSelf = async (dryRun: boolean): Promise<{ executablePath: string, scheduled: boolean, dryRun: boolean }> => {
  const executablePath = resolveSelfExecutable()
  if (basename(executablePath).toLowerCase().startsWith('bun')) throw new RunXError('Self-management requires a native RunX executable installed from a GitHub release.')
  if (dryRun) return { executablePath, scheduled: false, dryRun: true }
  if (process.platform === 'win32') {
    const command = `ping 127.0.0.1 -n 2 > nul & del /f /q "${executablePath}"`
    Bun.spawn(['cmd.exe', '/d', '/s', '/c', command], { detached: true, stdout: 'ignore', stderr: 'ignore', stdin: 'ignore' })
    return { executablePath, scheduled: true, dryRun: false }
  }
  await rm(executablePath)
  return { executablePath, scheduled: false, dryRun: false }
}

const envelope = (
  outcome: UpgradeEnvelope['outcome'],
  plan: UpgradePlan | null,
  events: UpgradeEvent[],
  result: UpgradeResult | null,
  recovery: RecoveryInstructions,
  error: UpgradeError | null = null,
): UpgradeEnvelope => ({ schemaVersion: 1, command: 'runx upgrade', outcome, plan, events, result, recovery, error })

const preventDowngrade = (currentVersion: string, stableTarget: ReleaseCatalogEntry, releases: ReleaseCatalogEntry[]): ReleaseCatalogEntry => {
  if (!valid(currentVersion) || !valid(stableTarget.version) || compare(currentVersion, stableTarget.version) < 0) return stableTarget
  return releases.find((release) => release.version === currentVersion) ?? {
    tag: currentVersion,
    version: currentVersion,
    channel: valid(currentVersion)?.includes('-') ? 'prerelease' : 'stable',
    prerelease: Boolean(valid(currentVersion)?.includes('-')),
    publishedAt: null,
    current: true,
    latestStable: false,
    compatibleAsset: null,
  }
}

const emitSkippedMutation = (events: UpgradeEvent[], onEvent?: (event: UpgradeEvent) => void): void => {
  for (const phase of ['download', 'validate', 'replace', 'verify', 'cache', 'cleanup'] as const) {
    const event: UpgradeEvent = { sequence: events.length + 1, phase, status: 'skipped' }
    events.push(event)
    onEvent?.(event)
  }
}

const rollbackReplacement = async (backupPath: string, executablePath: string, operations: UpgradeFileOperations): Promise<void> => {
  await operations.remove(executablePath)
  await operations.rename(backupPath, executablePath)
}

const verifyExecutableVersion = async (executablePath: string, expectedVersion: string): Promise<void> => {
  const verification = Bun.spawn([executablePath, '--version'], { stdin: 'ignore', stdout: 'pipe', stderr: 'pipe' })
  let timeout: ReturnType<typeof setTimeout> | undefined
  const exited = Promise.race([
    verification.exited,
    new Promise<never>((_, reject) => {
      timeout = setTimeout(() => {
        verification.kill()
        reject(new Error('replacement version check timed out after 10 seconds'))
      }, 10_000)
    }),
  ])
  const [stdout, stderr, exitCode] = await Promise.all([
    new Response(verification.stdout).text(),
    new Response(verification.stderr).text(),
    exited,
  ]).finally(() => {
    if (timeout !== undefined) clearTimeout(timeout)
  })
  const actualVersion = stdout.trim()
  if (exitCode !== 0) throw new Error(`replacement exited with code ${exitCode}${stderr.trim() ? `: ${stderr.trim()}` : ''}`)
  if (actualVersion !== expectedVersion) throw new Error(`replacement reported version ${actualVersion || '<empty>'}; expected ${expectedVersion}`)
}

const cleanupBackup = async (backupPath: string, os: UpgradePlan['os'], operations: UpgradeFileOperations): Promise<boolean> => {
  try {
    await operations.remove(backupPath)
    return false
  } catch {
    if (os !== 'windows') throw new UpgradeOperationError('replace_failed', `Could not delete old RunX executable: ${backupPath}`)
    scheduleWindowsBackupCleanup(backupPath)
    return true
  }
}

const scheduleWindowsBackupCleanup = (backupPath: string): void => {
  const command = 'for ($attempt = 0; $attempt -lt 300; $attempt += 1) { if (-not (Test-Path -LiteralPath $env:RUNX_BACKUP_PATH)) { exit 0 }; try { Remove-Item -LiteralPath $env:RUNX_BACKUP_PATH -Force -ErrorAction Stop; exit 0 } catch { Start-Sleep -Milliseconds 100 } }; exit 1'
  const encodedCommand = Buffer.from(command, 'utf16le').toString('base64')
  Bun.spawn(['cmd.exe', '/d', '/s', '/c', `powershell.exe -NoLogo -NoProfile -NonInteractive -WindowStyle Hidden -EncodedCommand ${encodedCommand}`], {
    detached: true,
    env: { ...process.env, RUNX_BACKUP_PATH: backupPath },
    stdout: 'ignore', stderr: 'ignore', stdin: 'ignore',
  })
}

const validateNativeBinary = (bytes: Uint8Array, os: UpgradePlan['os']): void => {
  const native = os === 'windows'
    ? bytes[0] === 0x4d && bytes[1] === 0x5a
    : os === 'linux'
      ? bytes[0] === 0x7f && bytes[1] === 0x45 && bytes[2] === 0x4c && bytes[3] === 0x46
      : isMachO(bytes)
  if (!native) throw new RunXError(`Downloaded file is not a native ${os} executable.`)
}

const isMachO = (bytes: Uint8Array): boolean => {
  const magic = Array.from(bytes.slice(0, 4)).map((value) => value.toString(16).padStart(2, '0')).join('')
  return ['feedface', 'feedfacf', 'cefaedfe', 'cffaedfe', 'cafebabe', 'bebafeca'].includes(magic)
}

const classifyPlanError = (error: unknown): UpgradeErrorCode => /malformed/.test(errorMessage(error)) ? 'release_payload_invalid' : 'release_lookup_failed'
const classifyFailure = (phase: UpgradePhase, mutationCode: UpgradeErrorCode | null): UpgradeErrorCode => {
  if (phase === 'download') return 'download_failed'
  if (phase === 'validate') return 'download_invalid'
  if (phase === 'replace') return mutationCode ?? 'replace_failed'
  if (phase === 'verify') return 'verification_failed'
  return 'release_lookup_failed'
}
const asOperationError = (code: UpgradeErrorCode, error: unknown): UpgradeOperationError => error instanceof UpgradeOperationError ? error : new UpgradeOperationError(code, errorMessage(error))
const resolveSelfExecutable = (): string => process.env['RUNX_SELF_PATH'] ?? process.execPath
const errorMessage = (error: unknown): string => error instanceof Error ? error.message : String(error)

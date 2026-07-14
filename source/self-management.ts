import { Buffer } from 'node:buffer'
import { chmod, rename, rm } from 'node:fs/promises'
import { basename } from 'node:path'
import { RunXError } from './errors.js'
import { readVersion } from './help.js'
import type { RunXUpgradeResult, UpdateResult } from './types.js'

export {
  checkForLatestVersion,
  listAvailableVersions,
  uninstallSelf,
  upgradeSelf,
}

const repositoryApi = 'https://api.github.com/repos/CGuiho/runx/releases'

type Release = { tag_name: string, assets: Array<{ name: string, browser_download_url: string }> }

const checkForLatestVersion = async (): Promise<UpdateResult> => {
  const currentVersion = readVersion()
  try {
    const response = await fetch(`${repositoryApi}/latest`, { headers: { accept: 'application/vnd.github+json' } })
    if (!response.ok) return { currentVersion, latestVersion: currentVersion, updateAvailable: false }
    const release = await response.json() as Release
    const latestVersion = release.tag_name.replace(/^@guiho\/runx@/, '').replace(/^v/, '')
    const asset = findAsset(release)
    return { currentVersion, latestVersion, updateAvailable: compareVersions(latestVersion, currentVersion) > 0, url: asset?.browser_download_url }
  } catch {
    return { currentVersion, latestVersion: currentVersion, updateAvailable: false }
  }
}

const listAvailableVersions = async (): Promise<string[]> => {
  const response = await fetch(`${repositoryApi}?per_page=20`, { headers: { accept: 'application/vnd.github+json' } })
  if (!response.ok) throw new RunXError(`Could not retrieve RunX releases: HTTP ${response.status}`)
  const releases = await response.json() as Release[]
  return releases.map((release) => release.tag_name.replace(/^@guiho\/runx@/, '').replace(/^v/, ''))
}

const upgradeSelf = async (dryRun: boolean): Promise<RunXUpgradeResult> => {
  const executablePath = requireNativeExecutable()
  const update = await checkForLatestVersion()
  const upToDate = !update.updateAvailable && update.latestVersion === update.currentVersion && Boolean(update.url)
  if (!update.updateAvailable || !update.url) return { ...update, executablePath, scheduled: false, upToDate }
  if (dryRun) return { ...update, executablePath, scheduled: false, upToDate: false }

  const response = await fetch(update.url)
  if (!response.ok) throw new RunXError(`Could not download RunX update: HTTP ${response.status}`)
  const temporaryPath = `${executablePath}.new`
  await Bun.write(temporaryPath, await response.arrayBuffer())
  if (process.platform !== 'win32') {
    await chmod(temporaryPath, 0o755)
    await rename(temporaryPath, executablePath)
    return { ...update, executablePath, scheduled: false, upToDate: false }
  }

  await replaceWindowsExecutable(temporaryPath, executablePath, update.latestVersion)
  return { ...update, executablePath, scheduled: false, upToDate: false }
}

const uninstallSelf = async (dryRun: boolean): Promise<{ executablePath: string, scheduled: boolean, dryRun: boolean }> => {
  const executablePath = requireNativeExecutable()
  if (dryRun) return { executablePath, scheduled: false, dryRun: true }
  if (process.platform === 'win32') {
    const command = `ping 127.0.0.1 -n 2 > nul & del /f /q "${executablePath}"`
    Bun.spawn(['cmd.exe', '/d', '/s', '/c', command], { detached: true, stdout: 'ignore', stderr: 'ignore', stdin: 'ignore' })
    return { executablePath, scheduled: true, dryRun: false }
  }

  await rm(executablePath)
  return { executablePath, scheduled: false, dryRun: false }
}

const findAsset = (release: Release): Release['assets'][number] | undefined => release.assets.find((asset) => asset.name === assetName())

const assetName = (): string => `runx-${process.platform === 'win32' ? 'windows' : process.platform === 'darwin' ? 'darwin' : 'linux'}-${process.arch}${process.platform === 'win32' ? '.exe' : ''}`

const requireNativeExecutable = (): string => {
  const executablePath = process.env['RUNX_SELF_PATH'] ?? process.execPath
  if (basename(executablePath).toLowerCase().startsWith('bun')) throw new RunXError('Self-management requires a native RunX executable installed from a GitHub release.')
  return executablePath
}

const replaceWindowsExecutable = async (temporaryPath: string, executablePath: string, expectedVersion: string): Promise<void> => {
  const backupPath = `${executablePath}.old`
  let originalMoved = false

  try {
    await rm(backupPath, { force: true })
    await rename(executablePath, backupPath)
    originalMoved = true
    await rename(temporaryPath, executablePath)
    await verifyExecutableVersion(executablePath, expectedVersion)
    await cleanupWindowsBackup(backupPath)
  } catch (error) {
    const replacementFailure = errorMessage(error)
    await rm(temporaryPath, { force: true }).catch(() => undefined)

    if (!originalMoved) {
      throw new RunXError(`Could not replace the Windows RunX executable at ${executablePath}: ${replacementFailure}`)
    }

    try {
      await rm(executablePath, { force: true })
      await rename(backupPath, executablePath)
    } catch (rollbackError) {
      throw new RunXError(`Windows RunX upgrade failed at ${executablePath}: ${replacementFailure}. Automatic rollback also failed: ${errorMessage(rollbackError)}`)
    }

    throw new RunXError(`Windows RunX upgrade failed; the previous executable was restored at ${executablePath}: ${replacementFailure}`)
  }
}

const verifyExecutableVersion = async (executablePath: string, expectedVersion: string): Promise<void> => {
  const verification = Bun.spawn([executablePath, '--version'], { stdin: 'ignore', stdout: 'pipe', stderr: 'pipe' })
  const [stdout, stderr, exitCode] = await Promise.all([
    new Response(verification.stdout).text(),
    new Response(verification.stderr).text(),
    verification.exited,
  ])
  const actualVersion = stdout.trim()
  if (exitCode !== 0) throw new Error(`replacement exited with code ${exitCode}${stderr.trim() ? `: ${stderr.trim()}` : ''}`)
  if (actualVersion !== expectedVersion) throw new Error(`replacement reported version ${actualVersion || '<empty>'}; expected ${expectedVersion}`)
}

const cleanupWindowsBackup = async (backupPath: string): Promise<void> => {
  try {
    await rm(backupPath, { force: true })
  } catch {}
  scheduleWindowsBackupCleanup(backupPath)
}

const scheduleWindowsBackupCleanup = (backupPath: string): void => {
  const command = 'for ($attempt = 0; $attempt -lt 300; $attempt += 1) { if (-not (Test-Path -LiteralPath $env:RUNX_BACKUP_PATH)) { exit 0 }; try { Remove-Item -LiteralPath $env:RUNX_BACKUP_PATH -Force -ErrorAction Stop; exit 0 } catch { Start-Sleep -Milliseconds 100 } }; exit 1'
  const encodedCommand = Buffer.from(command, 'utf16le').toString('base64')
  Bun.spawn(['cmd.exe', '/d', '/s', '/c', `powershell.exe -NoLogo -NoProfile -NonInteractive -WindowStyle Hidden -EncodedCommand ${encodedCommand}`], {
    detached: true,
    env: { ...process.env, RUNX_BACKUP_PATH: backupPath },
    stdout: 'ignore',
    stderr: 'ignore',
    stdin: 'ignore',
  })
}

const errorMessage = (error: unknown): string => error instanceof Error ? error.message : String(error)

const compareVersions = (left: string, right: string): number => {
  const parse = (value: string): number[] => value.split(/[.+-]/).map((part) => Number.parseInt(part, 10) || 0)
  const a = parse(left)
  const b = parse(right)
  for (let index = 0; index < Math.max(a.length, b.length); index += 1) {
    const difference = (a[index] ?? 0) - (b[index] ?? 0)
    if (difference !== 0) return difference
  }
  return 0
}

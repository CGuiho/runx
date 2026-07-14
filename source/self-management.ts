import { chmod, rename, rm } from 'node:fs/promises'
import { basename } from 'node:path'
import { RunXError } from './errors.js'
import { readVersion } from './help.js'
import type { RunXUpgradeResult, UpdateResult } from './types.js'

const repositoryApi = 'https://api.github.com/repos/CGuiho/runx/releases'

type Release = { tag_name: string, assets: Array<{ name: string, browser_download_url: string }> }

export const checkForLatestVersion = async (): Promise<UpdateResult> => {
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

export const listAvailableVersions = async (): Promise<string[]> => {
  const response = await fetch(`${repositoryApi}?per_page=20`, { headers: { accept: 'application/vnd.github+json' } })
  if (!response.ok) throw new RunXError(`Could not retrieve RunX releases: HTTP ${response.status}`)
  const releases = await response.json() as Release[]
  return releases.map((release) => release.tag_name.replace(/^@guiho\/runx@/, '').replace(/^v/, ''))
}

export const upgradeSelf = async (dryRun: boolean): Promise<RunXUpgradeResult> => {
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

  scheduleWindowsReplacement(temporaryPath, executablePath)
  return { ...update, executablePath, scheduled: true, upToDate: false }
}

export const uninstallSelf = async (dryRun: boolean): Promise<{ executablePath: string, scheduled: boolean, dryRun: boolean }> => {
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

const scheduleWindowsReplacement = (temporaryPath: string, executablePath: string): void => {
  const command = `ping 127.0.0.1 -n 2 > nul & move /y "${temporaryPath}" "${executablePath}" > nul`
  Bun.spawn(['cmd.exe', '/d', '/s', '/c', command], { detached: true, stdout: 'ignore', stderr: 'ignore', stdin: 'ignore' })
}

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

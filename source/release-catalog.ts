/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { compare, parse, valid } from 'semver'
import { RunXError } from './errors.js'
import type { ReleaseAsset, ReleaseCatalog, ReleaseCatalogEntry, UpgradeArch, UpgradeOs, UpgradeVariant } from './upgrade-types.js'

export {
  assetCandidates,
  fetchReleaseCatalog,
  findCompatibleAsset,
  normalizeReleaseVersion,
  resolveUpgradePlatform,
}

export type GitHubRelease = {
  tag_name: string
  draft?: boolean
  prerelease?: boolean
  published_at?: string | null
  created_at?: string | null
  assets: Array<{ name: string, browser_download_url: string }>
}

export type ReleasePlatform = {
  os: UpgradeOs
  arch: UpgradeArch
  variant: UpgradeVariant
}

type FetchCatalogOptions = ReleasePlatform & {
  currentVersion: string
  fetchImpl?: typeof fetch
}

const repositoryApi = 'https://api.github.com/repos/CGuiho/runx/releases'

const fetchReleaseCatalog = async (options: FetchCatalogOptions): Promise<ReleaseCatalog> => {
  const fetchImpl = options.fetchImpl ?? fetch
  const releases: GitHubRelease[] = []
  let page = 1
  let next = true

  while (next) {
    const response = await fetchImpl(`${repositoryApi}?per_page=100&page=${page}`, {
      headers: { accept: 'application/vnd.github+json' },
    })
    if (!response.ok) throw new RunXError(`Could not retrieve RunX releases page ${page}: HTTP ${response.status}`)
    const payload: unknown = await response.json()
    if (!Array.isArray(payload)) throw new RunXError(`Could not retrieve RunX releases page ${page}: malformed response`)
    const pageReleases = payload as GitHubRelease[]
    if (pageReleases.some((release) => !isGitHubRelease(release))) {
      throw new RunXError(`Could not retrieve RunX releases page ${page}: malformed release`)
    }
    releases.push(...pageReleases.filter((release) => !release.draft))
    next = hasNextPage(response.headers.get('link'))
    page += 1
  }

  const normalized = releases.map((release) => normalizeRelease(release, options))
  normalized.sort(compareCatalogEntries)
  const latestStableVersion = normalized.find((entry) => valid(entry.version) && !entry.prerelease)?.version ?? null

  for (const entry of normalized) entry.latestStable = entry.version === latestStableVersion

  return {
    schemaVersion: 1,
    command: 'runx upgrade list',
    currentVersion: options.currentVersion,
    latestStableVersion,
    releases: normalized,
  }
}

const normalizeRelease = (release: GitHubRelease, platform: ReleasePlatform & { currentVersion: string }): ReleaseCatalogEntry => {
  const version = normalizeReleaseVersion(release.tag_name)
  const parsed = parse(version)
  const prerelease = parsed ? parsed.prerelease.length > 0 : Boolean(release.prerelease)
  return {
    tag: release.tag_name,
    version,
    channel: releaseChannel(version, prerelease),
    prerelease,
    publishedAt: release.published_at ?? release.created_at ?? null,
    current: version === platform.currentVersion,
    latestStable: false,
    compatibleAsset: findCompatibleAsset(release, platform),
  }
}

const findCompatibleAsset = (release: GitHubRelease, platform: ReleasePlatform): ReleaseAsset | null => {
  for (const candidate of assetCandidates(platform)) {
    const asset = release.assets.find((item) => item.name === candidate)
    if (asset) return { name: asset.name, url: asset.browser_download_url }
  }
  return null
}

const assetCandidates = ({ os, arch, variant }: ReleasePlatform): string[] => {
  const suffix = os === 'windows' ? '.exe' : ''
  if (arch === 'arm64') return [`runx-${os}-arm64${suffix}`]
  const prefix = `runx-${os}-x64`
  const names = {
    baseline: [`${prefix}-baseline${suffix}`, `${prefix}${suffix}`, `${prefix}-modern${suffix}`],
    default: [`${prefix}${suffix}`, `${prefix}-baseline${suffix}`, `${prefix}-modern${suffix}`],
    modern: [`${prefix}-modern${suffix}`, `${prefix}${suffix}`, `${prefix}-baseline${suffix}`],
  }
  return names[variant]
}

const resolveUpgradePlatform = (platform = process.platform, architecture = process.arch): ReleasePlatform => {
  const os = platform === 'win32' ? 'windows' : platform === 'darwin' ? 'macos' : platform === 'linux' ? 'linux' : null
  const arch = architecture === 'x64' ? 'x64' : architecture === 'arm64' ? 'arm64' : null
  if (!os) throw new RunXError(`Self-upgrade is not supported on ${platform}.`)
  if (!arch) throw new RunXError(`Self-upgrade is not supported on architecture ${architecture}.`)
  return { os, arch, variant: 'baseline' }
}

const normalizeReleaseVersion = (tag: string): string => tag.replace(/^@guiho\/runx@/, '').replace(/^v/, '')

const releaseChannel = (version: string, prerelease: boolean): string => {
  const parsed = parse(version)
  if (!parsed) return 'other'
  if (!prerelease) return 'stable'
  const identifier = parsed.prerelease[0]
  return typeof identifier === 'string' ? identifier : 'prerelease'
}

const compareCatalogEntries = (left: ReleaseCatalogEntry, right: ReleaseCatalogEntry): number => {
  const leftValid = valid(left.version)
  const rightValid = valid(right.version)
  if (leftValid && rightValid) return compare(rightValid, leftValid)
  if (leftValid) return -1
  if (rightValid) return 1
  return timestamp(right.publishedAt) - timestamp(left.publishedAt)
}

const timestamp = (value: string | null): number => value ? Date.parse(value) || 0 : 0

const hasNextPage = (link: string | null): boolean => Boolean(link?.split(',').some((part) => /rel="next"/.test(part)))

const isGitHubRelease = (value: unknown): value is GitHubRelease => {
  if (!value || typeof value !== 'object') return false
  const release = value as Partial<GitHubRelease>
  return typeof release.tag_name === 'string'
    && Array.isArray(release.assets)
    && release.assets.every((asset) => typeof asset?.name === 'string' && typeof asset?.browser_download_url === 'string')
}

import { describe, expect, test } from 'bun:test'
import { assetCandidates, fetchReleaseCatalog, normalizeReleaseVersion } from './release-catalog.js'

describe('release catalog', () => {
  test('paginates every release and sorts valid SemVer before invalid tags', async () => {
    const urls: string[] = []
    const pages = [
      [release('@guiho/runx@1.0.0-alpha.10', '2026-07-15T03:00:00Z'), release('@guiho/runx@1.0.0', '2026-07-15T02:00:00Z')],
      [release('@guiho/runx@1.0.0-alpha.2', '2026-07-15T01:00:00Z'), release('nightly', '2026-07-16T00:00:00Z')],
    ]
    const fetchImpl = (async (input: string | URL | Request) => {
      urls.push(String(input))
      const page = pages[urls.length - 1] ?? []
      return Response.json(page, { headers: urls.length === 1 ? { link: '<https://api.github.test/releases?page=2>; rel="next"' } : {} })
    }) as typeof fetch

    const catalog = await fetchReleaseCatalog({ currentVersion: '1.0.0-alpha.2', os: 'windows', arch: 'x64', variant: 'baseline', fetchImpl })

    expect(urls).toHaveLength(2)
    expect(catalog.latestStableVersion).toBe('1.0.0')
    expect(catalog.releases.map((entry) => entry.version)).toEqual(['1.0.0', '1.0.0-alpha.10', '1.0.0-alpha.2', 'nightly'])
    expect(catalog.releases.map((entry) => entry.channel)).toEqual(['stable', 'alpha', 'alpha', 'other'])
    expect(catalog.releases[2]?.current).toBe(true)
    expect(catalog.releases[0]?.latestStable).toBe(true)
    expect(catalog.releases[0]?.compatibleAsset?.name).toBe('runx-windows-x64-baseline.exe')
  })

  test('rejects transport and malformed payload failures', async () => {
    const unavailable = (async () => new Response('no', { status: 503 })) as typeof fetch
    const malformed = (async () => Response.json({ releases: [] })) as typeof fetch
    await expect(fetchReleaseCatalog({ currentVersion: '1.0.0', os: 'linux', arch: 'x64', variant: 'baseline', fetchImpl: unavailable })).rejects.toThrow('HTTP 503')
    await expect(fetchReleaseCatalog({ currentVersion: '1.0.0', os: 'linux', arch: 'x64', variant: 'baseline', fetchImpl: malformed })).rejects.toThrow('malformed response')
  })

  test('uses the shared candidate policy', () => {
    expect(assetCandidates({ os: 'linux', arch: 'x64', variant: 'modern' })).toEqual([
      'runx-linux-x64-modern',
      'runx-linux-x64',
      'runx-linux-x64-baseline',
    ])
    expect(assetCandidates({ os: 'macos', arch: 'arm64', variant: 'baseline' })).toEqual(['runx-macos-arm64'])
    expect(normalizeReleaseVersion('v2.0.0-rc.1')).toBe('2.0.0-rc.1')
  })
})

const release = (tag: string, publishedAt: string) => ({
  tag_name: tag,
  draft: false,
  prerelease: tag.includes('-'),
  published_at: publishedAt,
  assets: [
    { name: 'runx-windows-x64-baseline.exe', browser_download_url: `https://example.test/${tag}/runx.exe` },
  ],
})

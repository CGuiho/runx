import { afterEach, describe, expect, test } from 'bun:test'
import { mkdir, mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { RunXError } from './errors.js'
import { runCli } from './cli.js'
import { readVersion } from './help.js'
import { readManifest, resolveCommand } from './manifest.js'
import { upgradeSelf } from './self-management.js'

const directories: string[] = []

afterEach(async () => {
  await Promise.all(directories.splice(0).map((directory) => rm(directory, { recursive: true, force: true })))
})

describe('RunX manifests', () => {
  test('enforces SemVer 1.x, a scripts directory, and the public group', async () => {
    const accepted = await fixture(emptyManifest('1.2.3-beta.1+build.7'))
    await expect(readManifest(accepted)).resolves.toMatchObject({ manifest: { version: '1.2.3-beta.1+build.7', commands: [] } })

    const invalidManifests = [
      emptyManifest('2.0.0'),
      emptyManifest('1.0'),
      emptyManifest('1').replace('"1"', '1'),
      emptyManifest().replace('scripts:\n  directory: scripts\n', ''),
      emptyManifest().replace('  directory: scripts', '  directory: ../scripts'),
      emptyManifest().replace('  public:\n    summary: Default public project commands.\n', ''),
      `${emptyManifest()}\n  - uid: private-check\n    id: check\n    group: private\n    summary: Check private work.\n    description: Checks private work.\n    command: echo check\n`,
    ]

    for (const content of invalidManifests) {
      const root = await fixture(content)
      await expect(readManifest(root)).rejects.toThrow(RunXError)
    }
  })

  test('discovers a parent manifest and resolves stable selectors', async () => {
    const root = await fixture(manifest())
    const nested = join(root, 'nested')
    await mkdir(nested)
    const { manifest: loaded, path } = await readManifest(nested)

    expect(resolveCommand(loaded, path, 'dev-start').selector).toBe('development/start')
    expect(resolveCommand(loaded, path, 'development/check').uid).toBe('check')
    expect(resolveCommand(loaded, path, '2').uid).toBe('check')
  })

  test('rejects ambiguous IDs', async () => {
    const root = await fixture(`${manifest()}\n  - uid: maintenance-start\n    id: start\n    group: maintenance\n    summary: Start maintenance.\n    description: Starts local maintenance.\n    command: echo maintenance\n`.replace('  development:\n    summary: Development commands.', '  development:\n    summary: Development commands.\n  maintenance:\n    summary: Maintenance commands.'))
    const { manifest: loaded, path } = await readManifest(root)
    expect(() => resolveCommand(loaded, path, 'start')).toThrow(RunXError)
  })

  test('runs the r alias and keeps dry runs read-only', async () => {
    const root = await fixture(manifest())
    const manifestPath = join(root, 'runx.yaml')
    const cliPath = Bun.fileURLToPath(new URL('./guiho-runx-bin.ts', import.meta.url))
    const dryRun = Bun.spawn([process.execPath, cliPath, 'r', 'dev-start', '--file', manifestPath, '--dry-run'], { stdout: 'pipe', stderr: 'pipe' })
    expect(await dryRun.exited).toBe(0)
    expect(await new Response(dryRun.stdout).text()).toContain('uid: dev-start')

    const run = Bun.spawn([process.execPath, cliPath, 'r', 'dev-start', '--file', manifestPath], { stdout: 'pipe', stderr: 'pipe' })
    expect(await run.exited).toBe(0)
    expect(await new Response(run.stdout).text()).toContain('start')
  })

  test('treats the installed release as already up to date without downloading it', async () => {
    const directory = await mkdtemp(join(tmpdir(), 'runx-upgrade-'))
    directories.push(directory)
    const executable = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    const previousSelfPath = process.env['RUNX_SELF_PATH']
    const originalFetch = globalThis.fetch
    const originalStdoutWrite = process.stdout.write
    let fetchCount = 0

    try {
      process.env['RUNX_SELF_PATH'] = executable
      await Bun.write(executable, process.platform === 'win32' ? 'MZ' : '\x7fELF')
      globalThis.fetch = async () => {
        fetchCount += 1
        return Response.json([{
          tag_name: `@guiho/runx@${readVersion()}`,
          draft: false,
          prerelease: false,
          published_at: '2026-07-15T00:00:00Z',
          assets: [{ name: updateAssetName(), browser_download_url: 'https://example.test/runx' }],
        }])
      }

      const result = await upgradeSelf(false)
      expect(result.outcome).toBe('up-to-date')
      expect(result.error).toBeNull()

      const stdout: string[] = []
      process.stdout.write = ((chunk: string | Uint8Array) => {
        stdout.push(String(chunk))
        return true
      }) as typeof process.stdout.write
      await runCli(['upgrade'])
      const output = stdout.join('')
      expect(output).toContain('Already up to date:')
      expect(output).toContain(process.platform === 'win32' ? '-Version' : '--version')
      expect(fetchCount).toBe(2)
    } finally {
      globalThis.fetch = originalFetch
      process.stdout.write = originalStdoutWrite
      if (previousSelfPath === undefined) delete process.env['RUNX_SELF_PATH']
      else process.env['RUNX_SELF_PATH'] = previousSelfPath
    }
  })

  test('routes upgrade inspection and uninstall dry runs through Citty', async () => {
    const directory = await mkdtemp(join(tmpdir(), 'runx-self-management-'))
    directories.push(directory)
    const executable = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    const previousSelfPath = process.env['RUNX_SELF_PATH']
    const originalFetch = globalThis.fetch
    const originalStdoutWrite = process.stdout.write

    try {
      process.env['RUNX_SELF_PATH'] = executable
      await Bun.write(executable, process.platform === 'win32' ? 'MZ' : '\x7fELF')
      globalThis.fetch = async (input) => {
        const release = {
          tag_name: `@guiho/runx@${readVersion()}`,
          draft: false,
          prerelease: false,
          published_at: '2026-07-15T00:00:00Z',
          assets: [{ name: updateAssetName(), browser_download_url: 'https://example.test/runx' }],
        }
        return Response.json([release])
      }

      const stdout: string[] = []
      process.stdout.write = ((chunk: string | Uint8Array) => {
        stdout.push(String(chunk))
        return true
      }) as typeof process.stdout.write

      await runCli(['upgrade', 'check', '--format=json'])
      expect(JSON.parse(stdout.splice(0).join('')).updateAvailable).toBe(false)

      await runCli(['upgrade', 'list', '--format=json'])
      expect(JSON.parse(stdout.splice(0).join('')).releases.map((release: { version: string }) => release.version)).toEqual([readVersion()])

      await runCli(['uninstall', '--dry-run', '--format=json'])
      expect(JSON.parse(stdout.join('')).dryRun).toBe(true)
      expect(await Bun.file(executable).exists()).toBe(true)
    } finally {
      globalThis.fetch = originalFetch
      process.stdout.write = originalStdoutWrite
      if (previousSelfPath === undefined) delete process.env['RUNX_SELF_PATH']
      else process.env['RUNX_SELF_PATH'] = previousSelfPath
    }
  })

  test('emits one exact failed JSON envelope without a generic trailing error', async () => {
    const previousSelfPath = process.env['RUNX_SELF_PATH']
    const originalFetch = globalThis.fetch
    const originalStdoutWrite = process.stdout.write
    const originalStderrWrite = process.stderr.write
    const originalExitCode = process.exitCode
    const stdout: string[] = []
    const stderr: string[] = []
    try {
      delete process.env['RUNX_SELF_PATH']
      globalThis.fetch = (async () => Response.json([{
        tag_name: `@guiho/runx@${readVersion()}`,
        draft: false,
        prerelease: false,
        published_at: '2026-07-15T00:00:00Z',
        assets: [{ name: updateAssetName(), browser_download_url: 'https://example.test/runx' }],
      }])) as typeof fetch
      process.stdout.write = ((chunk: string | Uint8Array) => { stdout.push(String(chunk)); return true }) as typeof process.stdout.write
      process.stderr.write = ((chunk: string | Uint8Array) => { stderr.push(String(chunk)); return true }) as typeof process.stderr.write
      process.exitCode = undefined

      await runCli(['upgrade', '--format=json'])

      const document = JSON.parse(stdout.join(''))
      expect(document).toMatchObject({
        schemaVersion: 1,
        command: 'runx upgrade',
        outcome: 'failed',
        result: null,
        recovery: { targetVersion: readVersion(), targetSource: 'resolved' },
        error: { code: 'verification_failed', phase: 'plan' },
      })
      expect(Object.keys(document.recovery)).toEqual(['targetVersion', 'targetSource', 'installCommand', 'stopProcessCommand'])
      expect(stderr.join('')).toContain('requires a native RunX executable')
      expect(process.exitCode).toBe(1)
    } finally {
      globalThis.fetch = originalFetch
      process.stdout.write = originalStdoutWrite
      process.stderr.write = originalStderrWrite
      process.exitCode = originalExitCode
      if (previousSelfPath === undefined) delete process.env['RUNX_SELF_PATH']
      else process.env['RUNX_SELF_PATH'] = previousSelfPath
    }
  })
})

const fixture = async (content: string): Promise<string> => {
  const directory = await mkdtemp(join(tmpdir(), 'runx-'))
  directories.push(directory)
  await Bun.write(join(directory, 'runx.yaml'), content)
  return directory
}

const updateAssetName = (): string => {
  const os = process.platform === 'win32' ? 'windows' : process.platform === 'darwin' ? 'macos' : 'linux'
  const suffix = process.platform === 'win32' ? '.exe' : ''
  return process.arch === 'arm64' ? `runx-${os}-arm64${suffix}` : `runx-${os}-x64-baseline${suffix}`
}

const manifest = (): string => `version: "1.0.0"
scripts:
  directory: scripts
groups:
  public:
    summary: Default public project commands.
  development:
    summary: Development commands.
commands:
  - uid: dev-start
    id: start
    group: development
    summary: Start development.
    description: Starts local development.
    command: echo start
  - uid: check
    id: check
    group: development
    summary: Check source.
    description: Checks source files.
    command: echo check
`

const emptyManifest = (version = '1.0.0'): string => `version: "${version}"
scripts:
  directory: scripts
groups:
  public:
    summary: Default public project commands.
commands: []
`

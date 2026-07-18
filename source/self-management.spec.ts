import { afterEach, describe, expect, test } from 'bun:test'
import { chmod, mkdtemp, readdir, rename, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { upgradeSelf, validateNativeBinary } from './self-management.js'
import { readVersion } from './help.js'

const directories: string[] = []
const originalFetch = globalThis.fetch
const originalSelfPath = process.env['RUNX_SELF_PATH']

afterEach(async () => {
  globalThis.fetch = originalFetch
  if (originalSelfPath === undefined) delete process.env['RUNX_SELF_PATH']
  else process.env['RUNX_SELF_PATH'] = originalSelfPath
  await Promise.all(directories.splice(0).map((directory) => rm(directory, { recursive: true, force: true })))
})

describe('RunX self-management', () => {
  test('returns a nullable-plan fallback envelope for discovery failures', async () => {
    const directory = await temporaryDirectory()
    process.env['RUNX_SELF_PATH'] = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    globalThis.fetch = (async () => new Response('offline', { status: 503 })) as typeof fetch
    const result = await upgradeSelf(false)
    expect(result).toMatchObject({
      schemaVersion: 1,
      command: 'runx upgrade',
      outcome: 'failed',
      plan: null,
      result: null,
      recovery: { targetVersion: readVersion(), targetSource: 'fallback-current' },
      error: { code: 'release_lookup_failed', phase: 'plan' },
    })
    expect(result.recovery.stopProcessCommand).toMatch(/runx/)
  })

  test('uses the stable payload error code for malformed discovery data', async () => {
    const directory = await temporaryDirectory()
    process.env['RUNX_SELF_PATH'] = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    globalThis.fetch = (async () => Response.json({ releases: [] })) as typeof fetch
    const result = await upgradeSelf(false)
    expect(result.error?.code).toBe('release_payload_invalid')
    expect(result.plan).toBeNull()
    expect(result.recovery.targetSource).toBe('fallback-current')
  })

  test('returns pinned recovery when invoked through the Bun launcher', async () => {
    delete process.env['RUNX_SELF_PATH']
    mockUpgrade('99.0.0', Bun.file(process.execPath))
    const result = await upgradeSelf(false)
    expect(result.outcome).toBe('failed')
    expect(result.plan?.targetVersion).toBe('99.0.0')
    expect(result.result).toBeNull()
    expect(result.recovery.targetSource).toBe('resolved')
    expect(result.error?.message).toContain('requires a native RunX executable')
    expect(result.recovery.targetVersion).toBe('99.0.0')
    expect(result.recovery.installCommand).toContain('99.0.0')
  })

  test('emits the plan before awaiting the download and verifies the canonical executable', async () => {
    const directory = await temporaryDirectory()
    const executablePath = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    await Bun.write(executablePath, Bun.file(process.execPath))
    if (process.platform !== 'win32') await chmod(executablePath, 0o755)
    process.env['RUNX_SELF_PATH'] = executablePath
    const order: string[] = []
    mockUpgrade(process.versions.bun, Bun.file(process.execPath), order)

    const result = await upgradeSelf({
      dryRun: false,
      onPlan: () => order.push('plan-rendered'),
      onEvent: (event) => {
        if (event.status === 'started' && event.phase !== 'plan') order.push(event.phase)
      },
    })

    expect(result.outcome).toBe('upgraded')
    expect(result.error).toBeNull()
    expect(await executableVersion(executablePath)).toBe(process.versions.bun)
    expect(order.indexOf('plan-rendered')).toBeLessThan(order.indexOf('download'))
    expect(order.indexOf('download')).toBeLessThan(order.indexOf('download-body'))
    expect(order.indexOf('download-body')).toBeLessThan(order.indexOf('validate'))
    expect(order.indexOf('validate')).toBeLessThan(order.indexOf('replace'))
    expect(order.indexOf('replace')).toBeLessThan(order.indexOf('verify'))
    expect((await readdir(directory)).some((name) => name.includes('.new-'))).toBe(false)
  }, 15_000)

  test('rejects a non-native download before replacing the executable', async () => {
    const directory = await temporaryDirectory()
    const executablePath = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    const originalBytes = new Uint8Array([1, 2, 3, 4])
    await Bun.write(executablePath, originalBytes)
    process.env['RUNX_SELF_PATH'] = executablePath
    mockUpgrade('99.0.0', new Blob(['not-an-executable']))

    const result = await upgradeSelf(false)

    expect(result.outcome).toBe('failed')
    expect(result.error?.phase).toBe('validate')
    expect(result.error?.code).toBe('download_invalid')
    expect(await Bun.file(executablePath).bytes()).toEqual(originalBytes)
  })

  test('restores the previous executable when target verification fails', async () => {
    const directory = await temporaryDirectory()
    const executablePath = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    const originalBytes = new Uint8Array([0x4d, 0x5a, 0x01, 0x02])
    await Bun.write(executablePath, originalBytes)
    process.env['RUNX_SELF_PATH'] = executablePath
    mockUpgrade('99.0.0', new Blob([nativeGarbage()]))

    const result = await upgradeSelf(false)

    expect(result.outcome).toBe('rolled-back')
    expect(result.error?.phase).toBe('verify')
    expect(result.error?.code).toBe('verification_failed')
    expect(result.result?.installedVersion).toBe(readVersion())
    expect(await Bun.file(executablePath).bytes()).toEqual(originalBytes)
    expect((await readdir(directory)).some((name) => name.includes('.new-') || name.includes('.old-'))).toBe(false)
  })

  test('rolls back deterministically when the second rename fails', async () => {
    const directory = await temporaryDirectory()
    const executablePath = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    const originalBytes = new Uint8Array([1, 2, 3, 4])
    await Bun.write(executablePath, originalBytes)
    process.env['RUNX_SELF_PATH'] = executablePath
    mockUpgrade('99.0.0', new Blob([nativeGarbage()]))
    let renameCount = 0

    const result = await upgradeSelf({
      dryRun: false,
      fileOperations: {
        rename: async (from, to) => {
          renameCount += 1
          if (renameCount === 2) throw new Error('controlled second rename failure')
          await rename(from, to)
        },
        remove: async (path) => rm(path, { force: true }),
        makeExecutable: async (path) => chmod(path, 0o755),
      },
    })

    expect(result.outcome).toBe('rolled-back')
    expect(result.error?.code).toBe('replace_failed')
    expect(result.error?.message).toContain('controlled second rename failure')
    expect(await Bun.file(executablePath).bytes()).toEqual(originalBytes)
    expect(renameCount).toBe(3)
  })

  test('reports both paths and causes when second rename and rollback fail', async () => {
    const directory = await temporaryDirectory()
    const executablePath = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    await Bun.write(executablePath, new Uint8Array([1, 2, 3, 4]))
    process.env['RUNX_SELF_PATH'] = executablePath
    mockUpgrade('99.0.0', new Blob([nativeGarbage()]))
    let renameCount = 0

    const result = await upgradeSelf({
      dryRun: false,
      fileOperations: {
        rename: async (from, to) => {
          renameCount += 1
          if (renameCount === 2) throw new Error('controlled replacement failure')
          if (renameCount === 3) throw new Error('controlled rollback failure')
          await rename(from, to)
        },
        remove: async (path) => rm(path, { force: true }),
        makeExecutable: async (path) => chmod(path, 0o755),
      },
    })

    expect(result.outcome).toBe('failed')
    expect(result.error?.code).toBe('rollback_failed')
    expect(result.error?.message).toContain('controlled replacement failure')
    expect(result.error?.message).toContain('controlled rollback failure')
    expect(result.error?.message).toContain(`Canonical path: ${executablePath}`)
    expect(result.error?.message).toContain('Backup path:')
    expect(renameCount).toBe(3)
  })

  test('does not downgrade stable or prerelease installations', async () => {
    const directory = await temporaryDirectory()
    const executablePath = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    await Bun.write(executablePath, new Uint8Array([1]))
    process.env['RUNX_SELF_PATH'] = executablePath

    for (const currentVersion of ['2.0.0', '2.0.0-alpha.2']) {
      mockReleases([
        releaseFixture('1.9.0', false),
        releaseFixture(currentVersion, currentVersion.includes('-')),
      ])
      const result = await upgradeSelf({ dryRun: false, currentVersion })
      expect(result.outcome).toBe('up-to-date')
      expect(result.plan?.targetVersion).toBe(currentVersion)
      expect(result.result?.installedVersion).toBe(currentVersion)
    }
  })

  test('still selects a stable release above the installed prerelease', async () => {
    const directory = await temporaryDirectory()
    process.env['RUNX_SELF_PATH'] = join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')
    await Bun.write(process.env['RUNX_SELF_PATH'], new Uint8Array([1]))
    mockReleases([releaseFixture('2.0.0', false), releaseFixture('2.0.0-alpha.2', true)])
    const result = await upgradeSelf({ dryRun: true, currentVersion: '2.0.0-alpha.2' })
    expect(result.outcome).toBe('dry-run')
    expect(result.plan?.targetVersion).toBe('2.0.0')
  })

  test('validates the native magic for the selected operating system', () => {
    expect(() => validateNativeBinary(new Uint8Array([0x4d, 0x5a]), 'windows')).not.toThrow()
    expect(() => validateNativeBinary(new Uint8Array([0x7f, 0x45, 0x4c, 0x46]), 'linux')).not.toThrow()
    expect(() => validateNativeBinary(new Uint8Array([0xcf, 0xfa, 0xed, 0xfe]), 'darwin')).not.toThrow()
    expect(() => validateNativeBinary(new Uint8Array([0x3c, 0x68, 0x74, 0x6d]), 'windows')).toThrow('not a native')
  })

  if (process.platform === 'win32') {
    test('replaces a running Windows executable before reporting success', async () => {
      const directory = await temporaryDirectory()
      const executablePath = join(directory, 'runx.exe')
      const powershellPath = Bun.which('powershell.exe')
      if (!powershellPath) throw new Error('powershell.exe is required for the Windows replacement test')
      await Bun.write(executablePath, Bun.file(powershellPath))
      const runningExecutable = Bun.spawn([executablePath, '-NoProfile', '-Command', 'Start-Sleep -Seconds 30'], {
        stdin: 'ignore', stdout: 'ignore', stderr: 'ignore',
      })

      try {
        await Bun.sleep(250)
        process.env['RUNX_SELF_PATH'] = executablePath
        mockUpgrade(process.versions.bun, Bun.file(process.execPath))
        const result = await upgradeSelf(false)
        expect(result.outcome).toBe('upgraded')
        expect(await executableVersion(executablePath)).toBe(process.versions.bun)
        expect(result.events.find((event) => event.phase === 'cleanup')?.status).toBe('started')
      } finally {
        runningExecutable.kill()
        await runningExecutable.exited
      }
    }, 15_000)
  }
})

const temporaryDirectory = async (): Promise<string> => {
  const directory = await mkdtemp(join(tmpdir(), 'runx-self-management-'))
  directories.push(directory)
  return directory
}

const mockUpgrade = (version: string, binary: Blob, order?: string[]): void => {
  let request = 0
  globalThis.fetch = (async () => {
    request += 1
    if (request === 1) {
      return Response.json([{
        tag_name: `@guiho/runx@${version}`,
        draft: false,
        prerelease: false,
        published_at: '2026-07-15T00:00:00Z',
        assets: [{ name: currentAssetName(), browser_download_url: 'https://example.test/runx' }],
      }])
    }
    return {
      ok: true,
      status: 200,
      arrayBuffer: async () => {
        order?.push('download-body')
        return binary.arrayBuffer()
      },
    } as Response
  }) as typeof fetch
}

const mockReleases = (releases: unknown[]): void => {
  globalThis.fetch = (async () => Response.json(releases)) as typeof fetch
}

const releaseFixture = (version: string, prerelease: boolean) => ({
  tag_name: `@guiho/runx@${version}`,
  draft: false,
  prerelease,
  published_at: '2026-07-15T00:00:00Z',
  assets: [{ name: currentAssetName(), browser_download_url: `https://example.test/${version}` }],
})

const currentAssetName = (): string => {
  const os = process.platform === 'win32' ? 'windows' : process.platform === 'darwin' ? 'darwin' : 'linux'
  const suffix = process.platform === 'win32' ? '.exe' : ''
  return process.arch === 'arm64' ? `runx-${os}-arm64${suffix}` : `runx-${os}-x64-baseline${suffix}`
}

const nativeGarbage = (): Uint8Array => {
  if (process.platform === 'win32') return new Uint8Array([0x4d, 0x5a, 0, 0])
  if (process.platform === 'darwin') return new Uint8Array([0xcf, 0xfa, 0xed, 0xfe, 0, 0])
  return new Uint8Array([0x7f, 0x45, 0x4c, 0x46, 0, 0])
}

const executableVersion = async (executablePath: string): Promise<string> => {
  const child = Bun.spawn([executablePath, '--version'], { stdin: 'ignore', stdout: 'pipe', stderr: 'pipe' })
  const [stdout, exitCode] = await Promise.all([new Response(child.stdout).text(), child.exited])
  if (exitCode !== 0) throw new Error(`Replacement verification exited with code ${exitCode}`)
  return stdout.trim()
}
